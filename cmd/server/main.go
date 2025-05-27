package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/iktakahiro/oniongo/internal/infrastructure/di"
	v1connect "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1/oniongov1connect"
	"github.com/iktakahiro/oniongo/internal/infrastructure/grpc/interceptor"
	"github.com/rs/cors"
	"github.com/samber/do"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	injector := di.DependencyInjection()

	// Get port from environment variable, default to 8080
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Set service names for reflection
	reflector := grpcreflect.NewStaticReflector(
		v1connect.TodoServiceName,
	)

	todoServiceHandler, err := do.Invoke[v1connect.TodoServiceHandler](injector)
	if err != nil {
		log.Fatalf("failed to invoke todo service handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(v1connect.NewTodoServiceHandler(
		todoServiceHandler,
		[]connect.HandlerOption{
			connect.WithCompressMinBytes(2048),
			connect.WithSendMaxBytes(4 * 1024 * 1024),
			connect.WithReadMaxBytes(4 * 1024 * 1024),
			connect.WithInterceptors(
				interceptor.NewLoggingInterceptor(),
			),
		}...,
	))

	corsOption := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		MaxAge:         int(2 * time.Hour / time.Second),
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: h2c.NewHandler(
			corsOption.Handler(mux),
			&http2.Server{},
		),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 8 * 1024,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Printf("start server, port: %v\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err)
	}
}
