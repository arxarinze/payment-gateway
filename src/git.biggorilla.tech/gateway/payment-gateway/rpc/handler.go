package rpc

import (
	"context"
	"encoding/json"

	"git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
	"google.golang.org/grpc/metadata"
	_ "google.golang.org/grpc/metadata"
)

type server struct {
	identity helpers.Identity
	repo     repo.MerchantRepo
}

func NewRPCInterface(ctx context.Context, identity helpers.Identity, repo repo.MerchantRepo) pb.PaymentGatewayServiceServer {
	return &server{
		identity,
		repo,
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
	out, err2 := json.Marshal(data)
	if err2 != nil || err1 != nil {
		return &pb.GenericResponse{
			Code:    500,
			Message: err2.Error() + err1.Error(),
		}, nil
	}
	return &pb.GenericResponse{
		Code:    200,
		Message: string(out),
	}, nil
}
func (s *server) CreateMerchant(ctx context.Context, in *pb.MerchantRequest) (*pb.GenericResponse, error) {
	auth, _ := metadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	data, err1 := s.repo.CreateMerchant(ctx, in.GetName(), in.GetEmail(), id)
	out, err2 := json.Marshal(data)
	if err2 != nil || err1 != nil {
		return &pb.GenericResponse{
			Code:    500,
			Message: err2.Error() + err1.Error(),
		}, nil
	}
	return &pb.GenericResponse{
		Code:    200,
		Message: string(out),
	}, nil
}

func (s *server) GenerateDepositAddress(ctx context.Context, in *pb.DepositAddressRequest) (*pb.DepositAddressResponse, error) {
	return nil, nil
}
