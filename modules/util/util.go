package util

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
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

	cache.Data, err = io.ReadAll(resp.Body)
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
func GetFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func NewSet[T cmp.Ordered](items ...T) *AnySet[T] {
	res := &AnySet[T]{
		meetMap: make(map[T]struct{}),
	}
	res.Add(items...)

	return res
}

// AnySet makes sure there are no duplicate in its values.
type AnySet[T cmp.Ordered] struct {
	values  AnyArray[T]
	meetMap map[T]struct{}
}

func (s *AnySet[T]) Add(items ...T) {
	if s.meetMap == nil {
		s.meetMap = make(map[T]struct{})
	}
	for _, item := range items {
		if _, ok := s.meetMap[item]; !ok {
			s.values = append(s.values, item)
			s.meetMap[item] = struct{}{}
		}
	}
}

func (s *AnySet[T]) Values() AnyArray[T] {
	return s.values
}

// AnyArray if a generic slice with a set of member methods that can be chained
type AnyArray[T cmp.Ordered] []T

// Remove removes input from the array
func (a AnyArray[T]) Remove(item T) AnyArray[T] {
	var res = make(AnyArray[T], 0, cap(a))
	for _, it := range a {
		if it != item {
			res = append(res, it)
		}
	}

	return res
}

// Dedup keeps each item in current array appears once only.
func (a AnyArray[T]) Dedup() AnyArray[T] {
	meetMap := map[T]struct{}{}
	res := AnyArray[T]{}

	for _, item := range a {
		_, ok := meetMap[item]
		if !ok {
			res = append(res, item)
			meetMap[item] = struct{}{}
		}
	}

	return res
}

// InterSection returns items that appear in both current array and given others
func (s AnyArray[T]) InterSection(others []T) AnyArray[T] {
	var res AnyArray[T]
	meetMap := map[T]struct{}{}

	for _, item := range s {
		meetMap[item] = struct{}{}
	}

	for _, item := range others {
		_, ok := meetMap[item]
		if ok {
			res = append(res, item)
		}
	}

	return res
}

// Sum adds up items in current array and returns the result
func (a AnyArray[T]) Sum() T {
	var res T
	for _, item := range a {
		res += item
	}
	return res
}

func (a AnyArray[T]) Len() int {
	return len(a)
}

// Map loops through current string slice and applies mapFunc to each index-item pair
//
// E.g
//
//	AnyArray{"a", "b", "c"}.Map(func(_ int, s string) string { return s + s })
func (a AnyArray[T]) Map(fn func(index int, item T) T) AnyArray[T] {
	res := make([]T, len(a), cap(a))

	for idx, item := range a {
		res[idx] = fn(idx, item)
	}

	return res
}

// check if array of strings contains given input
func (sa AnyArray[T]) Contains(input T) bool {
	for _, item := range sa {
		if item == input {
			return true
		}
	}
	return false
}

func (sa AnyArray[T]) ContainsAny(input ...T) bool {
	for _, item := range input {
		if sa.Contains(item) {
			return true
		}
	}
	return false
}

// Equals checks if two arrays of strings have same length and contains the same elements at each index
func (sa AnyArray[T]) Equals(input []T) bool {
	if len(sa) != len(input) {
		return false
	}

	for idx, item := range sa {
		if item != input[idx] {
			return false
		}
	}

	return true
}

// HasDuplicates checks if there are duplicates in current array
func (sa AnyArray[T]) HasDuplicates() bool {
	meetMap := map[T]struct{}{}

	for _, item := range sa {
		_, met := meetMap[item]
		if met {
			return true
		}
		meetMap[item] = struct{}{}
	}

	return false
}

// Join
func (sa AnyArray[T]) Join(sep string) string {
	var builder strings.Builder

	for i, item := range sa {
		if i == len(sa)-1 {
			sep = ""
		}
		builder.WriteString(fmt.Sprintf("%v%s", item, sep))
	}

	return builder.String()
}
