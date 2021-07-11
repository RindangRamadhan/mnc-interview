package helper

import (
	"net/http"

	"github.com/RindangRamadhan/mnc-interview/internal/tracing"

	"gitlab.com/pt-mai/maihelper"
	"google.golang.org/grpc"

	grpc_middlewares "gitlab.com/pt-mai/maihelper/maitracing/grpc-middlewares"
)

// DialService Dial service for HTTP request
func DialServiceWithRequest(service string, r *http.Request) (conn *grpc.ClientConn, err error) {
	serviceAddr := maihelper.Env.GrpcAddress(service)
	conn, err = grpc.Dial(
		serviceAddr,
		grpc.WithInsecure(),
		grpc_middlewares.ClientUnaryWithRequest(r, tracing.Tracer),
	)
	return
}

func DialService(service string) (conn *grpc.ClientConn, err error) {
	serviceAddr := maihelper.Env.GrpcAddress(service)
	conn, err = grpc.Dial(
		serviceAddr,
		grpc.WithInsecure(),
		grpc_middlewares.ClientUnary(tracing.Tracer),
	)
	return
}
