package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"mytro-backend-tile/internal/app/repository"
	"mytro-backend-tile/internal/app/service"
	"mytro-backend-tile/internal/infrastructure/config"
	storagev1 "mytro-backend-tile/internal/infrastructure/grpc/pb/content_storage/v1"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	GRPC   *storagev1.ContentStorageServiceClient
	DB     *gorm.DB

	TxManager    repository.TxManager
	Services     *AppServices
	repositories *AppRepositories
}

type AppServices struct {
	TileService service.TileService
}

type AppRepositories struct {
	TileRepository repository.TileRepository
	FileRepository repository.FileRepository
}

func NewApp(config *config.Config, db *gorm.DB, logger *zap.Logger, grpcClient *storagev1.ContentStorageServiceClient) *App {
	app := &App{
		Config: config,
		Logger: logger,
		GRPC:   grpcClient,
		DB:     db,

		TxManager:    repository.NewTxManager(db),
		Services:     &AppServices{},
		repositories: &AppRepositories{},
	}

	app.createRepositories()
	app.createServices()

	return app
}

func (app *App) createRepositories() {
	app.repositories.TileRepository = repository.NewTileRepository(app.Logger, app.DB)
	app.repositories.FileRepository = repository.NewFileRepository(*app.GRPC)
}

func (app *App) createServices() {
	app.Services.TileService = service.NewTileService(app.repositories.TileRepository, app.repositories.FileRepository, app.TxManager, app.Config.Upload)
}

func (app *App) GetPublicKey() (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode(([]byte)(app.Config.Keys.PublicKey()))

	if publicKeyBlock == nil {
		app.Logger.Fatal("failed to parse PEM block containing the key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		app.Logger.Fatal("failed to parse public key: " + err.Error())
	}

	return publicKey.(*rsa.PublicKey), nil
}
