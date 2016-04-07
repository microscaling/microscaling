FROM alpine:3.2
MAINTAINER Ross Fairbanks "ross@microscaling.com"

ENV BUILD_PACKAGES bash curl-dev

RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    rm -rf /var/cache/apk/*

# needs to be built for Linux:
# GOOS=linux go build -o microscaling .
ADD microscaling /

ADD run.sh /
RUN chmod +x /run.sh

# Needs a run.sh wrapper to run the microscaling binary successfully
ENTRYPOINT ["/run.sh"]
