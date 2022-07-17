package middleware

import (
	"context"
	"log"
	"strings"

	helper "git.biggorilla.tech/gateway/payment-gateway/helpers"
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
	values := md["authorization"]
	accessToken := values[0]
	token := strings.Split(accessToken, " ")
	claims, err := r.jwtManager.Verify(token[1])
	if err == nil {
		return handler(ctx, req)
	}
	log.Println("love --> unary interceptor: ", info.FullMethod, req, claims)
	return nil, err
}
