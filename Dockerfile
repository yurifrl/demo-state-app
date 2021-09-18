FROM golang:alpine as builder

ARG TARGETPLATFORM
ENV TARGETPLATFORM=${TARGETPLATFORM:-linux/amd64}

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /

COPY main.go main.go

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) \
    && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) \
    && \
    GOARM=$(echo ${TARGETPLATFORM} | cut -d / -f3); export GOARM=${GOARM:1} \
    && \
    go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o app ./main.go

# Use distroless as minimal base image to package the app binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM scratch
LABEL org.opencontainers.image.source https://github.com/yurifrl/demo-state-app
WORKDIR /
COPY --from=builder /app .

ENTRYPOINT ["/app"]