package service

// 数据处理可以在这哦
import (
	"HelpStudent/core/errorx"
	pingV1 "HelpStudent/gen/proto/ping/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type S struct {
	pingV1.UnimplementedPingServiceServer
}

func (s S) Ping(ctx context.Context, in *pingV1.PingRequest) (*pingV1.PingResponse, error) {
	return &pingV1.PingResponse{Value: in.Value}, nil
}

func (s S) PingErr(ctx context.Context, _ *pingV1.ABitOfEverything) (*emptypb.Empty, error) {
	//_, span := tracex.TracerFromContext(ctx).Start(ctx, "permissionErr")
	//defer span.End()
	//span.RecordError(errors.New("permissionErr not implemented"))
	//return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "permissionErr not implemented")
	return nil, errorx.NewErrCodeMsg(401000, "这是不可以的操作")
}

func (s S) PingAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}
