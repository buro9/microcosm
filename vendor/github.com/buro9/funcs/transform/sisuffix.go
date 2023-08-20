package transform

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// SiSuffix takes a number and returns an SI symbol as a suffix and the number
// rounded to 2 decimal places.
//
// If the value provided is not a number, an empty string is returned
func SiSuffix(value interface{}) string {
	var f float64
	switch v := value.(type) {
	case float32:
		f = float64(v)
	case float64:
		f = v
	case int:
		f = float64(v)
	case int32:
		f = float64(v)
	case int64:
		f = float64(v)
	default:
		return ""
	}

	return si(f)
}

// The code below comes from https://github.com/dustin/go-humanize and is
// licensed under the MIT license as follows. It it reproduced here so that the
// output could be varied to the needs of this software.
//
// Copyright (c) 2005-2008  Dustin Sallings <dustin@spy.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// <http://www.opensource.org/licenses/mit-license.php>

func stripTrailingZeros(s string) string {
	offset := len(s) - 1
	for offset > 0 {
		if s[offset] == '.' {
			offset--
			break
		}
		if s[offset] != '0' {
			break
		}
		offset--
	}
	return s[:offset+1]
}

var siPrefixTable = map[float64]string{
	-24: "y", // yocto
	-21: "z", // zepto
	-18: "a", // atto
	-15: "f", // femto
	-12: "p", // pico
	-9:  "n", // nano
	-6:  "µ", // micro
	-3:  "m", // milli
	0:   "",
	3:   "k", // kilo
	6:   "M", // mega
	9:   "G", // giga
	12:  "T", // tera
	15:  "P", // peta
	18:  "E", // exa
	21:  "Z", // zetta
	24:  "Y", // yotta
}

var revSIPrefixTable = revfmap(siPrefixTable)

// revfmap reverses the map and precomputes the power multiplier
func revfmap(in map[float64]string) map[string]float64 {
	rv := map[string]float64{}
	for k, v := range in {
		rv[v] = math.Pow(10, k)
	}
	return rv
}

var riParseRegex *regexp.Regexp

func init() {
	ri := `^([\-0-9.]+)\s?([`
	for _, v := range siPrefixTable {
		ri += v
	}
	ri += `]?)(.*)`

	riParseRegex = regexp.MustCompile(ri)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

// ComputeSI finds the most appropriate SI prefix for the given number
// and returns the prefix along with the value adjusted to be within
// that prefix.
//
// See also: SI, ParseSI.
//
// e.g. ComputeSI(2.2345e-12) -> (2.2345, "p")
func computeSI(input float64) (float64, string) {
	if input == 0 {
		return 0, ""
	}
	mag := math.Abs(input)
	exponent := math.Floor(logn(mag, 10))
	exponent = math.Floor(exponent/3) * 3

	value := mag / math.Pow(10, exponent)

	// Handle special case where value is exactly 1000.0
	// Should return 1 M instead of 1000 k
	if value == 1000.0 {
		exponent += 3
		value = mag / math.Pow(10, exponent)
	}

	value = math.Copysign(value, input)

	prefix := siPrefixTable[exponent]
	return value, prefix
}

// SI returns a string with default formatting.
//
// SI uses Ftoa to format float value, removing trailing zeros.
func si(input float64) string {
	value, prefix := computeSI(input)
	n := fmt.Sprintf("%v", value)
	if strings.Index(n, ".") > -1 && len(n)-1-strings.Index(n, ".") > 2 {
		n = n[:strings.Index(n, ".")+3]
	}
	return stripTrailingZeros(n) + prefix
}
