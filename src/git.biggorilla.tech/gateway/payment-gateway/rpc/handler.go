package rpc

import (
	"context"
	"strings"

	"git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/model"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
	"git.biggorilla.tech/gateway/payment-gateway/services/web3"
	"google.golang.org/grpc/metadata"
)

type server struct {
	identity       helpers.Identity
	repo           repo.MerchantRepo
	ethereumClient services.EthereumService
}

func NewRPCInterface(ctx context.Context, identity helpers.Identity, repo repo.MerchantRepo, ethereumClient services.EthereumService) pb.PaymentGatewayServiceServer {
	return &server{
		identity,
		repo,
		ethereumClient,
	}
}
func (s *server) GetPluginLink(ctx context.Context, in *pb.PluginLinkRequest) (*pb.GenericResponse, error) {

	return &pb.GenericResponse{}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.GenericResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	in.MerchantId = id
	data, err1 := s.repo.GenerateLink(ctx, in.GetMerchantId())
	if err1 != nil {
		return &pb.GenericResponse{
			Code:    500,
			Message: err1.Error(),
		}, nil
	}
	a := data.(*model.GenericResponse)
	return &pb.GenericResponse{
		Code:    int32(a.Code),
		Message: a.Message,
	}, nil
}
func (s *server) CreateMerchant(ctx context.Context, in *pb.MerchantRequest) (*pb.GenericResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	data, err1 := s.repo.CreateMerchant(ctx, in.GetName(), in.GetEmail(), id)
	if err1 != nil {
		a := data.(*model.GenericResponse)
		return &pb.GenericResponse{
			Code:    int32(a.Code),
			Message: a.Message,
		}, nil
	}
	return &pb.GenericResponse{
		Code:    200,
		Message: "Successfully Created Merchant",
	}, nil
}

func (s *server) GenerateDepositAddress(ctx context.Context, in *pb.DepositAddressRequest) (*pb.DepositAddressResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	switch strings.ToLower(in.Network) {
	case "ethereum":
		publicKey := s.repo.GenerateDepositAddress(ctx, s.ethereumClient, in.Network, in.Cryptosymbol, id)
		return &pb.DepositAddressResponse{
			Address: publicKey,
		}, nil
	case "bitcoin":
		return nil, nil
	}
	return nil, nil
}
