package grpc

import (
	"context"
	"fmt"
	authv1 "github.com/TauAdam/sso/contracts/gen/go/sso"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api authv1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	address string,
	retriesNumber int,
	timeout time.Duration,
) (*Client, error) {
	const op = "grpc.New"

	retryOptions := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		retry.WithMax(uint(retriesNumber)),
		retry.WithPerRetryTimeout(timeout),
	}

	logOptions := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(
				InterceptorsLogger(log),
				logOptions...,
			),
			retry.UnaryClientInterceptor(retryOptions...)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s failed to dial: %w", op, err)
	}

	return &Client{
		api: authv1.NewAuthClient(conn),
	}, nil
}

// InterceptorsLogger wraps slog.Logger to grpc logging.Logger
func InterceptorsLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "grpc.IsAdmin"

	response, err := c.api.IsAdmin(ctx, &authv1.IsAdminRequest{UserId: userID})
	if err != nil {
		return false, fmt.Errorf("%s failed to call: %w", op, err)
	}

	return response.IsAdmin, nil
}
