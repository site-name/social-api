package util

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// OptionalBool a boolean that can be "null"
type OptionalBool byte

const (
	// OptionalBoolNone a "null" boolean value
	OptionalBoolNone = iota
	// OptionalBoolTrue a "true" boolean value
	OptionalBoolTrue
	// OptionalBoolFalse a "false" boolean value
	OptionalBoolFalse
)

// IsTrue return true if equal to OptionalBoolTrue
func (o OptionalBool) IsTrue() bool {
	return o == OptionalBoolTrue
}

// IsFalse return true if equal to OptionalBoolFalse
func (o OptionalBool) IsFalse() bool {
	return o == OptionalBoolFalse
}

// IsNone return true if equal to OptionalBoolNone
func (o OptionalBool) IsNone() bool {
	return o == OptionalBoolNone
}

// OptionalBoolOf get the corresponding OptionalBool of a bool
func OptionalBoolOf(b bool) OptionalBool {
	if b {
		return OptionalBoolTrue
	}
	return OptionalBoolFalse
}

type Ordered interface {
	uint8 |
		int |
		uint |
		int8 |
		int16 |
		int32 |
		int64 |
		uint16 |
		uint32 |
		uint64 |
		float32 |
		float64 |
		~string
}

// Max accepts any number of arguments of any type and returns max value
func Max[T Ordered](a ...T) T {
	if len(a) == 0 {
		var res T
		return res
	}

	res := a[0]
	for i := range a {
		if a[i] > res {
			res = a[i]
		}
	}

	return res
}

// Min accepts any number of arguments of any type and returns min value
func Min[T Ordered](a ...T) T {
	if len(a) == 0 {
		var res T
		return res
	}

	res := a[0]
	for i := range a {
		if a[i] < res {
			res = a[i]
		}
	}

	return res
}

// IsEmptyString checks if the provided string is empty
func IsEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// NormalizeEOL will convert Windows (CRLF) and Mac (CR) EOLs to UNIX (LF)
func NormalizeEOL(input []byte) []byte {
	var right, left, pos int
	if right = bytes.IndexByte(input, '\r'); right == -1 {
		return input
	}
	length := len(input)
	tmp := make([]byte, length)

	// We know that left < length because otherwise right would be -1 from IndexByte.
	copy(tmp[pos:pos+right], input[left:left+right])
	pos += right
	tmp[pos] = '\n'
	left += right + 1
	pos++

	for left < length {
		if input[left] == '\n' {
			left++
		}

		right = bytes.IndexByte(input[left:], '\r')
		if right == -1 {
			copy(tmp[pos:], input[left:])
			pos += length - left
			break
		}
		copy(tmp[pos:pos+right], input[left:left+right])
		pos += right
		tmp[pos] = '\n'
		left += right + 1
		pos++
	}
	return tmp[:pos]
}

// ItemInSlice checks if given item a resides in given slice
func ItemInSlice[T Ordered](item T, slice []T) bool {
	for _, b := range slice {
		if b == item {
			return true
		}
	}
	return false
}

// RemoveItemsFromSlice removes all occurrences of items from slice
func RemoveItemsFromSlice[T Ordered](slice []T, items ...T) []T {
	res := make([]T, 0, cap(slice))

	for _, item := range slice {
		if !ItemInSlice(item, items) {
			res = append(res, item)
		}
	}

	return res
}

// SlicesIntersection returns a slice of common items of both given slices
func SlicesIntersection[T Ordered](slice1, slice2 []T) []T {
	meetMap := map[T]bool{}
	res := []T{}

	for _, value := range slice1 {
		meetMap[value] = true
	}

	for _, value := range slice2 {
		if meetMap[value] {
			res = append(res, value)
		}
	}

	return res
}

// Dedup return a slice of unique elements of any type
func Dedup[T Ordered](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})

	j := 0
	for i := 1; i < len(slice); i++ {
		if slice[j] == slice[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		slice[j] = slice[i]
	}

	return slice[:j+1]
}

// SumOfSlice returns sum of item in given array
func SumOfSlice[T Ordered](slice ...T) T {
	if len(slice) == 0 {
		var res T
		return res
	}

	sum := slice[0]
	for i := range slice {
		if i == 0 {
			continue
		}

		sum += slice[i]
	}

	return sum
}

func GetIPAddress(r *http.Request, trustedProxyIPHeader []string) string {
	var address string

	for _, proxyHeader := range trustedProxyIPHeader {
		header := r.Header.Get(proxyHeader)
		if header != "" {
			addresses := strings.Fields(header)
			if len(addresses) > 0 {
				address = strings.TrimRight(addresses[0], ",")
			}
		}

		if address != "" {
			return address
		}
	}

	if address == "" {
		address, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return address
}

func GetHostnameFromSiteURL(siteURL string) string {
	u, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

type RequestCache struct {
	Data []byte
	Date string
	Key  string
}

// Fetch JSON data from the notices server
// if skip is passed, does a fetch without touching the cache
func GetURLWithCache(url string, cache *RequestCache, skip bool) ([]byte, error) {
	// Build a GET Request, including optional If-None-Match header.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		cache.Data = nil
		return nil, err
	}
	if !skip && cache.Data != nil {
		req.Header.Add("If-None-Match", cache.Key)
		req.Header.Add("If-Modified-Since", cache.Date)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		cache.Data = nil
		return nil, err
	}
	defer resp.Body.Close()
	// No change from latest known Etag?
	if resp.StatusCode == http.StatusNotModified {
		return cache.Data, nil
	}

	if resp.StatusCode != 200 {
		cache.Data = nil
		return nil, errors.Errorf("Fetching notices failed with status code %d", resp.StatusCode)
	}

	cache.Data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		cache.Data = nil
		return nil, err
	}

	// If etags headers are missing, ignore.
	cache.Key = resp.Header.Get("ETag")
	cache.Date = resp.Header.Get("Date")
	return cache.Data, err
}

// Append tokens to passed baseUrl as query params
func AppendQueryParamsToURL(baseURL string, params map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return ""
	}
	for key, value := range params {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// GetFunctionName returns a string name of given function
//
// E.g
//
//	func hello() {}
//	name := GetFunctionName(hello)
//	fmt.Println(name) == "hello"
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
