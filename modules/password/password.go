package password

import (
	goContext "context"
	"crypto/rand"
	"math/big"
	"strings"
	"sync"

	"github.com/sitename/sitename/modules/setting"
)

// complexity contains information about a particular kind of password complexity
type complexity struct {
	ValidChars string
	TrNameOne  string
}

var (
	matchComplexityOnce sync.Once
	validChars          string
	requiredList        []complexity

	charComplexities = map[string]complexity{
		"lower": {
			`abcdefghijklmnopqrstuvwxyz`,
			"form.password_lowercase_one",
		},
		"upper": {
			`ABCDEFGHIJKLMNOPQRSTUVWXYZ`,
			"form.password_uppercase_one",
		},
		"digit": {
			`0123456789`,
			"form.password_digit_one",
		},
		"spec": {
			` !"#$%&'()*+,-./:;<=>?@[\]^_{|}~` + "`",
			"form.password_special_one",
		},
	}
)

// NewComplexity for preparation
func NewComplexity() {
	matchComplexityOnce.Do(func() {
		setupComplexity(setting.PasswordComplexity)
	})
}

func setupComplexity(values []string) {
	if len(values) != 1 || values[0] != "off" {
		for _, val := range values {
			if complex, ok := charComplexities[val]; ok {
				validChars += complex.ValidChars
				requiredList = append(requiredList, complex)
			}
		}
		if len(requiredList) == 0 {
			// No valid character classes found; use all classes as default
			for _, complex := range charComplexities {
				validChars += complex.ValidChars
				requiredList = append(requiredList, complex)
			}
		}
	}
	if validChars == "" {
		// No complexities to check; provide a sensible default for password generation
		validChars = charComplexities["lower"].ValidChars + charComplexities["upper"].ValidChars + charComplexities["digit"].ValidChars
	}
}

// IsComplexEnough returns true if password meets complexity settings
func IsComplexEnough(password string) bool {
	NewComplexity()
	if len(validChars) > 0 {
		for _, req := range requiredList {
			if !strings.ContainsAny(req.ValidChars, password) {
				return false
			}
		}
	}
	return true
}

// Generate a random password
func Generate(n int) (string, error) {
	NewComplexity()
	buffer := make([]byte, n)
	max := big.NewInt(int64(len(validChars)))
	for {
		for j := 0; j < n; j++ {
			rnd, err := rand.Int(rand.Reader, max)
			if err != nil {
				return "", err
			}
			buffer[j] = validChars[rnd.Int64()]
		}
		pwned, err := IsPwned(goContext.Background(), string(buffer))
		if err != nil {
			return "", err
		}
		if IsComplexEnough(string(buffer)) && !pwned && string(buffer[0]) != " " && string(buffer[n-1]) != " " {
			return string(buffer), nil
		}
	}
}
