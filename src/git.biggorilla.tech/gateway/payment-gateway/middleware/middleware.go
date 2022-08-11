package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	helper "git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	_ "git.biggorilla.tech/gateway/payment-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Middleware interface {
	UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type middleware struct {
	jwtManager helper.JWTManager
}

func NewMiddleware(ctx context.Context, jwtManager helper.JWTManager) Middleware {
	return &middleware{
		jwtManager: jwtManager,
	}
}

func (r *middleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		panic(ok)
	}
	fmt.Println(info.FullMethod)
	if info.FullMethod == "/payment_gateway.v1.PaymentGatewayService/GetPublicMerchantInfo" || info.FullMethod == "/payment_gateway.v1.PaymentGatewayService/GenerateDepositAddress" {
		return handler(ctx, req)
	}
	values := md["authorization"]
	accessToken := values[0]
	client := http.Client{}
	requ, err := http.NewRequest("GET", "http://ec2-52-72-83-242.compute-1.amazonaws.com/global/verify-sso", nil)
	if err != nil {
		return &pb.GenericResponse{
			Code:    403,
			Message: "Forbidden",
		}, nil
	}

	requ.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {accessToken},
	}

	res, err := client.Do(requ)
	if err != nil {
		return &pb.GenericResponse{
			Code:    403,
			Message: "Forbidden",
		}, nil
	}
	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	if data["status"] == false {
		return &pb.GenericResponse{
			Code:    403,
			Message: "Forbidden",
		}, nil
	}
	//_, err := r.jwtManager.Verify(token[1])

	// if err != nil {
	// 	return &pb.GenericResponse{
	// 		Code:    403,
	// 		Message: "Forbidden",
	// 	}, nil
	// }
	//log.Println("love --> unary interceptor: ", claims)
	return handler(context.WithValue(ctx, "user", data), req)
}
