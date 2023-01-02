package server

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
)

func RunServers(lc fx.Lifecycle, log logger.ILogger, echoServer *server.EchoServer, grpcServer *grpc.GrpcServer, ctx context.Context) error {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoServer.RunHttpServer(ctx); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			go func() {
				if err := grpcServer.RunGrpcServer(ctx); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running grpc server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
