package modality

import "strings"

type Type string

const (
	Unknown    Type = "unknown"
	Text       Type = "text"
	Image      Type = "image"
	Audio      Type = "audio"
	Video      Type = "video"
	Document   Type = "document"
	Multimodal Type = "multimodal"
)

func Normalize(raw string) Type {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "text":
		return Text
	case "image", "vision":
		return Image
	case "audio", "speech":
		return Audio
	case "video":
		return Video
	case "document", "pdf":
		return Document
	case "multimodal", "multi":
		return Multimodal
	default:
		return Unknown
	}
}

func (t Type) IsKnown() bool {
	return t != "" && t != Unknown
}

func (t Type) String() string {
	if t == "" {
		return string(Unknown)
	}
	return string(t)
}
