// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package configservice

import (
	"crypto/ecdsa"

	"github.com/sitename/sitename/model_helper"
)

// An interface representing something that contains a Config, such as the app.App struct
type ConfigService interface {
	Config() *model_helper.Config
	AddConfigListener(func(old, current *model_helper.Config)) string
	RemoveConfigListener(string)
	AsymmetricSigningKey() *ecdsa.PrivateKey
}
