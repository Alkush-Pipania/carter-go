package utils

// AllowedContentTypes maps MIME types to source types
var AllowedContentTypes = map[string]string{
	"application/pdf":               "pdf",
	"application/vnd.ms-powerpoint": "ppt",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "ppt",
	"application/msword": "doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "doc",
}

// IsAllowedContentType checks if the content type is allowed for upload
func IsAllowedContentType(contentType string) bool {
	_, ok := AllowedContentTypes[contentType]
	return ok
}

// GetSourceTypeFromContentType returns the source type for a content type
func GetSourceTypeFromContentType(contentType string) string {
	return AllowedContentTypes[contentType]
}
