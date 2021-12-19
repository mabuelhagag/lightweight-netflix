# Pull base image
FROM golang:alpine

# Install git
RUN apk update && apk add --no-cache git build-base

# Where our files will be in the docker container 
WORKDIR /opt/go-app

# Get Dependency
COPY ./go.mod ./go.sum ./
RUN go mod download -x

# Install Air for hot reload
RUN go install github.com/cosmtrek/air@latest
# Insall Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy the source from the current directory to the working Directory inside the container
# Source also contains go.mod and go.sum which are dependency files
COPY . .

# The ENTRYPOINT defines the command that will be ran when the container starts up
# In this case air command for hot reload go apps on file changes
ENTRYPOINT air