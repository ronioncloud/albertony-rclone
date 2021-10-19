package fs

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"path"
	"regexp"
	"strings"
)

var mimeTypeDefinitionRegexp *regexp.Regexp

func init() {
	mimeTypeDefinitionRegexp = regexp.MustCompile(`(?:[^\s:,]+)+`)
	if err := initDefaultMimeTypes(); err != nil {
		panic(err)
	}
}

// initDefaultMimeTypes adds a minimal number of default mime types
//
// Main purpose is to augment go's built in types for environments which
// don't have access to a mime.types file (e.g. Termux on android).
func initDefaultMimeTypes() error {
	for _, t := range []struct {
		mimeType   string
		extensions []string
	}{
		{"audio/flac", []string{".flac"}},
		{"audio/mpeg", []string{".mpga", ".mpega", ".mp2", ".mp3", ".m4a"}},
		{"audio/ogg", []string{".oga", ".ogg", ".opus", ".spx"}},
		{"audio/x-wav", []string{".wav"}},
		{"image/tiff", []string{".tiff", ".tif"}},
		{"video/dv", []string{".dif", ".dv"}},
		{"video/fli", []string{".fli"}},
		{"video/mpeg", []string{".mpeg", ".mpg", ".mpe"}},
		{"video/MP2T", []string{".ts"}},
		{"video/mp4", []string{".mp4"}},
		{"video/quicktime", []string{".qt,.mov"}},
		{"video/ogg", []string{".ogv"}},
		{"video/webm", []string{".webm"}},
		{"video/x-msvideo", []string{".avi"}},
		{"video/x-matroska", []string{".mpv", ".mkv"}},
		{"text/srt", []string{".srt"}},
	} {
		if err := AddExtensionsForMimeType(t.mimeType, t.extensions, false); err != nil {
			return err
		}
	}
	return nil
}

// ImportMimeTypeFile imports mime type definitions from text file
//
// File format is similar to Linux mime.types, but much simpler:
// - Empty/whitespace-only lines are skipped
// - Comment lines, starting with the # character, are skipped
// - Mime definitions start with the mime media type followed by a series of file extensions
// - Separators can be space, ':' or ','
// - File extensions can be given with or without leading '.' (mime.types file format is without)
//
// Caller decides if any existing extensions associated with the types
// should be replaced or not.
func ImportMimeTypeFile(filePath string, replace bool) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(bytes.NewBuffer(fileBytes))
	var line string
	for err != io.EOF {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line = strings.TrimSpace(line)
		if line != "" && line[0] != '#' {
			if err = ImportMimeTypeString(line, replace); err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportMimeTypeString imports a mime type definition from string
//
// String format is media type followed by a series of file extensions, where:
// - Separators can be space, ':' or ','
// - Extensions can be given with or without leading '.'
// Examples:
// - "image/jpeg:jpg"
// - "image/jpeg    .jpeg .jpe .jpg"
//
// Caller decides if any existing extensions associated with the types
// should be replaced or not.
func ImportMimeTypeString(definition string, replace bool) error {
	matches := mimeTypeDefinitionRegexp.FindAllStringSubmatch(definition, -1)
	if len(matches) < 2 {
		return fmt.Errorf("invalid format in mime type definition %q", definition)
	}
	typ := matches[0][0]
	for _, submatches := range matches[1:] {
		ext := submatches[0]
		if ext[0] != '.' {
			ext = "." + ext
		}
		if err := SetMimeTypeForExtension(ext, typ, replace); err != nil {
			return err
		}
	}
	return nil
}

// AddExtensionsForMimeType sets extensions associated with a mime type
//
// Caller decides if any existing mime types associated with the
// extensions should be replaced or not.
func AddExtensionsForMimeType(mimeType string, extensions []string, replace bool) error {
	for _, extension := range extensions {
		if err := SetMimeTypeForExtension(extension, mimeType, replace); err != nil {
			return err
		}
	}
	return nil
}

// SetMimeTypeForExtension sets the mime type associated with an extension
//
// Caller decides if any existing extension associated with the type should
// be replaced or not.
func SetMimeTypeForExtension(extension, mimeType string, replace bool) error {
	if replace || mime.TypeByExtension(extension) == "" {
		if err := mime.AddExtensionType(extension, mimeType); err != nil {
			return err
		}
	}
	return nil
}

// MimeTypeFromName returns a guess at the mime type from the name
func MimeTypeFromName(remote string) (mimeType string) {
	mimeType = mime.TypeByExtension(path.Ext(remote))
	if !strings.ContainsRune(mimeType, '/') {
		mimeType = "application/octet-stream"
	}
	return mimeType
}

// MimeType returns the MimeType from the object, either by calling
// the MimeTyper interface or using MimeTypeFromName
func MimeType(ctx context.Context, o ObjectInfo) (mimeType string) {
	// Read the MimeType from the optional interface if available
	if do, ok := o.(MimeTyper); ok {
		mimeType = do.MimeType(ctx)
		// Debugf(o, "Read MimeType as %q", mimeType)
		if mimeType != "" {
			return mimeType
		}
	}
	return MimeTypeFromName(o.Remote())
}

// MimeTypeDirEntry returns the MimeType of a DirEntry
//
// It returns "inode/directory" for directories, or uses
// MimeType(Object)
func MimeTypeDirEntry(ctx context.Context, item DirEntry) string {
	switch x := item.(type) {
	case Object:
		return MimeType(ctx, x)
	case Directory:
		return "inode/directory"
	}
	return ""
}
