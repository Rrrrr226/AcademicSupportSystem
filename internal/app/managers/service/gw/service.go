package service

// 数据处理可以在这哦
import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"HelpStudent/core/errorx"
	"HelpStudent/core/tracex"
	managersV1 "HelpStudent/gen/proto/managers/v1"
)

type S struct {
	managersV1.UnimplementedPingServiceServer
}

func (s S) Managers(ctx context.Context, in *managersV1.PingRequest) (*managersV1.PingResponse, error) {
	return &managersV1.PingResponse{Value: in.Value}, nil
}

func (s S) ManagersErr(ctx context.Context, _ *managersV1.ABitOfEverything) (*emptypb.Empty, error) {
	_, span := tracex.TracerFromContext(ctx).Start(ctx, "managersErr")
	defer span.End()
	span.RecordError(errors.New("managersErr not implemented"))
	//return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "managersErr not implemented")
	return nil, errorx.NewErrCodeMsg(401000, "这是不可以的操作")
}

func (s S) ManagersAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}