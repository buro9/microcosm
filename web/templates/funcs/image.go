package funcs

import (
	"strings"
)

// isImage returns true if the path supplied has a file extension that matches
// a known image type according to
// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types#common_image_file_types
func isImage(u string) bool {
	u = strings.ToLower(u)
	switch {
	case strings.HasSuffix(u, `.apng`):
		return true

	case strings.HasSuffix(u, `.avif`):
		return true

	case strings.HasSuffix(u, `.gif`):
		return true

	case strings.HasSuffix(u, `.jpg`):
		return true

	case strings.HasSuffix(u, `.jpeg`):
		return true

	case strings.HasSuffix(u, `.jpe`):
		return true

	case strings.HasSuffix(u, `.jif`):
		return true

	case strings.HasSuffix(u, `.jfif`):
		return true

	case strings.HasSuffix(u, `.png`):
		return true

	case strings.HasSuffix(u, `.svg`):
		return true

	case strings.HasSuffix(u, `.webp`):
		return true

	default:
		return false
	}
}
