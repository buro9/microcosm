// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cast"
	bp "github.com/spf13/hugo/bufferpool"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// Strips carriage returns from third-party / external processes (useful for Windows)
func normalizeExternalHelperLineFeeds(content []byte) []byte {
	return bytes.Replace(content, []byte("\r"), []byte(""), -1)
}

// FindAvailablePort returns an available and valid TCP port.
func FindAvailablePort() (*net.TCPAddr, error) {
	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer l.Close()
		addr := l.Addr()
		if a, ok := addr.(*net.TCPAddr); ok {
			return a, nil
		}
		return nil, fmt.Errorf("Unable to obtain a valid tcp port. %v", addr)
	}
	return nil, err
}

// InStringArray checks if a string is an element of a slice of strings
// and returns a boolean value.
func InStringArray(arr []string, el string) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

// GuessType attempts to guess the type of file from a given string.
func GuessType(in string) string {
	switch strings.ToLower(in) {
	case "md", "markdown", "mdown":
		return "markdown"
	case "asciidoc", "adoc", "ad":
		return "asciidoc"
	case "mmark":
		return "mmark"
	case "rst":
		return "rst"
	case "html", "htm":
		return "html"
	}

	return "unknown"
}

// FirstUpper returns a string with the first character as upper case.
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// UniqueStrings returns a new slice with any duplicates removed.
func UniqueStrings(s []string) []string {
	var unique []string
	set := map[string]interface{}{}
	for _, val := range s {
		if _, ok := set[val]; !ok {
			unique = append(unique, val)
			set[val] = val
		}
	}
	return unique
}

// ReaderToBytes takes an io.Reader argument, reads from it
// and returns bytes.
func ReaderToBytes(lines io.Reader) []byte {
	if lines == nil {
		return []byte{}
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	b.ReadFrom(lines)

	bc := make([]byte, b.Len(), b.Len())
	copy(bc, b.Bytes())
	return bc
}

// ToLowerMap makes all the keys in the given map lower cased and will do so
// recursively.
// Notes:
// * This will modify the map given.
// * Any nested map[interface{}]interface{} will be converted to map[string]interface{}.
func ToLowerMap(m map[string]interface{}) {
	for k, v := range m {
		switch v.(type) {
		case map[interface{}]interface{}:
			v = cast.ToStringMap(v)
			ToLowerMap(v.(map[string]interface{}))
		case map[string]interface{}:
			ToLowerMap(v.(map[string]interface{}))
		}

		lKey := strings.ToLower(k)
		if k != lKey {
			delete(m, k)
			m[lKey] = v
		}

	}
}

// ReaderToString is the same as ReaderToBytes, but returns a string.
func ReaderToString(lines io.Reader) string {
	if lines == nil {
		return ""
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	b.ReadFrom(lines)
	return b.String()
}

// ReaderContains reports whether subslice is within r.
func ReaderContains(r io.Reader, subslice []byte) bool {

	if r == nil || len(subslice) == 0 {
		return false
	}

	bufflen := len(subslice) * 4
	halflen := bufflen / 2
	buff := make([]byte, bufflen)
	var err error
	var n, i int

	for {
		i++
		if i == 1 {
			n, err = io.ReadAtLeast(r, buff[:halflen], halflen)
		} else {
			if i != 2 {
				// shift left to catch overlapping matches
				copy(buff[:], buff[halflen:])
			}
			n, err = io.ReadAtLeast(r, buff[halflen:], halflen)
		}

		if n > 0 && bytes.Contains(buff, subslice) {
			return true
		}

		if err != nil {
			break
		}
	}
	return false
}

// ThemeSet checks whether a theme is in use or not.
func ThemeSet() bool {
	return viper.GetString("theme") != ""
}

type logPrinter interface {
	// Println is the only common method that works in all of JWWs loggers.
	Println(a ...interface{})
}

// DistinctLogger ignores duplicate log statements.
type DistinctLogger struct {
	sync.RWMutex
	logger logPrinter
	m      map[string]bool
}

// Println will log the string returned from fmt.Sprintln given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Println(v ...interface{}) {
	// fmt.Sprint doesn't add space between string arguments
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.print(logStatement)
}

// Printf will log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
// Note: A newline is appended.
func (l *DistinctLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.print(logStatement)
}

func (l *DistinctLogger) print(logStatement string) {
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	if !l.m[logStatement] {
		l.logger.Println(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}

// NewDistinctErrorLogger creates a new DistinctLogger that logs ERRORs
func NewDistinctErrorLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.ERROR}
}

// NewDistinctFeedbackLogger creates a new DistinctLogger that can be used
// to give feedback to the user while not spamming with duplicates.
func NewDistinctFeedbackLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.FEEDBACK}
}

var (
	// DistinctErrorLog can be used to avoid spamming the logs with errors.
	DistinctErrorLog = NewDistinctErrorLogger()

	// DistinctFeedbackLog can be used to avoid spamming the logs with info messages.
	DistinctFeedbackLog = NewDistinctFeedbackLogger()
)

// InitLoggers sets up the global distinct loggers.
func InitLoggers() {
	DistinctErrorLog = NewDistinctErrorLogger()
}

