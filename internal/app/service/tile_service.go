package service

import (
	"context"
	"errors"
	"mytro-backend-tile/internal/app/repository"
	"mytro-backend-tile/internal/domain/models"
	"mytro-backend-tile/internal/infrastructure/config"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnsupportedContentType = errors.New("unsupported content type")

	AvailableContentTypes = map[string]string{
		"image/gif":       ".gif",
		"image/png":       ".png",
		"image/jpeg":      ".jpg",
		"image/webp":      ".webp",
		"video/mp4":       ".mp4",
		"video/quicktime": ".mov",
		"video/webm":      ".webm",
	}
)

var (
	videoContentTypes = typesWithPrefix("video/")
	imageContentTypes = typesWithPrefix("image/")
)

type TileService struct {
	tileRepository repository.TileRepository
	fileRepository repository.FileRepository
	txManager      repository.TxManager
	upload         config.UploadConfig
}

func NewTileService(tileRepository repository.TileRepository, fileRepository repository.FileRepository, txManager repository.TxManager, upload config.UploadConfig) TileService {
	return TileService{
		tileRepository: tileRepository,
		fileRepository: fileRepository,
		txManager:      txManager,
		upload:         upload,
	}
}

func (s TileService) CreateUpload(ctx context.Context, upload models.CreateUpload) (*models.CreatedUpload, error) {
	id := uuid.New()

	_, exists := AvailableContentTypes[upload.FileType]
	if !exists {
		return nil, ErrUnsupportedContentType
	}

	isVideo := slices.Contains(videoContentTypes, upload.FileType)

	err := s.tileRepository.Create(ctx, &models.Tile{
		ID:         id,
		OwnerID:    upload.OwnerID,
		Type:       getFileType(isVideo),
		Status:     models.TileStatusDraft,
		Visibility: models.TileVisibilityPrivate,
	})
	if err != nil {
		return nil, err
	}

	maxSize, acceptedTypes := s.upload.MaxImageSize, imageContentTypes
	var maxDuration *uint32
	if isVideo {
		maxSize, acceptedTypes = s.upload.MaxVideoSize, videoContentTypes
		maxDuration = &s.upload.MaxVideoDuration
	}

	url, err := s.fileRepository.CreateUploadURL(ctx, models.UploadURLData{
		TileID:               id,
		MaxSize:              maxSize,
		MaxDuration:          maxDuration,
		ContentType:          upload.FileType,
		AcceptedContentTypes: acceptedTypes,
		Tasks:                getTasks(isVideo),
		ExternalTasks:        []string{"mytro.nsfw", "mytro.embedding"},
		TTL:                  &s.upload.URLTTL,
	})
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		delErr := s.tileRepository.DeleteDraft(cleanupCtx, id)
		if delErr != nil {
			return nil, errors.Join(err, delErr)
		}

		return nil, err
	}

	createdURL := models.CreatedUpload{
		TileID:     id,
		FileID:     url.FileID,
		UploadURL:  url.UploadURL,
		FormFields: url.FormFields,
		ExpiresAt:  url.ExpiresAt,
	}

	return &createdURL, err
}

func (s TileService) CheckTileOwner(ctx context.Context, tileID uuid.UUID, userID uuid.UUID) (bool, error) {
	ownerID, err := s.tileRepository.GetOwnerByID(ctx, tileID)

	return ownerID == userID, err
}

func (s TileService) GetNewTiles(ctx context.Context, userID uuid.UUID) ([]models.Tile, error) {
	tileList, err := s.tileRepository.GetAllNewByOwner(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tileList, nil
}

func getFileType(isVideo bool) models.TileType {
	if isVideo {
		return models.TileTypeVideo
	} else {
		return models.TileTypeImage
	}
}

func getTasks(isVideo bool) []models.ProcessingTask {
	if isVideo {
		return []models.ProcessingTask{
			models.ProcessingTaskThumbnails,
			models.ProcessingTaskTranscode,
			models.ProcessingTaskFrames,
		}
	} else {
		return []models.ProcessingTask{
			models.ProcessingTaskThumbnails,
		}
	}
}

func typesWithPrefix(prefix string) []string {
	types := []string{}

	for t := range AvailableContentTypes {
		if strings.HasPrefix(t, prefix) {
			types = append(types, t)
		}
	}

	return types
}
