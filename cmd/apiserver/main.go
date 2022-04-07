package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	linuxcncv1 "github.com/thetechnick/linuxcnc-ui/api/v1"
	"github.com/thetechnick/linuxcnc-ui/internal/apiserver"
	apiserverv1 "github.com/thetechnick/linuxcnc-ui/internal/apiserver/v1"
	"github.com/thetechnick/linuxcnc-ui/internal/rs274ngc/rs274ngcinterop"
)

func main() {
	var c apiserver.APIServerConfig

	flag.StringVar(&c.Address, "address", ":8080", "server secure address")
	flag.StringVar(&c.TLSCertFile, "tls-cert-file", "cert.pem", "server certificate")
	flag.StringVar(&c.TLSKeyFile, "tls-key-file", "key.pem", "server key")

	flag.Parse()

	if err := run(c); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(c apiserver.APIServerConfig) error {
	register := []apiserver.RegisterFn{
		func(ctx context.Context, grpcServer *grpc.Server,
			grpcGatewayMux *gwruntime.ServeMux, grpcClient *grpc.ClientConn) error {

			healthService := &apiserverv1.HealthServiceServer{}
			linuxcncv1.RegisterHealthServiceServer(grpcServer, healthService)
			if err := linuxcncv1.RegisterHealthServiceHandler(ctx, grpcGatewayMux, grpcClient); err != nil {
				return fmt.Errorf("register health service handler: %w", err)
			}

			fileRoot := "/tmp/test"

			fileService := &apiserverv1.FilesServiceServer{
				Root: fileRoot,
			}
			linuxcncv1.RegisterFilesServiceServer(grpcServer, fileService)
			if err := linuxcncv1.RegisterFilesServiceHandler(ctx, grpcGatewayMux, grpcClient); err != nil {
				return fmt.Errorf("register files service handler: %w", err)
			}

			interpreterService := &apiserverv1.InterpreterServiceServer{
				Interpreter: rs274ngcinterop.Interpreter,
				Root:        fileRoot,
			}
			linuxcncv1.RegisterInterpreterServiceServer(grpcServer, interpreterService)
			if err := linuxcncv1.RegisterInterpreterServiceHandler(ctx, grpcGatewayMux, grpcClient); err != nil {
				return fmt.Errorf("register interpreter service handler: %w", err)
			}
			return nil
		},
	}

	s, err := apiserver.NewServer(c, register...)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	// Signal Handler
	serverStopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		close(serverStopCh)
	}()

	if err := s.Run(serverStopCh); err != nil {
		return fmt.Errorf("running server: %w", err)
	}
	return nil
}
