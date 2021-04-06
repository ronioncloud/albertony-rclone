package fs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Interface which flags must satisfy - only defined for _test.go
// since we don't want to pull in pflag here
type flagger interface {
	pflag.Value
	json.Unmarshaler
}

// Check it satisfies the interface
var _ flagger = (*SizeSuffix)(nil)

func TestSizeSuffixString(t *testing.T) {
	for _, test := range []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{102, "102"},
		{1024, "1K"},
		{1024 * 1024, "1M"},
		{1024 * 1024 * 1024, "1G"},
		{10 * 1024 * 1024 * 1024, "10G"},
		{10.1 * 1024 * 1024 * 1024, "10.100G"},
		{-1, "off"},
		{-100, "off"},
	} {
		ss := SizeSuffix(test.in)
		got := ss.String()
		assert.Equal(t, test.want, got)
	}
}

func TestSizeSuffixUnit(t *testing.T) {
	for _, test := range []struct {
		in   float64
		want string
	}{
		{0, "0 Bytes"},
		{102, "102 Bytes"},
		{1024, "1 KiBytes"},
		{1024 * 1024, "1 MiBytes"},
		{1024 * 1024 * 1024, "1 GiBytes"},
		{10 * 1024 * 1024 * 1024, "10 GiBytes"},
		{10.1 * 1024 * 1024 * 1024, "10.100 GiBytes"},
		{10 * 1024 * 1024 * 1024 * 1024, "10 TiBytes"},
		{10 * 1024 * 1024 * 1024 * 1024 * 1024, "10 PiBytes"},
		{1 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1 EiBytes"},
		{-1, "off"},
		{-100, "off"},
	} {
		ss := SizeSuffix(test.in)
		got := ss.Unit("Bytes")
		assert.Equal(t, test.want, got)
	}
}

func TestSizeSuffixSet(t *testing.T) {
	for _, test := range []struct {
		in   string
		want int64
		err  bool
	}{
		{"0", 0, false},
		{"1b", 1, false},
		{"102B", 102, false},
		{"0.1k", 102, false},
		{"0.1", 102, false},
		{"1K", 1024, false},
		{"1", 1024, false},
		{"2.5", 1024 * 2.5, false},
		{"1M", 1024 * 1024, false},
		{"1.g", 1024 * 1024 * 1024, false},
		{"10G", 10 * 1024 * 1024 * 1024, false},
		{"10T", 10 * 1024 * 1024 * 1024 * 1024, false},
		{"10P", 10 * 1024 * 1024 * 1024 * 1024 * 1024, false},
		{"off", -1, false},
		{"OFF", -1, false},
		{"", 0, true},
		{"1q", 0, true},
		{"1.q", 0, true},
		{"1q", 0, true},
		{"-1K", 0, true},
	} {
		ss := SizeSuffix(0)
		err := ss.Set(test.in)
		if test.err {
			require.Error(t, err, test.in)
		} else {
			require.NoError(t, err, test.in)
		}
		assert.Equal(t, test.want, int64(ss))
	}
}

func TestSizeSuffixScan(t *testing.T) {
	var v SizeSuffix
	n, err := fmt.Sscan(" 17M ", &v)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, SizeSuffix(17<<20), v)
}

func TestSizeSuffixUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		in   string
		want int64
		err  bool
	}{
		{`"0"`, 0, false},
		{`"102B"`, 102, false},
		{`"1K"`, 1024, false},
		{`"2.5"`, 1024 * 2.5, false},
		{`"1M"`, 1024 * 1024, false},
		{`"1.g"`, 1024 * 1024 * 1024, false},
		{`"10G"`, 10 * 1024 * 1024 * 1024, false},
		{`"off"`, -1, false},
		{`""`, 0, true},
		{`"1q"`, 0, true},
		{`"-1K"`, 0, true},
		{`0`, 0, false},
		{`102`, 102, false},
		{`1024`, 1024, false},
		{`1000000000`, 1000000000, false},
		{`1.1.1`, 0, true},
	} {
		var ss SizeSuffix
		err := json.Unmarshal([]byte(test.in), &ss)
		if test.err {
			require.Error(t, err, test.in)
		} else {
			require.NoError(t, err, test.in)
		}
		assert.Equal(t, test.want, int64(ss))
	}
}
