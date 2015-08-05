FROM golang:1.3-onbuild
MAINTAINER Ross Fairbanks "ross.fairbanks@gmail.com"

ADD dynamo-config.json /etc/aws-config.json

CMD ["go", "run", "force12.go"]
