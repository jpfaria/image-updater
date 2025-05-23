package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jpfaria/image-updater/internal/model"
	"github.com/xgodev/boost/wrapper/log"
)

// Client handles Git repository operations
type Client struct {
	auth        transport.AuthMethod
	cloneDir    string
	defaultUser string
	defaultEmail string
}

// NewClient creates a new Git client
func NewClient(authType, username, password, sshKeyPath, cloneDir string) (*Client, error) {
	var auth transport.AuthMethod
	var err error

	// Set up authentication
	if authType == "ssh" {
		auth, err = ssh.NewPublicKeysFromFile("git", sshKeyPath, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create SSH auth: %w", err)
		}
	} else {
		// Default to HTTPS
		auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	// Create clone directory if it doesn't exist
	if err := os.MkdirAll(cloneDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create clone directory: %w", err)
	}

	return &Client{
		auth:        auth,
		cloneDir:    cloneDir,
		defaultUser: "Image Updater",
		defaultEmail: "image-updater@example.com",
	}, nil
}

// CloneRepository clones a Git repository
func (c *Client) CloneRepository(ctx context.Context, url, branch string) (string, error) {
	log.Infof("Cloning repository %s (branch: %s)", url, branch)

	// Create a unique directory name based on the repository URL
	repoDir := filepath.Join(c.cloneDir, fmt.Sprintf("repo-%d", time.Now().UnixNano()))

	// Clone options
	options := &git.CloneOptions{
		URL:           url,
		Auth:          c.auth,
		Progress:      nil,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Depth:         1,
	}

	// Clone the repository
	_, err := git.PlainClone(repoDir, false, options)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return repoDir, nil
}

// GetFile gets a file from a Git repository
func (c *Client) GetFile(ctx context.Context, repoDir, filePath string) (*model.File, error) {
	log.Infof("Getting file %s from repository at %s", filePath, repoDir)

	// Open the repository
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	// Get the file tree
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	// Get the file
	file, err := tree.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Read the file content
	content, err := file.Contents()
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Get the last commit for this file
	lastCommit, err := c.getLastCommitForFile(repo, filePath)
	if err != nil {
		log.Warnf("Failed to get last commit for file: %v", err)
	}

	return &model.File{
		Path:       filePath,
		Type:       "file",
		LastCommit: lastCommit.Hash.String(),
		LastUpdate: lastCommit.Author.When.Format(time.RFC3339),
		Content:    content,
	}, nil
}

// UpdateFile updates a file in a Git repository
func (c *Client) UpdateFile(ctx context.Context, repoDir, filePath, content, commitMessage string) error {
	log.Infof("Updating file %s in repository at %s", filePath, repoDir)

	// Open the repository
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create the full file path
	fullPath := filepath.Join(repoDir, filePath)

	// Ensure the directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file content
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Add the file to the index
	if _, err := worktree.Add(filePath); err != nil {
		return fmt.Errorf("failed to add file to index: %w", err)
	}

	// Commit the changes
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  c.defaultUser,
			Email: c.defaultEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Push the changes
	if err := repo.Push(&git.PushOptions{
		Auth: c.auth,
	}); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}

// getLastCommitForFile gets the last commit for a file
func (c *Client) getLastCommitForFile(repo *git.Repository, filePath string) (*object.Commit, error) {
	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get the commit history
	logOptions := &git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
		PathFilter: func(path string) bool {
			return path == filePath
		},
	}

	// Get the commit iterator
	commits, err := repo.Log(logOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	// Get the first commit (most recent)
	commit, err := commits.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get next commit: %w", err)
	}

	return commit, nil
}
