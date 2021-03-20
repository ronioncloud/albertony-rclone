package fs

// SizeSuffix is parsed by flag with K/M/G binary suffixes
import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SizeSuffix is an int64 with a friendly way of printing setting
type SizeSuffix int64

// Common multipliers for SizeSuffix
const (
	Byte SizeSuffix = 1 << (iota * 10)
	KibiByte
	MebiByte
	GibiByte
	TebiByte
	PebiByte
	ExbiByte
)

// Turn SizeSuffix into a string and a suffix
func (x SizeSuffix) string() (string, string) {
	scaled := float64(0)
	suffix := ""
	switch {
	case x < 0:
		return "off", ""
	case x == 0:
		return "0", ""
	case x < KibiByte:
		scaled = float64(x)
		suffix = ""
	case x < MebiByte:
		scaled = float64(x) / float64(KibiByte)
		suffix = "K"
	case x < GibiByte:
		scaled = float64(x) / float64(MebiByte)
		suffix = "M"
	case x < TebiByte:
		scaled = float64(x) / float64(GibiByte)
		suffix = "G"
	case x < PebiByte:
		scaled = float64(x) / float64(TebiByte)
		suffix = "T"
	case x < ExbiByte:
		scaled = float64(x) / float64(PebiByte)
		suffix = "P"
	default:
		scaled = float64(x) / float64(ExbiByte)
		suffix = "E"
	}
	if math.Floor(scaled) == scaled {
		return fmt.Sprintf("%.0f", scaled), suffix
	}
	return fmt.Sprintf("%.3f", scaled), suffix
}

// String turns SizeSuffix into a string
func (x SizeSuffix) String() string {
	val, suffix := x.string()
	return val + suffix
}

// Unit turns SizeSuffix into a string with a unit
func (x SizeSuffix) Unit(unit string) string {
	val, suffix := x.string()
	if val == "off" {
		return val
	}
	var suffixUnit string
	if suffix != "" && unit != "" {
		suffixUnit = suffix + "i" + unit
	} else {
		suffixUnit = suffix + unit
	}
	return val + " " + suffixUnit
}

// Set a SizeSuffix
func (x *SizeSuffix) Set(s string) error {
	if len(s) == 0 {
		return errors.New("empty string")
	}
	if strings.ToLower(s) == "off" {
		*x = -1
		return nil
	}
	suffix := s[len(s)-1]
	suffixLen := 1
	var multiplier float64
	switch suffix {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		suffixLen = 0
		multiplier = float64(KibiByte)
	case 'b', 'B':
		multiplier = float64(Byte)
	case 'k', 'K':
		multiplier = float64(KibiByte)
	case 'm', 'M':
		multiplier = float64(MebiByte)
	case 'g', 'G':
		multiplier = float64(GibiByte)
	case 't', 'T':
		multiplier = float64(TebiByte)
	case 'p', 'P':
		multiplier = float64(PebiByte)
	case 'e', 'E':
		multiplier = float64(ExbiByte)
	default:
		return errors.Errorf("bad suffix %q", suffix)
	}
	s = s[:len(s)-suffixLen]
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	if value < 0 {
		return errors.Errorf("size can't be negative %q", s)
	}
	value *= multiplier
	*x = SizeSuffix(value)
	return nil
}

// Type of the value
func (x *SizeSuffix) Type() string {
	return "SizeSuffix"
}

// Scan implements the fmt.Scanner interface
func (x *SizeSuffix) Scan(s fmt.ScanState, ch rune) error {
	token, err := s.Token(true, nil)
	if err != nil {
		return err
	}
	return x.Set(string(token))
}

// SizeSuffixList is a slice SizeSuffix values
type SizeSuffixList []SizeSuffix

func (l SizeSuffixList) Len() int           { return len(l) }
func (l SizeSuffixList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SizeSuffixList) Less(i, j int) bool { return l[i] < l[j] }

// Sort sorts the list
func (l SizeSuffixList) Sort() {
	sort.Sort(l)
}

// UnmarshalJSONFlag unmarshals a JSON input for a flag. If the input
// is a string then it calls the Set method on the flag otherwise it
// calls the setInt function with a parsed int64.
func UnmarshalJSONFlag(in []byte, x interface{ Set(string) error }, setInt func(int64) error) error {
	// Try to parse as string first
	var s string
	err := json.Unmarshal(in, &s)
	if err == nil {
		return x.Set(s)
	}
	// If that fails parse as integer
	var i int64
	err = json.Unmarshal(in, &i)
	if err != nil {
		return err
	}
	return setInt(i)
}

// UnmarshalJSON makes sure the value can be parsed as a string or integer in JSON
func (x *SizeSuffix) UnmarshalJSON(in []byte) error {
	return UnmarshalJSONFlag(in, x, func(i int64) error {
		*x = SizeSuffix(i)
		return nil
	})
}
