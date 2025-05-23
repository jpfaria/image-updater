package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/log"
)

// WebhookHandler handles webhook related requests
type WebhookHandler struct {
	// Add service dependencies here
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// DockerWebhook handles Docker registry webhooks
func (h *WebhookHandler) DockerWebhook(c echo.Context) error {
	log.Info("Received Docker webhook")
	
	// Parse request body
	var req struct {
		Repository string `json:"repository"`
		Tag        string `json:"tag"`
		Digest     string `json:"digest"`
		Namespace  string `json:"namespace"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	
	log.Infof("Docker webhook for %s/%s:%s", req.Namespace, req.Repository, req.Tag)
	
	// Mock response for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Webhook processed successfully",
	})
}
