package interceptor

import (
	"context"
	"log"
	"time"

	"connectrpc.com/connect"
)

func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			log.Printf("RPC Start: %s", req.Spec().Procedure)

			res, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				log.Printf("RPC Error: %s, Duration: %v, Error: %v",
					req.Spec().Procedure, duration, err)
			} else {
				log.Printf("RPC Success: %s, Duration: %v",
					req.Spec().Procedure, duration)
			}

			return res, err
		}
	}
}
