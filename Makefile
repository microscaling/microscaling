# Image can be overidden with an env var.
DOCKER_IMAGE?=microscaling/microscaling

# Version tag from the code.
VERSION=`cat VERSION`

default: build

build:
	# Compile for Linux
	GOOS=linux go build -o microscaling

	# Build Docker image
	docker build \
  --build-arg VERSION=$(VERSION) \
  --build-arg VCS_URL=`git config --get remote.origin.url` \
  --build-arg VCS_REF=`git rev-parse --short HEAD` \
  --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
  -t $(DOCKER_IMAGE):$(VERSION) .
