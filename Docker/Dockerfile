# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/Mobisense/virtproject

# Resolving dependencies
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/julienschmidt/httprouter

# Build the golang-docker command inside the container.
RUN go install Mobisense/virtproject

# Run the golang-docker command when the container starts.
ENTRYPOINT /go/bin/virtproject

# http server listens on port 8080.
EXPOSE 9999