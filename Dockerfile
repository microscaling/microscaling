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
LABEL org.label-schema.vendor="Microscaling Systems" \
      org.label-schema.license="Apache-2.0" \
      org.label-schema.url="https://microscaling.com" \
      org.label-schema.vcs-type="git" \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.docker.dockerfile="/Dockerfile"

CMD ["/microscaling"]
