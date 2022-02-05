package util

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"
)

const (
	defaultKeySalt       = "ThisIsKeySaltForToken123456@#Lol"
	charSet              = "0123456789abcdefghijklmnopqrstuvwxyz"
	passwordResetTimeout = 24 * 60 * 60
)

var DefaultTokenGenerator = NewTokenGenerator("", crypto.SHA1, "")

type Hashable interface {
	GetId() string
	GetPassword() string
	GetLastLogin() time.Time
	GetEmail() string
}

// TokenGenerator is borrowed from django's
type TokenGenerator struct {
	keySalt   string
	algorithm crypto.Hash
	secret    string
}

func NewTokenGenerator(salt string, algo crypto.Hash, secret string) *TokenGenerator {
	if salt == "" {
		salt = defaultKeySalt
	}
	if secret == "" {
		secret = defaultKeySalt
	}

	return &TokenGenerator{
		keySalt:   salt,
		algorithm: algo,
		secret:    secret,
	}
}

func (p *TokenGenerator) MakeToken(h Hashable) string {
	return p.makeTokenWithTimestamp(h, p.numSeconds(p.now()), false)
}

// Check that a password reset token is correct for a given user.
func (p *TokenGenerator) CheckToken(h Hashable, token string) bool {
	if h == nil && token == "" {
		return false
	}

	splitToken := strings.Split(token, "-")
	if len(splitToken) != 2 {
		return false
	}

	tsB36 := splitToken[0]
	legacyToken := len(tsB36) < 4

	ts, err := base36ToInt(tsB36)
	if err != nil {
		return false
	}

	// Check that the timestamp/uid has not been tampered with
	if subtle.ConstantTimeCompare(
		[]byte(p.makeTokenWithTimestamp(h, ts, false)),
		[]byte(token),
	) == 0 {

		if subtle.ConstantTimeCompare(
			[]byte(p.makeTokenWithTimestamp(h, ts, true)),
			[]byte(token),
		) == 0 {
			return false
		}
	}

	now := p.now()
	if legacyToken {
		ts *= 24 * 60 * 60
		ts += now.Sub(
			time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				0, 0, 0, 0,
				now.Location(),
			),
		).Milliseconds() / 1000
	}

	if p.numSeconds(now)-ts > int64(passwordResetTimeout) {
		return false
	}

	return true
}

func (p *TokenGenerator) makeTokenWithTimestamp(h Hashable, timestamp int64, legacy bool) string {
	algo := crypto.SHA1
	if !legacy {
		algo = p.algorithm
	}

	hasher := saltedHmac(
		p.keySalt,
		p.makeHashValue(h, timestamp),
		p.secret,
		algo,
	)

	hashString := ""

	for i, item := range hex.EncodeToString(hasher.Sum(nil)) {
		if i%2 == 0 {
			hashString += string(item)
		}
	}

	return intToBase36(timestamp) + "-" + hashString
}

func (p *TokenGenerator) makeHashValue(h Hashable, timestamp int64) string {
	l := h.GetLastLogin()
	loginTimestamp := time.Date(
		l.Year(),
		l.Month(),
		l.Day(),
		l.Hour(),
		l.Minute(),
		l.Second(),
		0,
		l.Location(),
	).String()

	return fmt.Sprintf("%s%s%s%d%s", h.GetId(), h.GetPassword(), loginTimestamp, timestamp, h.GetEmail())
}

func (p *TokenGenerator) numSeconds(t time.Time) int64 {
	return t.
		Sub(time.Date(2001, 1, 1, 0, 0, 0, 0, t.Location())).
		Milliseconds() / 1000
}

func (p *TokenGenerator) now() time.Time {
	return time.Now()
}

// intToBase36 Convert an integer to a base36 string
func intToBase36(i int64) string {
	if i < 0 {
		i = 0
	}
	if i < 36 {
		return string(charSet[i])
	}

	var (
		b36 = ""
		n   int64
	)
	for i != 0 {
		i, n = i/36, i%36
		b36 = string(charSet[n]) + b36
	}

	return b36
}

func saltedHmac(keySalt, value, secret string, algorithm crypto.Hash) hash.Hash {
	if secret == "" {
		secret = defaultKeySalt
	}
	if keySalt == "" {
		keySalt = defaultKeySalt
	}

	var newHashFunc func() hash.Hash

	switch algorithm {
	case crypto.MD5:
		newHashFunc = md5.New
	case crypto.SHA1:
		newHashFunc = sha1.New
	case crypto.SHA256:
		newHashFunc = sha256.New

	default:
		newHashFunc = sha512.New
	}

	hasher := hmac.New(newHashFunc, []byte(keySalt+secret))
	hasher.Write([]byte(value))

	return hasher
}

func base36ToInt(s string) (int64, error) {
	if len(s) > 13 {
		return 0, errors.New("base36 input too large")
	}
	return strconv.ParseInt(s, 36, 36)
}
