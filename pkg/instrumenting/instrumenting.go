package instrumenting

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"my/sberAuto/test_auto/pkg/grpcerrors"
	"my/sberAuto/test_auto/pkg/metric"
)

// Manager...
type Manager struct {
	logger  log.Logger
	metrics metric.Metrics
}

// Manager constructor...
func NewInstrumentingManager(m metric.Metrics, l log.Logger) *Manager {
	return &Manager{
		metrics: m,
		logger:  l,
	}
}

// Logger Instrumenting...
func (im *Manager) Logger(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{},
	err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Log("Method: %s, Time: %v, Metadata: %v, Err: %v\n",
		info.FullMethod,
		time.Since(start),
		md,
		err)

	return reply, err
}

func (im *Manager) Metrics(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	var status = http.StatusOK
	if err != nil {
		status = grpcerrors.ErrCodeToHTTPStatus(grpcerrors.ParseErrStatusCode(err))
	}
	im.metrics.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.metrics.IncHits(status, info.FullMethod, info.FullMethod)

	return resp, err
}
