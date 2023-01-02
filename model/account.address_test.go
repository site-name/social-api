package model

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressIsValid(t *testing.T) {
	addr := new(Address)
	err := addr.IsValid()

	require.True(t, HasExpectedAddressIsValidError(err, "id", "", addr.Id), "expected address is valid error: %s", err.Error())
}

func HasExpectedAddressIsValidError(err *AppError, fieldName, addressID string, fieldValue any) bool {
	if err == nil {
		return false
	}

	return err.Where == "Address.IsValid" &&
		err.Id == fmt.Sprintf("address.is_valid.%s.app_error", fieldName) &&
		err.StatusCode == http.StatusBadRequest
}
