package model

// Image represents a Docker image
type Image struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Registry  string `json:"registry"`
	Namespace string `json:"namespace"`
	LatestTag string `json:"latest_tag,omitempty"`
}

// Tag represents a Docker image tag
type Tag struct {
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	CreatedAt string `json:"created_at"`
}

// Environment represents a deployment environment
type Environment struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Application  string `json:"application"`
	ValuesPath   string `json:"values_path"`
	CurrentImage string `json:"current_image,omitempty"`
}

// Deployment represents a deployment record
type Deployment struct {
	ID        string `json:"id"`
	ImageTag  string `json:"image_tag"`
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
	Status    string `json:"status"`
}

// Repository represents a Git repository
type Repository struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Branch   string `json:"branch"`
	TeamName string `json:"team_name"`
}

// File represents a file in a Git repository
type File struct {
	Path       string `json:"path"`
	Type       string `json:"type"` // "file" or "directory"
	LastCommit string `json:"last_commit"`
	LastUpdate string `json:"last_update"`
	Content    string `json:"content,omitempty"`
}
