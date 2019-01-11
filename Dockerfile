FROM golang:1.11 as builder
ENV CGO_ENABLED 0
COPY . resource
WORKDIR /go/resource
RUN go build -o /assets/in github.com/antonu17/concourse-docker-image-tags/cmd/in \
 && go build -o /assets/check github.com/antonu17/concourse-docker-image-tags/cmd/check

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder assets/ /opt/resource/
COPY assets/out /opt/resource/
