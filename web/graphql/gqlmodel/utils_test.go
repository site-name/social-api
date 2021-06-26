package gqlmodel

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/site-name/i18naddress"
)

func TestI18nAddressValidationRulesToGraphql(t *testing.T) {
	params := &i18naddress.Params{
		CountryCode: "vn",
	}
	rules, err := i18naddress.GetValidationRules(params)
	if err != nil {
		t.Fatal(err)
	}

	dt, err := json.Marshal(rules)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dt))

	res := I18nAddressValidationRulesToGraphql(rules)
	dt, err = json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(dt))
}
