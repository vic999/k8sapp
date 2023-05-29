FROM scratch



#COPY certs /etc/ssl/certs/



FROM golang:1.20-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v


# Build the Go app
RUN go build -o ./out/k8sapp ./cmd/.

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

# This container exposes port 8080 to the outside world
ENV K8SAPP_LOCAL_HOST 0.0.0.0
ENV K8SAPP_LOCAL_PORT 8080
ENV K8SAPP_LOG_LEVEL 0

EXPOSE $K8SAPP_LOCAL_PORT

COPY --from=build_base /tmp/go-sample-app/out/k8sapp /k8sapp


CMD ["/k8sapp"]
