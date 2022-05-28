FROM golang:1.18-alpine

ENV MYCDK_HOME /go/src/mycdk
WORKDIR $MYCDK_HOME
ADD .. $MYCDK_HOME

RUN apk update && apk add git openssh gcc libc-dev make aws-cli

RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest && go install -v golang.org/x/tools/gopls@latest
