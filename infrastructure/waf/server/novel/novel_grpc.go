package novel

import (
	"context"

	"narou/adapter/logger"
	metadataRp "narou/adapter/repository/metadata"
	"narou/infrastructure/database"
	metadataUc "narou/usecase/interactor/metadata"
	pb "narou/usecase/port/boudary/proto/novel"
)

func (s *Service) Get(ctx context.Context, _ *pb.Req) (*pb.Novels, error) {
	lg := logger.NewLogger(ctx)
	lg.Info("grpc server get start")
	defer lg.Info("grpc server get end")

	con, _ := database.GetConn()
	lst, _ := metadataUc.NewMetaDataListInteractor(ctx,
		lg, metadataRp.NewRepository(con),
		nil).Execute()

	return pb.Convert2ProtoBuf(lst), nil
}

type Service struct {
	pb.UnimplementedNovelListServer
}

func NewGrpcService() *Service {
	return &Service{}
}
