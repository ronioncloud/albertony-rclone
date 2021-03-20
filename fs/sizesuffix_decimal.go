package fs

// SizeSuffixDecimal is parsed by flag with k/M/G decimal suffixes
import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SizeSuffixDecimal is an int64 with a friendly way of printing setting
type SizeSuffixDecimal int64

// Common multipliers for SizeSuffix
const (
	DecimalByte SizeSuffixDecimal = 1
	KiloByte                      = 1000 * DecimalByte
	MegaByte                      = 1000 * KiloByte
	GigaByte                      = 1000 * MegaByte
	TeraByte                      = 1000 * GigaByte
	PetaByte                      = 1000 * TeraByte
	ExaByte                       = 1000 * PetaByte
)

// Turn SizeSuffixDecimal into a string and a suffix
func (x SizeSuffixDecimal) string() (string, string) {
	scaled := float64(0)
	suffix := ""
	switch {
	case x < 0:
		return "off", ""
	case x == 0:
		return "0", ""
	case x < KiloByte:
		scaled = float64(x)
		suffix = ""
	case x < MegaByte:
		scaled = float64(x) / float64(KiloByte)
		suffix = "k"
	case x < GigaByte:
		scaled = float64(x) / float64(MegaByte)
		suffix = "M"
	case x < TeraByte:
		scaled = float64(x) / float64(GigaByte)
		suffix = "G"
	case x < PetaByte:
		scaled = float64(x) / float64(TeraByte)
		suffix = "T"
	case x < ExaByte:
		scaled = float64(x) / float64(PetaByte)
		suffix = "P"
	default:
		scaled = float64(x) / float64(ExaByte)
		suffix = "E"
	}
	if math.Floor(scaled) == scaled {
		return fmt.Sprintf("%.0f", scaled), suffix
	}
	return fmt.Sprintf("%.3f", scaled), suffix
}

// String turns SizeSuffixDecimal into a string
func (x SizeSuffixDecimal) String() string {
	val, suffix := x.string()
	return val + suffix
}

// Unit turns SizeSuffixDecimal into a string with a unit
func (x SizeSuffixDecimal) Unit(unit string) string {
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

// Set a SizeSuffixDecimal
func (x *SizeSuffixDecimal) Set(s string) error {
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
		multiplier = float64(KiloByte)
	case 'b', 'B':
		multiplier = float64(Byte)
	case 'k', 'K':
		multiplier = float64(KiloByte)
	case 'm', 'M':
		multiplier = float64(MegaByte)
	case 'g', 'G':
		multiplier = float64(GigaByte)
	case 't', 'T':
		multiplier = float64(TeraByte)
	case 'p', 'P':
		multiplier = float64(PetaByte)
	case 'e', 'E':
		multiplier = float64(ExaByte)
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
	*x = SizeSuffixDecimal(value)
	return nil
}

// Type of the value
func (x *SizeSuffixDecimal) Type() string {
	return "SizeSuffixDecimal"
}

// Scan implements the fmt.Scanner interface
func (x *SizeSuffixDecimal) Scan(s fmt.ScanState, ch rune) error {
	token, err := s.Token(true, nil)
	if err != nil {
		return err
	}
	return x.Set(string(token))
}

// SizeSuffixDecimalList is a slice SizeSuffixDecimal values
type SizeSuffixDecimalList []SizeSuffixDecimal

func (l SizeSuffixDecimalList) Len() int           { return len(l) }
func (l SizeSuffixDecimalList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SizeSuffixDecimalList) Less(i, j int) bool { return l[i] < l[j] }

// Sort sorts the list
func (l SizeSuffixDecimalList) Sort() {
	sort.Sort(l)
}

// UnmarshalJSON makes sure the value can be parsed as a string or integer in JSON
func (x *SizeSuffixDecimal) UnmarshalJSON(in []byte) error {
	return UnmarshalJSONFlag(in, x, func(i int64) error {
		*x = SizeSuffixDecimal(i)
		return nil
	})
}
