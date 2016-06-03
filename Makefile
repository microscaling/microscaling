default: test

# Build Docker image
build: docker_build output

# Build and push Docker image
release: docker_build docker_push output

# Image can be overidden with an env var.
DOCKER_IMAGE ?= microscaling/microscaling

# Get the latest commit.
GIT_COMMIT = `git rev-parse --short HEAD` 

ifeq ($(MAKECMDGOALS),release)
# Use the version number as the release tag.
DOCKER_TAG = `cat VERSION`
else
# Add the commit ref for development builds.
DOCKER_TAG = `cat VERSION`-$(GIT_COMMIT)
endif

test:
	go test -v ./...

get-deps:
	go get -t -v ./...

docker_build:
	# Compile for Linux
	GOOS=linux go build -o microscaling

	# Build Docker image
	docker build \
  --build-arg VCS_URL=`git config --get remote.origin.url` \
  --build-arg VCS_REF=$(GIT_COMMIT) \
  --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	-t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker_push:
	# Tag image as latest
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

	# Push to DockerHub
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

output:
	@echo Docker Image: $(DOCKER_IMAGE):$(DOCKER_TAG)
