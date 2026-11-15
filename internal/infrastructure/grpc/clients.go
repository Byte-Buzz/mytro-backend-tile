package grpc

import (
	storagev1 "mytro-backend-tile/internal/infrastructure/grpc/pb/content_storage/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewContextClient(addr string) (storagev1.ContentStorageServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return storagev1.NewContentStorageServiceClient(conn), nil
}
