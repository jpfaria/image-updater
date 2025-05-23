package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/log"
)

// EnvironmentHandler handles environment related requests
type EnvironmentHandler struct {
	// Add service dependencies here
}

// NewEnvironmentHandler creates a new environment handler
func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{}
}

// ListEnvironments lists all environments
func (h *EnvironmentHandler) ListEnvironments(c echo.Context) error {
	log.Info("Listing environments")
	
	// Mock response for now
	environments := []map[string]interface{}{
		{
			"id":            "1",
			"name":          "production",
			"application":   "my-app",
			"values_path":   "my-app/production/values.yaml",
			"current_image": "nginx:1.24.0",
		},
		{
			"id":            "2",
			"name":          "staging",
			"application":   "my-app",
			"values_path":   "my-app/staging/values.yaml",
			"current_image": "nginx:1.25.0",
		},
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   environments,
	})
}

// GetEnvironment gets an environment by ID
func (h *EnvironmentHandler) GetEnvironment(c echo.Context) error {
	id := c.Param("id")
	log.Infof("Getting environment with ID: %s", id)
	
	// Mock response for now
	environment := map[string]interface{}{
		"id":            id,
		"name":          "production",
		"application":   "my-app",
		"values_path":   "my-app/production/values.yaml",
		"current_image": "nginx:1.24.0",
		"deployments": []map[string]interface{}{
			{
				"id":        "1",
				"image_tag": "1.24.0",
				"timestamp": "2025-04-15T10:00:00Z",
				"user":      "admin",
				"status":    "success",
			},
			{
				"id":        "2",
				"image_tag": "1.23.0",
				"timestamp": "2025-03-10T10:00:00Z",
				"user":      "admin",
				"status":    "success",
			},
		},
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   environment,
	})
}

// DeployToEnvironment deploys an image to an environment
func (h *EnvironmentHandler) DeployToEnvironment(c echo.Context) error {
	id := c.Param("id")
	log.Infof("Deploying to environment with ID: %s", id)
	
	// Parse request body
	var req struct {
		ImageTag string `json:"image_tag"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	
	// Mock response for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Deployment initiated",
		"data": map[string]interface{}{
			"deployment_id": "123",
			"environment":   "production",
			"image_tag":     req.ImageTag,
			"status":        "pending",
		},
	})
}
