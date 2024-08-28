package novel

import (
	"context"

	"narou/domain/metadata"
	"narou/infrastructure/database"
	"narou/sdk/logger"
	metadataUc "narou/usecase/metadata"
	pb "narou/usecase/port/boudary/proto/novel"
)

func (s *Service) Get(ctx context.Context, _ *pb.Req) (*pb.Novels, error) {
	lg, err := logger.NewServerLogger(ctx)
	if err != nil {
		return nil, err
	}
	lg.Info("grpc server get start")
	defer lg.Info("grpc server get end")

	con, err := database.GetConn()
	if err != nil {
		return nil, err
	}
	lst, _ := metadataUc.NewMetaDataListUseCase(lg, metadata.NewRepository(con), nil).Execute(ctx)

	return pb.Convert2ProtoBuf(lst), nil
}

type Service struct {
	pb.UnimplementedNovelListServer
}

func NewGrpcService() *Service {
	return &Service{}
}
