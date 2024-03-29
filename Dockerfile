# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder


# Set destination for COPY
WORKDIR /quizard-backend-go

COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code in this directory
COPY . .

# Build executable with name quizard-backend into this directory
RUN CGO_ENABLED=0 GOOS=linux go build -o quizard-backend .


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
# would like to do from scratch but got nonroot user errors docker restart may fix but it interrupts service
FROM alpine:edge

WORKDIR /quizard-backend-go

# copy exe from builder into this empty image at this directory. exe will have name quizard-backend
COPY --from=builder /quizard-backend-go/quizard-backend .

# EXPOSE 8080

# enter into workdir of this image, slash executable name
# docker build -t quizard-be:multistage -f Dockerfile.multistage .
# docer run --publish 8080:8080 quizard-be:multistage
ENTRYPOINT [ "/quizard-backend-go/quizard-backend" ]


