package models

import (
	"time"

	"github.com/google/uuid"
)

type UploadURLData struct {
	TileID               uuid.UUID
	ContentType          string
	AcceptedContentTypes []string
	Tasks                []ProcessingTask
	ExternalTasks        []string
	MaxSize              uint64
	MaxDuration          *uint32
	TTL                  *uint32
}

type UploadURL struct {
	FileID     string
	UploadURL  string
	FormFields map[string]string
	ExpiresAt  time.Time
}

type CreateUpload struct {
	FileType string
	OwnerID  uuid.UUID
}

type CreatedUpload struct {
	TileID     uuid.UUID
	FileID     string
	UploadURL  string
	FormFields map[string]string
	ExpiresAt  time.Time
}
