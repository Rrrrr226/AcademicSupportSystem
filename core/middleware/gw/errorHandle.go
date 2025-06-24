package gw

import (
	"HelpStudent/core/jsonx"
	"HelpStudent/core/middleware/response"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func GrpcGatewayError(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Internal, err.Error())
	}

	httpError := response.JsonResponse{Code: int32(s.Code()), Message: s.Message(), Error: s.Details()}

	resp, _ := jsonx.Marshal(httpError)
	w.Header().Set("Content-type", "application/json")

	code := int(s.Code()) / 1000
	if isGrpcErr(s.Code()) {
		code = runtime.HTTPStatusFromCode(s.Code())
	} else {
		if code == 0 {
			code = 500
		}
	}

	w.WriteHeader(code)
	_, _ = w.Write(resp)
}

func isGrpcErr(code codes.Code) bool {
	return code < 16
}
