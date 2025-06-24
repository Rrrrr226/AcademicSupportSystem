package rpc

import (
	"HelpStudent/core/auth"
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

func AuthInterceptor(ctx context.Context) (context.Context, error) {

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, nil
	}
	entry, err := auth.ParseToken(token)
	if err == nil {
		return context.WithValue(ctx, "uid", entry.Info.Uid), nil
	}
	return ctx, nil
}

func GetUid(ctx context.Context) string {
	uid := ctx.Value("uid")
	if id, ok := uid.(string); ok {
		return id
	}
	return ""
}
