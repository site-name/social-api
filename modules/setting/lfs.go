package setting

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"time"

	"code.gitea.io/gitea/modules/generate"
	"code.gitea.io/gitea/modules/log"
	"github.com/sitename/sitename/modules/util"

	ini "gopkg.in/ini.v1"
)

// LFS represents the configuration for Git LFS
var LFS = struct {
	StartServer     bool          `ini:"LFS_START_SERVER"`
	JWTSecretBase64 string        `ini:"LFS_JWT_SECRET"`
	JWTSecretBytes  []byte        `ini:"-"`
	HTTPAuthExpiry  time.Duration `ini:"LFS_HTTP_AUTH_EXPIRY"`
	MaxFileSize     int64         `ini:"LFS_MAX_FILE_SIZE"`
	LocksPagingNum  int           `ini:"LFS_LOCKS_PAGING_NUM"`

	Storage
}{}

func newLFSService() {
	sec := Cfg.Section("server")
	if err := sec.MapTo(&LFS); err != nil {
		log.Fatal("Failed to map LFS settings: %v", err)
	}

	lfsSec := Cfg.Section("lfs")
	storageType := lfsSec.Key("STORAGE_TYPE").MustString("")

	// Specifically default PATH to LFS_CONTENT_PATH
	lfsSec.Key("PATH").MustString(
		sec.Key("LFS_CONTENT_PATH").String())

	LFS.Storage = getStorage("lfs", storageType, lfsSec)

	// Rest of LFS service settings
	if LFS.LocksPagingNum == 0 {
		LFS.LocksPagingNum = 50
	}

	LFS.HTTPAuthExpiry = sec.Key("LFS_HTTP_AUTH_EXPIRY").MustDuration(20 * time.Minute)

	if LFS.StartServer {
		LFS.JWTSecretBytes = make([]byte, 32)
		n, err := base64.RawURLEncoding.Decode(LFS.JWTSecretBytes, []byte(LFS.JWTSecretBase64))

		if err != nil || n != 32 {
			LFS.JWTSecretBase64, err = generate.NewJwtSecret()
			if err != nil {
				log.Fatal("Error generating JWT Secret for custom config: %v", err)
				return
			}

			// Save secret
			cfg := ini.Empty()
			isFile, err := util.IsFile(CustomConf)
			if err != nil {
				log.Error("Unable to check if %s is a file. Error: %v", CustomConf, err)
			}
			if isFile {
				// Keeps custom settings if there is already something.
				if err := cfg.Append(CustomConf); err != nil {
					log.Error("Failed to load custom conf '%s': %v", CustomConf, err)
				}
			}

			cfg.Section("server").Key("LFS_JWT_SECRET").SetValue(LFS.JWTSecretBase64)

			if err := os.MkdirAll(filepath.Dir(CustomConf), os.ModePerm); err != nil {
				log.Fatal("Failed to create '%s': %v", CustomConf, err)
			}
			if err := cfg.SaveTo(CustomConf); err != nil {
				log.Fatal("Error saving generated JWT Secret to custom config: %v", err)
				return
			}
		}
	}
}
