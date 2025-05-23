package server

import (
	"context"

	"github.com/jpfaria/image-updater/internal/config"
	"github.com/jpfaria/image-updater/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
}

// New creates a new server instance
func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	// Create plugins
	corsPlugin := cors.New(ctx)
	recoverPlugin := recover.New(ctx)
	logPlugin := log.New(ctx)

	// Create Echo server with Boost integration
	echoServer, err := echoserver.NewServer(ctx, corsPlugin, recoverPlugin, logPlugin)
	if err != nil {
		return nil, err
	}

	// Create server instance
	server := &Server{
		echo:   echoServer,
		config: cfg,
	}

	// Register routes
	server.registerRoutes()

	return server, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	// Configure server options
	options := &echoserver.Options{
		Port:       s.config.Server.Port,
		Host:       s.config.Server.Host,
		Protocol:   "HTTP",
		HideBanner: false,
	}

	// Start the server with options
	return s.echo.StartWithOptions(ctx, options)
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes() {
	// API group
	api := s.echo.Group("/api")

	// Docker image routes
	dockerHandler := handler.NewDockerHandler()
	api.GET("/images", dockerHandler.ListImages)
	api.GET("/images/:id", dockerHandler.GetImage)
	api.GET("/images/:id/tags", dockerHandler.ListTags)
	api.POST("/images/:id/refresh", dockerHandler.RefreshTags)

	// Environment routes
	envHandler := handler.NewEnvironmentHandler()
	api.GET("/environments", envHandler.ListEnvironments)
	api.GET("/environments/:id", envHandler.GetEnvironment)
	api.POST("/environments/:id/deploy", envHandler.DeployToEnvironment)

	// Git routes
	gitHandler := handler.NewGitHandler()
	api.GET("/repositories", gitHandler.ListRepositories)
	api.GET("/repositories/:id/files", gitHandler.ListFiles)
	api.GET("/repositories/:id/files/:path", gitHandler.GetFile)

	// Webhook routes
	webhookHandler := handler.NewWebhookHandler()
	api.POST("/webhooks/docker", webhookHandler.DockerWebhook)

	// Health check
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}
