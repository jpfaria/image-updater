package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jpfaria/image-updater/internal/model"
	"github.com/xgodev/boost/wrapper/log"
)

// Client handles Docker registry operations
type Client struct {
	baseURL     string
	credentials *Credentials
	httpClient  *http.Client
}

// Credentials represents Docker registry credentials
type Credentials struct {
	Username string
	Password string
	Token    string
}

// NewClient creates a new Docker registry client
func NewClient(baseURL string, credentials *Credentials) *Client {
	return &Client{
		baseURL:     baseURL,
		credentials: credentials,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListTags lists all tags for a Docker image
func (c *Client) ListTags(ctx context.Context, namespace, repository string) ([]model.Tag, error) {
	log.Infof("Listing tags for %s/%s", namespace, repository)

	// Construct URL for Docker Registry API v2
	url := fmt.Sprintf("%s/v2/%s/%s/tags/list", c.baseURL, namespace, repository)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if provided
	if c.credentials != nil {
		if c.credentials.Token != "" {
			req.Header.Set("Authorization", "Bearer "+c.credentials.Token)
		} else if c.credentials.Username != "" && c.credentials.Password != "" {
			req.SetBasicAuth(c.credentials.Username, c.credentials.Password)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list tags: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var result struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to model.Tag
	tags := make([]model.Tag, 0, len(result.Tags))
	for _, tag := range result.Tags {
		// For each tag, we would normally get the digest and creation time
		// This would require additional API calls to get manifest for each tag
		// For simplicity, we're just using the tag name here
		tags = append(tags, model.Tag{
			Name:      tag,
			Digest:    "", // Would require additional API call
			CreatedAt: time.Now().Format(time.RFC3339), // Would require additional API call
		})
	}

	return tags, nil
}

// GetImageDigest gets the digest for a specific image tag
func (c *Client) GetImageDigest(ctx context.Context, namespace, repository, tag string) (string, error) {
	log.Infof("Getting digest for %s/%s:%s", namespace, repository, tag)

	// Construct URL for Docker Registry API v2
	url := fmt.Sprintf("%s/v2/%s/%s/manifests/%s", c.baseURL, namespace, repository, tag)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add accept header for manifest v2
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// Add authentication if provided
	if c.credentials != nil {
		if c.credentials.Token != "" {
			req.Header.Set("Authorization", "Bearer "+c.credentials.Token)
		} else if c.credentials.Username != "" && c.credentials.Password != "" {
			req.SetBasicAuth(c.credentials.Username, c.credentials.Password)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get digest: %s - %s", resp.Status, string(body))
	}

	// Get digest from header
	digest := resp.Header.Get("Docker-Content-Digest")
	if digest == "" {
		return "", fmt.Errorf("digest not found in response headers")
	}

	return digest, nil
}
