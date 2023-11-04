ARG GO_VERSION=1.21.3
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache git
WORKDIR /build_app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go test ./... -v
# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /helsinki-guide main.go
 
# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final

USER nonroot:nonroot
# copy compiled app
COPY --from=build --chown=nonroot:nonroot /helsinki-guide /helsinki-guide
ENTRYPOINT ["/helsinki-guide", "bot"]