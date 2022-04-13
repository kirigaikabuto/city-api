package common

var (
	imageTypes = []string{
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/tiff",
		"image/vnd.microsoft.icon",
		"image/x-icon",
		"image/vnd.djvu",
		"image/svg+xml",
	}
	videoTypes = []string{
		"video/mpeg",
		"video/mp4",
		"video/webm",
	}
)

func IsImage(contentType string) bool {
	for _, v := range imageTypes {
		if v == contentType {
			return true
		}
	}
	return false
}

func IsVideo(contentType string) bool {
	for _, v := range videoTypes {
		if v == contentType {
			return true
		}
	}
	return false
}
