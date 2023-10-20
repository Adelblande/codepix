package grpc

import (
	"context"

	"github.com/adelblande/codepix/application/grpc/pb"
	"github.com/adelblande/codepix/application/usecase"
)

type PixGrpcService struct {
	PixUseCase usecase.PixUseCase
	pb.UnimplementedPixServiceServer
}

func (p *PixGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	key, err := p.PixUseCase.RegisterKey(in.Key, in.Kind, in.AccountId)
	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error: err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Id: key.ID,
		Status: "created",
	}, nil
}

func (p *PixGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	info, err := p.PixUseCase.FindKey(in.Kind, in.Key)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id: info.ID,
		Kind: info.Kind,
		Key: info.Key,
		Account: &pb.Account{
			AccountId: info.AccountId,
			AccountNumber: info.Account.Number,
			BankId: info.Account.BankID,
			BankName: info.Account.BankID,
			OwnerName: info.Account.OwnerName,
			CreatedAt: info.CreatedAt.String(),
		},
		CreatedAt: info.CreatedAt.String(),
	}, nil
}

func NewPixGrpcService(usecase usecase.PixUseCase) *PixGrpcService {
	return &PixGrpcService{
		PixUseCase: usecase,
	}
}