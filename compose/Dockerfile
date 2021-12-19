# https://blog.jetbrains.com/go/2020/05/08/running-go-applications-using-docker-compose-in-goland/

# Compile stage
FROM golang:1.17.5-bullseye AS build-env

# Build Delve
RUN go get github.com/go-delve/delve/cmd/dlv

ADD . /dockerdev
WORKDIR /dockerdev

# Compile the application with the optimizations turned off
# This is important for the debugger to correctly work with the binary
RUN go build -gcflags "all=-N -l" -o /server ./cmd/server

# Final stage
FROM debian:bullseye

EXPOSE 8000 40000

WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server /

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]