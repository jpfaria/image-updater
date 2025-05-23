# Image Updater

A tool for managing Docker image deployments across different environments through GitOps.

## Overview

Image Updater monitors Docker registries for new image versions and provides a user interface to select which version to deploy to which environment. When a deployment is triggered, it automatically updates the corresponding manifest files in Git repositories, allowing ArgoCD or similar GitOps tools to detect the changes and apply them to the cluster.

## Features

- Monitor Docker registries for new image versions
- Web interface for selecting versions and environments
- Automatic updates of Git manifests (Helm values.yaml files)
- Support for multiple teams, applications, and environments
- Authentication and audit logging
- Webhook support for CI/CD integration

## Architecture

The application follows a clean architecture approach with:

- Echo web framework with Boost integration for the backend API
- Vue.js for the frontend interface
- Docker SDK for registry monitoring
- Git integration for manifest updates

## Development

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- Docker (for local testing)
- Git

### Getting Started

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Configure the application (see Configuration section)
4. Run `go run cmd/server/main.go` to start the server

## License

This project is licensed under the MIT License - see the LICENSE file for details.
