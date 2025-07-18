package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cliffordotieno/ai-context-gap-tracker/internal/config"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/contexttracker"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/database"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/logicengine"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/promptrewriter"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/responseauditor"
	"github.com/cliffordotieno/ai-context-gap-tracker/internal/server"
	"github.com/cliffordotieno/ai-context-gap-tracker/pkg/redis"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize modules
	contextTracker := contexttracker.New(db, redisClient)
	logicEngine := logicengine.New(db)
	responseAuditor := responseauditor.New(db)
	promptRewriter := promptrewriter.New(contextTracker, logicEngine)

	// Initialize HTTP server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	httpServer := server.NewHTTPServer(router, contextTracker, logicEngine, responseAuditor, promptRewriter)
	httpServer.SetupRoutes()

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	server.RegisterGRPCServices(grpcServer, contextTracker, logicEngine, responseAuditor, promptRewriter)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
		if err != nil {
			log.Fatal("Failed to listen on gRPC port:", err)
		}
		log.Printf("gRPC server listening on :%d", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: router,
	}

	go func() {
		log.Printf("HTTP server listening on :%d", cfg.Server.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	grpcServer.GracefulStop()
	log.Println("Server exited")
}