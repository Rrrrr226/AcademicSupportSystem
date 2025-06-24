package gw

import (
	"HelpStudent/core/errorx"
	"HelpStudent/core/tracex"
	pingV1 "HelpStudent/gen/proto/ping/v1"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type S struct {
	pingV1.UnimplementedPingServiceServer
}

func (s S) Ping(ctx context.Context, in *pingV1.PingRequest) (*pingV1.PingResponse, error) {
	return &pingV1.PingResponse{Value: in.Value}, nil
}

func (s S) PingErr(ctx context.Context, _ *pingV1.ABitOfEverything) (*emptypb.Empty, error) {
	_, span := tracex.TracerFromContext(ctx).Start(ctx, "PingErr")
	defer span.End()
	span.RecordError(errors.New("PingErr not implemented"))
	//return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "PingErr not implemented")
	return nil, errorx.NewErrCodeMsg(401000, "这是不可以的操作")
}

func (s S) PingAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}
