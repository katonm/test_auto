package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"my/sberAuto/test_auto/pb"
	"my/sberAuto/test_auto/pkg/endpoint"
	"my/sberAuto/test_auto/pkg/instrumenting"
	"my/sberAuto/test_auto/pkg/metric"
	"my/sberAuto/test_auto/pkg/service"
	"my/sberAuto/test_auto/pkg/transport"
)

type configuration struct {
	Port              string        `envconfig:"port" default:":50051"`
	MaxConnectionIdle time.Duration `envconfig:"max_idle" default:"5m"`
	Timeout           time.Duration `envconfig:"timeout" default:"15s"`
	MaxConnectionAge  time.Duration `envconfig:"max_age" default:"5m"`
	Time              time.Duration `envconfig:"time" default:"120m"`

	URL         string `envconfig:"url" default:"0.0.0.0:7071"`
	ServiceName string `envconfig:"service_name" default:"parenthesis_svc"`
}

func main() {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		_ = level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)
	}

	svc := service.NewService(logger)
	ep := endpoint.MakeEndpoints(svc)
	grpcServer := transport.NewGRPCServer(ep, logger)

	grpcListener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		_ = logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	metrics, err := metric.CreateMetrics(cfg.URL, cfg.ServiceName)
	if err != nil {
		level.Error(logger).Log("CreateMetrics Error: %s", err)
	}

	im := instrumenting.NewInstrumentingManager(metrics, logger)

	baseServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			Timeout:           cfg.Timeout,
			MaxConnectionAge:  cfg.MaxConnectionAge,
			Time:              cfg.Time,
		}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpcTags.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			grpcRecovery.UnaryServerInterceptor(),
		),
	)

	pb.RegisterParenthesesServiceServer(baseServer, grpcServer)
	grpcPrometheus.Register(baseServer)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		level.Info(logger).Log("msg", "Server started successfully")
		baseServer.Serve(grpcListener)
	}()

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Error(logger).Log("exit", <-errs)
}