// Deprecated informs about a deprecation, but only once for a given set of arguments' values.
// If the err flag is enabled, it logs as an ERROR (will exit with -1) and the text will
// point at the next Hugo release.
// The idea is two remove an item in two Hugo releases to give users and theme authors
// plenty of time to fix their templates.
func Deprecated(object, item, alternative string, err bool) {
	if err {
		DistinctErrorLog.Printf("%s's %s is deprecated and will be removed in Hugo %s. Use %s instead.", object, item, NextHugoReleaseVersion(), alternative)

	} else {
		// Make sure the users see this while avoiding build breakage. This will not lead to an os.Exit(-1)
		DistinctFeedbackLog.Printf("WARNING: %s's %s is deprecated and will be removed in a future release. Use %s instead.", object, item, alternative)
	}
}

// SliceToLower goes through the source slice and lowers all values.
func SliceToLower(s []string) []string {
	if s == nil {
		return nil
	}

	l := make([]string, len(s))
	for i, v := range s {
		l[i] = strings.ToLower(v)
	}

	return l
}

// Md5String takes a string and returns its MD5 hash.
func Md5String(f string) string {
	h := md5.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum([]byte{}))
}

// IsWhitespace determines if the given rune is whitespace.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// Seq creates a sequence of integers.
// It's named and used as GNU's seq.
// Examples:
// 3 => 1, 2, 3
// 1 2 4 => 1, 3
// -3 => -1, -2, -3
// 1 4 => 1, 2, 3, 4
// 1 -2 => 1, 0, -1, -2
func Seq(args ...interface{}) ([]int, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, errors.New("Seq, invalid number of args: 'first' 'increment' (optional) 'last' (optional)")
	}

	intArgs := cast.ToIntSlice(args)

	if len(intArgs) < 1 || len(intArgs) > 3 {
		return nil, errors.New("Invalid argument(s) to Seq")
	}

	var inc = 1
	var last int
	var first = intArgs[0]

	if len(intArgs) == 1 {
		last = first
		if last == 0 {
			return []int{}, nil
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(intArgs) == 2 {
		last = intArgs[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = intArgs[1]
		last = intArgs[2]
		if inc == 0 {
			return nil, errors.New("'increment' must not be 0")
		}
		if first < last && inc < 0 {
			return nil, errors.New("'increment' must be > 0")
		}
		if first > last && inc > 0 {
			return nil, errors.New("'increment' must be < 0")
		}
	}

	// sanity check
	if last < -100000 {
		return nil, errors.New("size of result exceeds limit")
	}
	size := ((last - first) / inc) + 1

	// sanity check
	if size <= 0 || size > 2000 {
		return nil, errors.New("size of result exceeds limit")
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq, nil
}

// DoArithmetic performs arithmetic operations (+,-,*,/) using reflection to
// determine the type of the two terms.
func DoArithmetic(a, b interface{}, op rune) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64
	var af, bf float64
	var au, bu uint64
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
			if ai >= 0 {
				au = uint64(ai)
				ai = 0
			} else {
				bi = int64(bu) // may overflow
				bu = 0
			}
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au = av.Uint()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
			if bi >= 0 {
				bu = uint64(bi)
				bi = 0
			} else {
				ai = int64(au) // may overflow
				au = 0
			}
		case reflect.Float32, reflect.Float64:
			af = float64(au) // may overflow
			au = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.String:
		as := av.String()
		if bv.Kind() == reflect.String && op == '+' {
			bs := bv.String()
			return as + bs, nil
		}
		return nil, errors.New("Can't apply the operator to the values")
	default:
		return nil, errors.New("Can't apply the operator to the values")
	}

	switch op {
	case '+':
		if ai != 0 || bi != 0 {
			return ai + bi, nil
		} else if af != 0 || bf != 0 {
			return af + bf, nil
		} else if au != 0 || bu != 0 {
			return au + bu, nil
		}
		return 0, nil
	case '-':
		if ai != 0 || bi != 0 {
			return ai - bi, nil
		} else if af != 0 || bf != 0 {
			return af - bf, nil
		} else if au != 0 || bu != 0 {
			return au - bu, nil
		}
		return 0, nil
	case '*':
		if ai != 0 || bi != 0 {
			return ai * bi, nil
		} else if af != 0 || bf != 0 {
			return af * bf, nil
		} else if au != 0 || bu != 0 {
			return au * bu, nil
		}
		return 0, nil
	case '/':
		if bi != 0 {
			return ai / bi, nil
		} else if bf != 0 {
			return af / bf, nil
		} else if bu != 0 {
			return au / bu, nil
		}
		return nil, errors.New("Can't divide the value by 0")
	default:
		return nil, errors.New("There is no such an operation")
	}
}

// NormalizeHugoFlags facilitates transitions of Hugo command-line flags,
// e.g. --baseUrl to --baseURL, --uglyUrls to --uglyURLs
func NormalizeHugoFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "baseUrl":
		name = "baseURL"
		break
	case "uglyUrls":
		name = "uglyURLs"
		break
	}
	return pflag.NormalizedName(name)
}

// DiffStringSlices returns the difference between two string slices.
// Useful in tests.
// See:
// http://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings-in-golang
func DiffStringSlices(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
