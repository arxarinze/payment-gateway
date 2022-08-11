package rpc

import (
	"context"
	"strconv"
	"strings"

	"git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/model"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/pb/transform"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
	services "git.biggorilla.tech/gateway/payment-gateway/services/web3"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	identity       helpers.Identity
	repo           repo.MerchantRepo
	ethereumClient services.EthereumService
}

// GetTransactions implements pb.PaymentGatewayServiceServer
func (s *server) GetTransactions(ctx context.Context, in *pb.TransactionRequest) (*pb.Transactions, error) {
	id := s.identity.GetIdentity(ctx)
	data, err := s.repo.GetTransactions(ctx, in.GetMerchantId(), id)
	if err != nil {
		return nil, err
	}
	return &pb.Transactions{
		Data: transform.TransactionToPbList(*data),
	}, nil
}

// UpdateMerchant implements pb.PaymentGatewayServiceServer
func (s *server) UpdateMerchant(ctx context.Context, in *pb.MerchantUpdateRequest) (*pb.GenericResponse, error) {
	id := s.identity.GetIdentity(ctx)
	data, err := s.repo.UpdateMerchant(ctx, in.GetName(), in.GetEmail(), in.GetAddress(), in.GetAvatar(), id, in.GetMerchantId())
	if err != nil {
		return nil, err
	}
	count := data.(int64)
	stringCount := strconv.Itoa(int(count))
	return &pb.GenericResponse{
		Code:    200,
		Message: "Updated " + stringCount + " Records",
	}, nil
}

// GetMerchant implements pb.PaymentGatewayServiceServer
func (s *server) GetMerchants(ctx context.Context, in *emptypb.Empty) (*pb.MerchantResponse, error) {
	id := s.identity.GetIdentity(ctx)
	data, err := s.repo.GetMerchants(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.MerchantResponse{
		Data: transform.MerchantToPbList(*data),
	}, nil
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
	//auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(ctx)
	link, err := s.repo.GetPluginLink(ctx, id, in.GetMerchantId(), in.GetType())
	if err != nil {
		return nil, err
	}
	return &pb.LinkResponse{
		Link: link,
	}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.LinkResponse, error) {

	id := s.identity.GetIdentity(ctx)
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
	// auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(ctx)
	data, err1 := s.repo.CreateMerchant(ctx, in.GetName(), in.GetEmail(), in.GetAddress(), in.GetAvatar(), id)
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
