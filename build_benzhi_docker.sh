#!/bin/bash
set -e
IMAGE_NAME=${1:-my-go-task}
PLATFORM=${2:-linux/amd64}
docker buildx build --platform "$PLATFORM" -f benzhi.Dockerfile -t "$IMAGE_NAME" .
echo ""
echo "✅ Docker image '$IMAGE_NAME' built successfully!"
