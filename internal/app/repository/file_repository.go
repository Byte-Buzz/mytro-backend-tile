package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mytro-backend-tile/internal/domain/models"
	storagev1 "mytro-backend-tile/internal/infrastructure/grpc/pb/content_storage/v1"
)

type FileRepository interface {
	CreateUploadURL(ctx context.Context, data models.UploadURLData) (*models.UploadURL, error)
}

type fileRepository struct {
	grpcClient storagev1.ContentStorageServiceClient
}

func NewFileRepository(grpcClient storagev1.ContentStorageServiceClient) FileRepository {
	return &fileRepository{
		grpcClient: grpcClient,
	}
}

func (r *fileRepository) CreateUploadURL(ctx context.Context, data models.UploadURLData) (*models.UploadURL, error) {
	info := map[string]any{
		"tile_id": data.TileID.String(),
	}

	infoJson, err := json.Marshal(info)

	if err != nil {
		return nil, err
	}

	tasks := make([]storagev1.ProcessingTask, len(data.Tasks))
	for i, task := range data.Tasks {
		t, err := toProtoTask(task)
		if err != nil {
			return nil, err
		}

		tasks[i] = t
	}

	response, err := r.grpcClient.CreateUpload(ctx, &storagev1.CreateUploadRequest{
		AppName:              "mytro",
		ContentType:          data.ContentType,
		AcceptedContentTypes: data.AcceptedContentTypes,
		MaxSize:              data.MaxSize,
		MaxDurationSeconds:   data.MaxDuration,
		Tasks:                tasks,
		ExternalTasks:        data.ExternalTasks,
		Info:                 infoJson,
		UploadUrlTtl:         data.TTL,
	})

	if err != nil {
		return nil, err
	}

	return &models.UploadURL{
		FileID:     response.FileId,
		UploadURL:  response.UploadUrl,
		FormFields: response.FormFields,
		ExpiresAt:  response.ExpiresAt.AsTime(),
	}, nil
}

func toProtoTask(t models.ProcessingTask) (storagev1.ProcessingTask, error) {
	switch t {
	case models.ProcessingTaskThumbnails:
		return storagev1.ProcessingTask_PROCESSING_TASK_THUMBNAILS, nil
	case models.ProcessingTaskFrames:
		return storagev1.ProcessingTask_PROCESSING_TASK_FRAMES, nil
	case models.ProcessingTaskTranscode:
		return storagev1.ProcessingTask_PROCESSING_TASK_TRANSCODE, nil
	default:
		return 0, fmt.Errorf("unknown processing task: %v", t)
	}
}
