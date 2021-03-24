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
		suffix = "Ki"
	case x < GibiByte:
		scaled = float64(x) / float64(MebiByte)
		suffix = "Mi"
	case x < TebiByte:
		scaled = float64(x) / float64(GibiByte)
		suffix = "Gi"
	case x < PebiByte:
		scaled = float64(x) / float64(TebiByte)
		suffix = "Ti"
	case x < ExbiByte:
		scaled = float64(x) / float64(PebiByte)
		suffix = "Pi"
	default:
		scaled = float64(x) / float64(ExbiByte)
		suffix = "Ei"
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
func (x SizeSuffix) unit(unit string) string {
	val, suffix := x.string()
	if val == "off" {
		return val
	}
	var suffixUnit string
	if suffix != "" && unit != "" {
		suffixUnit = suffix + unit
	} else {
		suffixUnit = suffix + unit
	}
	return val + " " + suffixUnit
}

// BitUnit turns SizeSuffix into a string with bit unit
func (x SizeSuffix) BitUnit() string {
	return x.unit("bit")
}

// BitRateUnit turns SizeSuffix into a string with bit rate unit
func (x SizeSuffix) BitRateUnit() string {
	return x.unit("bit/s")
}

// ByteUnit turns SizeSuffix into a string with byte unit
func (x SizeSuffix) ByteUnit() string {
	return x.unit("Byte")
}

// ByteRateUnit turns SizeSuffix into a string with byte rate unit
func (x SizeSuffix) ByteRateUnit() string {
	return x.unit("Byte/s")
}

// ByteShortUnit turns SizeSuffix into a string with byte unit short form
func (x SizeSuffix) ByteShortUnit() string {
	return x.unit("B")
}

// ByteRateShortUnit turns SizeSuffix into a string with byte rate unit short form
func (x SizeSuffix) ByteRateShortUnit() string {
	return x.unit("B/s")
}

func (x *SizeSuffix) symbolMultiplier(s byte) (found bool, multiplier float64) {
	switch s {
	case 'k', 'K':
		return true, float64(KibiByte)
	case 'm', 'M':
		return true, float64(MebiByte)
	case 'g', 'G':
		return true, float64(GibiByte)
	case 't', 'T':
		return true, float64(TebiByte)
	case 'p', 'P':
		return true, float64(PebiByte)
	case 'e', 'E':
		return true, float64(ExbiByte)
	default:
		return false, float64(Byte)
	}
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
	unitPrefix := false
	var multiplier float64
	switch suffix {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		suffixLen = 0
		multiplier = float64(KibiByte)
	case 'b', 'B':
		if len(s) > 2 && s[len(s)-2] == 'i' {
			suffix = s[len(s)-3]
			suffixLen = 3
			if unitPrefix, multiplier = x.symbolMultiplier(suffix); !unitPrefix {
				return errors.Errorf("bad suffix %q", suffix)
			}
			// TODO: Support SI form MB, treat it equivalent to MiB, or reserve it for the SizeSuffixDecimal only?
			//} else if len(s) > 1 {
			//	suffix = s[len(s)-2]
			//	if unitPrefix, multiplier = x.suffixUnitPrefix(suffix); unitPrefix {
			//		suffixLen = 2
			//	}
		} else {
			multiplier = float64(Byte)
		}
	case 'i', 'I':
		if len(s) > 1 {
			suffix = s[len(s)-2]
			suffixLen = 2
			unitPrefix, multiplier = x.symbolMultiplier(suffix)
		}
		if !unitPrefix {
			return errors.Errorf("bad suffix %q", suffix)
		}
	default:
		if unitPrefix, multiplier = x.symbolMultiplier(suffix); !unitPrefix {
			return errors.Errorf("bad suffix %q", suffix)
		}
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
