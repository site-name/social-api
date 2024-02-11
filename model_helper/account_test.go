package model_helper

import (
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/stretchr/testify/require"
)

func TestAddressIsValid(t *testing.T) {
	addr := model.Address{
		FirstName: "lol",
		LastName:  "lol",
	}
	AddressPreSave(&addr)
	err := AddressIsValid(addr)
	require.NotNil(t, err, "error should be non-nil")
}
