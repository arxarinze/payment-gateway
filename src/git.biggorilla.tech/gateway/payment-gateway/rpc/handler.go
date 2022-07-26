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

func (s *server) GetPublicMerchantInfo(ctx context.Context, in *pb.MerchantPublicRequest) (*pb.MerchantPublicResponse, error) {
	data, err := s.repo.GetPublicMerchantInfo(ctx, in.GetPluginId())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *server) GetPluginLink(ctx context.Context, in *pb.PluginLinkRequest) (*pb.LinkResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	link := s.repo.GetPluginLink(ctx, id, in.GetMerchantId(), in.GetType())
	return &pb.LinkResponse{
		Link: link,
	}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.LinkResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	data, err1 := s.repo.GenerateLink(ctx, in.GetMerchantId(), id)
	if err1 != nil {
		return &pb.LinkResponse{
			Link: "",
		}, err1
	}
	a := data.(*model.Link)
	return &pb.LinkResponse{
		Link: a.PluginID,
	}, err1
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
	switch strings.ToLower(in.Network) {
	case "ethereum":
		publicKey, err := s.repo.GenerateDepositAddress(ctx, s.ethereumClient, in.Network, in.Cryptosymbol, in.PluginId)
		if err != nil {
			return &pb.DepositAddressResponse{}, err
		}
		return &pb.DepositAddressResponse{
			Address: publicKey,
		}, nil
	case "bitcoin":
		return nil, nil
	}
	return nil, nil
}
