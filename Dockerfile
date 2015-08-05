FROM golang:1.3-onbuild
MAINTAINER Ross Fairbanks "ross.fairbanks@gmail.com"

CMD ["go", "run", "*.go"]
