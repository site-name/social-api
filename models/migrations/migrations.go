package migrations

import (
	"fmt"
	"os"

	"github.com/sitename/sitename/modules/log"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const minDBVersion = 70

// Migration describes on migration from lower version to high version
type Migration interface {
	Description() string
	Migrate(*xorm.Engine) error
}

type migration struct {
	description string
	migrate     func(*xorm.Engine) error
}

// NewMigration creates a new migration
func NewMigration(desc string, fn func(*xorm.Engine) error) Migration {
	return &migration{desc, fn}
}

// Description returns the migration's description
func (m *migration) Description() string {
	return m.description
}

// Migrate executes the migration
func (m *migration) Migrate(x *xorm.Engine) error {
	return m.migrate(x)
}

// Version describes the version table. Should have only one row with id==1
type Version struct {
	ID      int64 `xorm:"pk autoincr"`
	Version int64
}

// This is a sequence of migrations. Add new migrations to the bottom of the list.
// If you want to "retire" a migration, remove it from the top of the list and
// update minDBVersion accordingly
var migrations = []Migration{}

// Migrate database to current version
func Migrate(x *xorm.Engine) error {
	// Set a new clean the default mapper to GonicMapper as that is the default for Sitename.
	x.SetMapper(names.GonicMapper{})
	if err := x.Sync(new(Version)); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	currentVersion := &Version{ID: 1}
	has, err := x.Get(currentVersion)
	if err != nil {
		return fmt.Errorf("get: %v", err)
	} else if !has {
		// If the version record does not exist we think
		// it is a fresh installation and we can skip all migrations.
		currentVersion.ID = 0
		currentVersion.Version = int64(minDBVersion + len(migrations))

		if _, err = x.InsertOne(currentVersion); err != nil {
			return fmt.Errorf("insert: %v", err)
		}
	}

	v := currentVersion.Version
	if minDBVersion > v {
		log.Fatal(`Sitename no longer supports auto-migration from your previously installed version.
Please try upgrading to a lower version first (suggested v1.6.4), then upgrade to this version.`)
		return nil
	}

	// Downgrading Sitename's database version not supported
	if int(v-minDBVersion) > len(migrations) {
		msg := fmt.Sprintf("Downgrading database version from '%d' to '%d' is not supported and may result in loss of data integrity.\nIf you really know what you're doing, execute `UPDATE version SET version=%d WHERE id=1;`\n",
			v, minDBVersion+len(migrations), minDBVersion+len(migrations))
		fmt.Fprint(os.Stderr, msg)
		log.Fatal(msg)
		return nil
	}

	// Migrate
	for i, m := range migrations[v-minDBVersion:] {
		log.Info("Migration[%d]: %s", v+int64(i), m.Description())
		// Reset the mapper between each migration - migrations are not supposed to depend on each other
		x.SetMapper(names.GonicMapper{})
		if err = m.Migrate(x); err != nil {
			return fmt.Errorf("do migrate: %v", err)
		}
		currentVersion.Version = v + int64(i) + 1
		if _, err = x.ID(1).Update(currentVersion); err != nil {
			return err
		}
	}
	return nil
}
