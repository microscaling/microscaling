FROM alpine:3.3
MAINTAINER Ross Fairbanks "ross@microscaling.com"

ENV BUILD_PACKAGES ca-certificates

RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    rm -rf /var/cache/apk/*

# Add binary and Dockerfile
COPY microscaling Dockerfile /

RUN chmod +x /microscaling

# Metadata params
ARG VERSION
ARG VCS_URL
ARG VCS_REF
ARG BUILD_DATE

# Metadata
LABEL com.microscaling.vendor="Microscaling Systems" \
      com.microscaling.license="Apache-2.0" \
      com.microscaling.url="https://microscaling.com" \
      com.microscaling.vcs-type="git" \
      com.microscaling.vcs-url=$VCS_URL \
      com.microscaling.vcs-ref=$VCS_REF \
      com.microscaling.build-date=$BUILD_DATE \
      com.microscaling.dockerfile="/Dockerfile"

CMD ["/microscaling"]
