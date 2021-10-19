package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetMimeTypeForExtension(t *testing.T) {
	mimeType := "testing/aaa"
	extension := "aa1"
	err := SetMimeTypeForExtension(extension, mimeType, true)
	require.Error(t, err) // // "mime: missing leading dot"

	mimeType = "testing/aaa"
	extension = ".aa1"
	err = SetMimeTypeForExtension(extension, mimeType, true)
	require.NoError(t, err)
	require.Equal(t, mimeType, MimeTypeFromName("testing"+extension))

	mimeType = "testing/aaa"
	extension = ".aa2"
	err = SetMimeTypeForExtension(extension, mimeType, false)
	require.NoError(t, err)
	require.Equal(t, mimeType, MimeTypeFromName("testing"+extension))

	mimeType = "testing/bbb"
	extension = ".aa2"
	err = SetMimeTypeForExtension(extension, mimeType, false)
	require.NoError(t, err)
	require.Equal(t, "testing/aaa", MimeTypeFromName("testing"+extension))
}

func TestImportMimeTypeString(t *testing.T) {
	for _, test := range []struct {
		in         string
		mimeType   string
		extensions []string
	}{
		{in: "testing/aaa:aa1,aa2,aa3", mimeType: "testing/aaa", extensions: []string{".aa1", ".aa2", ".aa3"}},
		{in: "testing/bbb bb1,bb2,bb3", mimeType: "testing/bbb", extensions: []string{".bb1", ".bb2", ".bb3"}},
		{in: "testing/ccc:.cc1,.cc2,.cc3", mimeType: "testing/ccc", extensions: []string{".cc1", ".cc2", ".cc3"}},
		{in: "testing/ddd dd1 dd2 dd3", mimeType: "testing/ddd", extensions: []string{".dd1", ".dd2", ".dd3"}},
		{in: "testing/eee .ee1 .ee2 .ee3", mimeType: "testing/eee", extensions: []string{".ee1", ".ee2", ".ee3"}},
		{in: "testing/fff:ff1:ff2:ff3", mimeType: "testing/fff", extensions: []string{".ff1", ".ff2", ".ff3"}},
		{in: "testing/ggg:.gg1:.gg2:.gg3", mimeType: "testing/ggg", extensions: []string{".gg1", ".gg2", ".gg3"}},
		{in: "testing/hhh ,:hh1:,:hh2:,:hh3:,", mimeType: "testing/hhh", extensions: []string{".hh1", ".hh2", ".hh3"}},
		{in: "testing ii1", mimeType: "application/octet-stream", extensions: []string{".ii1"}}, // Will be successfully parsed as mime type "testing", but MimeTypeFromName returns "application/octet-stream" for mime types without "/"

		{in: "testing", mimeType: "", extensions: nil},
		{in: "ab:c", mimeType: "application/octet-stream", extensions: []string{".c"}},  // Parsed mime type "ab" without "/" again
		{in: "a:b/c", mimeType: "application/octet-stream", extensions: []string{".c"}}, // Parsed mime type "a" without "/" again
		{in: "a/:c", mimeType: "", extensions: nil},
		{in: "a/b:c", mimeType: "a/b", extensions: []string{".c"}},

		{in: "testing/jpeg:jpg", mimeType: "testing/jpeg", extensions: []string{".jpg"}},
		{in: "testing/mpeg    .mpga .mpega .mp2 .mp3 .m4a", mimeType: "testing/mpeg", extensions: []string{".mpga", ".mpega", ".mp2", ".mp3", ".m4a"}},
	} {
		err := ImportMimeTypeString(test.in, true)
		if test.mimeType != "" {
			require.NoError(t, err)
			for _, ext := range test.extensions {
				assert.Equal(t, test.mimeType, MimeTypeFromName("testing"+ext))
			}
		} else {
			require.Error(t, err)
		}
	}
}
