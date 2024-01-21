package plugin

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func (a *ServicePlugin) GetPluginPublicKeyFiles() ([]string, *model_helper.AppError) {
	return a.srv.Config().PluginSettings.SignaturePublicKeyFiles, nil
}

func (a *ServicePlugin) GetPublicKey(name string) ([]byte, *model_helper.AppError) {
	data, err := a.srv.ConfigStore.GetFile(name)
	if err != nil {
		return nil, model_helper.NewAppError("GetPublicKey", "app.plugin.get_public_key.get_file.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return data, nil
}

// AddPublicKey will add plugin public key to the config. Overwrites the previous file
func (a *ServicePlugin) AddPublicKey(name string, key io.Reader) *model_helper.AppError {
	if model.IsSamlFile(&a.srv.Config().SamlSettings, name) {
		return model_helper.NewAppError("AddPublicKey", "app.plugin.modify_saml.app_error", nil, "", http.StatusInternalServerError)
	}
	data, err := io.ReadAll(key)
	if err != nil {
		return model_helper.NewAppError("AddPublicKey", "app.plugin.write_file.read.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	err = a.srv.ConfigStore.SetFile(name, data)
	if err != nil {
		return model_helper.NewAppError("AddPublicKey", "app.plugin.write_file.saving.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.srv.UpdateConfig(func(cfg *model.Config) {
		if !cfg.PluginSettings.SignaturePublicKeyFiles.Contains(name) {
			cfg.PluginSettings.SignaturePublicKeyFiles = append(cfg.PluginSettings.SignaturePublicKeyFiles, name)
		}
	})

	return nil
}

// DeletePublicKey will delete plugin public key from the config.
func (a *ServicePlugin) DeletePublicKey(name string) *model_helper.AppError {
	if model.IsSamlFile(&a.srv.Config().SamlSettings, name) {
		return model_helper.NewAppError("AddPublicKey", "app.plugin.modify_saml.app_error", nil, "", http.StatusInternalServerError)
	}
	filename := filepath.Base(name)
	if err := a.srv.ConfigStore.RemoveFile(filename); err != nil {
		return model_helper.NewAppError("DeletePublicKey", "app.plugin.delete_public_key.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	a.srv.UpdateConfig(func(cfg *model.Config) {
		cfg.PluginSettings.SignaturePublicKeyFiles = cfg.PluginSettings.SignaturePublicKeyFiles.Remove(filename)
	})

	return nil
}

func (a *ServicePlugin) VerifyPlugin(plugin, signature io.ReadSeeker) *model_helper.AppError {
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
	return model_helper.NewAppError("VerifyPlugin", "api.plugin.verify_plugin.app_error", nil, "", http.StatusInternalServerError)
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
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't read the file")
	}
	block, err := armor.Decode(bytes.NewReader(readBytes))
	if err != nil {
		return bytes.NewReader(readBytes), nil
	}
	return block.Body, nil
}
