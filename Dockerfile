ARG opts
# Build the logapi binary
FROM golang:alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests

COPY main.go main.go

# Build
ARG CGO_ENABLED=0
ARG GOARCH=amd64
ARG GOARM=6
RUN env CGO_ENABLED=${CGO_ENABLED} GOARCH=${GOARCH} GOARM=${GOARM}  go build -o app ./main.go

# Use distroless as minimal base image to package the app binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM scratch
WORKDIR /
COPY --from=builder /workspace/app .

ENTRYPOINT ["/app"]