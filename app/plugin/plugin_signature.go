package plugin

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/modules/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func (a *AppPlugin) GetPluginPublicKeyFiles() ([]string, *model.AppError) {
	return a.Config().PluginSettings.SignaturePublicKeyFiles, nil
}

func (a *AppPlugin) GetPublicKey(name string) ([]byte, *model.AppError) {
	data, err := a.Srv().ConfigStore.GetFile(name)
	if err != nil {
		return nil, model.NewAppError("GetPublicKey", "app.plugin.get_public_key.get_file.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return data, nil
}

// AddPublicKey will add plugin public key to the config. Overwrites the previous file
func (a *AppPlugin) AddPublicKey(name string, key io.Reader) *model.AppError {
	if model.IsSamlFile(&a.Config().SamlSettings, name) {
		return model.NewAppError("AddPublicKey", "app.plugin.modify_saml.app_error", nil, "", http.StatusInternalServerError)
	}
	data, err := ioutil.ReadAll(key)
	if err != nil {
		return model.NewAppError("AddPublicKey", "app.plugin.write_file.read.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	err = a.Srv().ConfigStore.SetFile(name, data)
	if err != nil {
		return model.NewAppError("AddPublicKey", "app.plugin.write_file.saving.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.UpdateConfig(func(cfg *model.Config) {
		if !util.StringInSlice(name, cfg.PluginSettings.SignaturePublicKeyFiles) {
			cfg.PluginSettings.SignaturePublicKeyFiles = append(cfg.PluginSettings.SignaturePublicKeyFiles, name)
		}
	})

	return nil
}

// DeletePublicKey will delete plugin public key from the config.
func (a *AppPlugin) DeletePublicKey(name string) *model.AppError {
	if model.IsSamlFile(&a.Config().SamlSettings, name) {
		return model.NewAppError("AddPublicKey", "app.plugin.modify_saml.app_error", nil, "", http.StatusInternalServerError)
	}
	filename := filepath.Base(name)
	if err := a.Srv().ConfigStore.RemoveFile(filename); err != nil {
		return model.NewAppError("DeletePublicKey", "app.plugin.delete_public_key.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.SignaturePublicKeyFiles = util.RemoveStringFromSlice(filename, cfg.PluginSettings.SignaturePublicKeyFiles)
	})

	return nil
}

func (a *AppPlugin) VerifyPlugin(plugin, signature io.ReadSeeker) *model.AppError {
	if err := verifySignature(bytes.NewReader(sitenamePluginPublicKey), plugin, signature); err == nil {
		return nil
	}
	publicKeys, appErr := a.GetPluginPublicKeyFiles()
	if appErr != nil {
		return appErr
	}
	for _, pk := range publicKeys {
		pkBytes, appErr := a.GetPublicKey(pk)
		if appErr != nil {
			slog.Warn("Unable to get public key for ", slog.String("filename", pk))
			continue
		}
		publicKey := bytes.NewReader(pkBytes)
		plugin.Seek(0, 0)
		signature.Seek(0, 0)
		if err := verifySignature(publicKey, plugin, signature); err == nil {
			return nil
		}
	}
	return model.NewAppError("VerifyPlugin", "api.plugin.verify_plugin.app_error", nil, "", http.StatusInternalServerError)
}

func verifySignature(publicKey, message, signatrue io.Reader) error {
	pk, err := decodeIfArmored(publicKey)
	if err != nil {
		return errors.Wrap(err, "can't decode public key")
	}
	s, err := decodeIfArmored(signatrue)
	if err != nil {
		return errors.Wrap(err, "can't decode signature")
	}
	return verifyBinarySignature(pk, message, s)
}

func verifyBinarySignature(publicKey, signedFile, signature io.Reader) error {
	keyring, err := openpgp.ReadKeyRing(publicKey)
	if err != nil {
		return errors.Wrap(err, "can't read public key")
	}
	if _, err = openpgp.CheckDetachedSignature(keyring, signedFile, signature); err != nil {
		return errors.Wrap(err, "error while checking the signature")
	}
	return nil
}

func decodeIfArmored(reader io.Reader) (io.Reader, error) {
	readBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't read the file")
	}
	block, err := armor.Decode(bytes.NewReader(readBytes))
	if err != nil {
		return bytes.NewReader(readBytes), nil
	}
	return block.Body, nil
}
