# Compile image
FROM golang:1.14-alpine as builder

RUN apk --update-cache add upx

# Setup build directory
RUN mkdir -p /build
WORKDIR /build

# Go build options
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Fetch deps first as they don't change frequently
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
RUN go mod download

# Copy go code
COPY version/     /build/version
COPY pkg/         /build/pkg
COPY cmd/manager  /build/cmd/manager

# Build the binary
RUN go build \
  -o /tmp/vdi-manager \
  ./cmd/manager && upx /tmp/vdi-manager

##
# Build the runtime image
##
FROM scratch

# Install operator binary
COPY --from=builder /tmp/vdi-manager /vdi-manager

ENTRYPOINT ["/vdi-manager"]
