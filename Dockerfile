FROM alpine:3.2
MAINTAINER Ross Fairbanks "ross.fairbanks@gmail.com"

# We'll want this if we can move to alpine (well, maybe not the curl thing if we use python / pip instead)
ENV BUILD_PACKAGES bash curl-dev python-dev build-base
ENV PYTHON_PACKAGES python py-pip
RUN apk update && \
    apk upgrade && \
    apk add $BUILD_PACKAGES && \
    apk add $PYTHON_PACKAGES && \
    rm -rf /var/cache/apk/*

#RUN apt-get update && \
#    apt-get -y install python-dev python-pip

RUN pip install -U docker-compose==1.4.2

# force12 needs to be built for Linux:
#   GOOS=linux go build -o force12 .
ADD force12 /
ADD windtunnel-compose.yml /

# Needs a run.sh wrapper to run the force12 binary successfully
ENTRYPOINT ["/bin/bash"]
ADD run.sh /

CMD ["/run.sh"]