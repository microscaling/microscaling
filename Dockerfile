FROM alpine:3.4
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
ARG BUILD_DATE
ARG VERSION
ARG VCS_URL
ARG VCS_REF

# Metadata
LABEL org.label-schema.name="Microscaling Engine" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.version=$VERSION \
      org.label-schema.url="https://microscaling.com" \
      org.label-schema.vendor="Microscaling Systems" \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.schema-version="1.0" \
      org.label-schema.docker.dockerfile="/Dockerfile" \
      org.label-schema.license="Apache-2.0" \
      org.label-schema.description="Our Microscaling Engine provides automation, resilience and efficiency for microservice architectures. Experiment with microscaling at app.microscaling.com."

CMD ["/microscaling"]
