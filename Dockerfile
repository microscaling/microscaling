FROM alpine:3.2
MAINTAINER Ross Fairbanks "ross@force12.io"

ENV BUILD_PACKAGES bash curl-dev
RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    rm -rf /var/cache/apk/*

# force12 needs to be built for Linux:
# GOOS=linux go build -o force12 .
ADD force12 /

ADD compose-demo.yml /
ADD run.sh /
RUN chmod +x /run.sh

# Needs a run.sh wrapper to run the force12 binary successfully
ENTRYPOINT ["/run.sh"]
