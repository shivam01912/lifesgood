FROM golang:1.19

# Set destination for COPY
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build
RUN go mod -o /docker-lifes-good

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/docker-lifes-good"]