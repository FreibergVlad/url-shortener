package integration

import (
	"context"
	"net"
	"testing"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func MaybeSkipIntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}

func BufnetGrpcClient(listener *bufconn.Listener) *grpc.ClientConn {
	return must.Return(
		grpc.NewClient(
			"passthrough://bufnet",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return listener.Dial()
			}),
		),
	)
}
