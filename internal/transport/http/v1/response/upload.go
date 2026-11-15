package response

import "time"

type CreateUploadResponse struct {
	TileID     string            `json:"tile_id"`
	UploadURL  string            `json:"upload_url"`
	FormFields map[string]string `json:"form_fields"`
	ExpiresAt  time.Time         `json:"expires_at"`
}
