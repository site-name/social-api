package routers

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/sitename/sitename/models"
// 	"github.com/sitename/sitename/modules/cache"
// 	// "github.com/sitename/sitename/modules/cron"
// 	"github.com/sitename/sitename/modules/eventsource"
// 	"github.com/sitename/sitename/modules/log"
// 	"github.com/sitename/sitename/modules/setting"

// 	// "github.com/sitename/sitename/modules/storage"
// 	"github.com/sitename/sitename/services/mailer"
// )

// // NewServices init new services
// func NewServices() {
// 	// run setting services
// 	setting.NewServices()
// 	// if err := storage.Init(); err != nil {
// 	// 	log.Fatal("storage init failed: %v", err)
// 	// }

// 	// run mailer service
// 	mailer.NewContext()
// 	_ = cache.NewContext()
// }

// // GlobalInit is for global configuration reload-able.
// func GlobalInit(ctx context.Context) {
// 	setting.NewContext()

// 	log.Trace("AppPath: %s", setting.AppPath)
// 	log.Trace("AppWorkPath: %s", setting.AppWorkPath)
// 	log.Trace("Custom path: %s", setting.CustomPath)
// 	log.Trace("Log path: %s", setting.LogRootPath)

// 	NewServices()

// 	if err := initDBEngine(ctx); err == nil {
// 		log.Info("ORM engine initialization successful!")
// 	} else {
// 		log.Fatal("ORM engine initialization failed: %v", err)
// 	}

// 	if err := models.InitOAuth2(); err != nil {
// 		log.Fatal("Failed to initialize OAuth2 support: %v", err)
// 	}

// 	// Booting long running goroutines
// 	cron.NewContext()

// 	// eventsource
// 	eventsource.GetManager().Init()
// }

// // In case of problems connecting to DB, retry connection. Eg, PGSQL in Docker Container on Synology
// func initDBEngine(ctx context.Context) (err error) {
// 	log.Info("Beginning ORM engine initialization.")
// 	for i := 0; i < setting.Database.DBConnectRetries; i++ {
// 		select {
// 		case <-ctx.Done():
// 			return fmt.Errorf("Aborted due to shutdown:\nin retry ORM engine initialization")
// 		default:
// 		}
// 		log.Info("ORM engine initialization attempt #%d/%d...", i+1, setting.Database.DBConnectRetries)
// 		if err = models.NewEngine(ctx, nil); err == nil {
// 			break
// 		} else if err != nil && i == setting.Database.DBConnectRetries-1 {
// 			return err
// 		}
// 		log.Error("ORM engine initialization attempt #%d/%d failed. Error: %v", i+1, setting.Database.DBConnectRetries, err)
// 		log.Info("Backing off for %d seconds", int64(setting.Database.DBConnectBackoff/time.Second))
// 		time.Sleep(setting.Database.DBConnectBackoff)
// 	}
// 	models.HasEngine = true
// 	return nil
// }
