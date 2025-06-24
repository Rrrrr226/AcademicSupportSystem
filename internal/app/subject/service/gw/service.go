package service

// 数据处理可以在这哦
import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"HelpStudent/core/errorx"
	"HelpStudent/core/tracex"
	subjectV1 "HelpStudent/gen/proto/subject/v1"
)

type S struct {
	subjectV1.UnimplementedPingServiceServer
}

func (s S) Subject(ctx context.Context, in *subjectV1.PingRequest) (*subjectV1.PingResponse, error) {
	return &subjectV1.PingResponse{Value: in.Value}, nil
}

func (s S) SubjectErr(ctx context.Context, _ *subjectV1.ABitOfEverything) (*emptypb.Empty, error) {
	_, span := tracex.TracerFromContext(ctx).Start(ctx, "subjectErr")
	defer span.End()
	span.RecordError(errors.New("subjectErr not implemented"))
	//return nil, errors.Wrapf(errorx.NewErrCodeMsg(502000, "这是不可以的操作"), "subjectErr not implemented")
	return nil, errorx.NewErrCodeMsg(401000, "这是不可以的操作")
}

func (s S) SubjectAuth(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, errorx.UnAuthorizedError
}