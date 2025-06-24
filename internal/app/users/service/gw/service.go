package service

// 数据处理可以在这哦
import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"HelpStudent/core/errorx"
	"HelpStudent/core/tracex"
	usersV1 "HelpStudent/gen/proto/users/v1"
)

type S struct {
	usersV1.UnimplementedPingServiceServer
}

func (s S) Users(ctx context.Context, in *usersV1.PingRequest) (*usersV1.PingResponse, error) {
	return &usersV1.PingResponse{Value: in.Value}, nil
}

func (s S) UsersErr(ctx context.Context, _ *usersV1.ABitOfEverything) (*emptypb.Empty, error) {
	_, span := tracex.TracerFromContext(ctx).Start(ctx, "usersErr")
	defer span.End()
	span.RecordError(errors.New("usersErr not implemented"))
	//return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "usersErr not implemented")
	return nil, errorx.NewErrCodeMsg(401000, "这是不可以的操作")
}

func (s S) UsersAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}