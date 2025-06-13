package schemas

type S3UploadRequest struct {
	Folder   string `json:"folder"             validate:"required"`
	FileName string `json:"fileName"           validate:"required"`
	EntityID string `json:"entityId,omitempty"` // Optional: ID of the entity this image belongs to
}

type S3DownloadRequest struct {
	Key string `json:"key" validate:"required"`
}
