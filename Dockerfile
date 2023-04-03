FROM golang:1.20.2-alpine3.17 AS builder

# Change work directory to one that we are going to use as our porject directory
WORKDIR /usr/src/app

# Download all the dependencies specified in the "go.mod" file of our module.
# This command downloads dependency modules into a kind of local dependency cache, directory "~/go/pkg/mod/"
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy all files from current directory (context for the building of image)
# to the current WORKDIR in image
COPY . .

# Build the web application
RUN go build -o /usr/bin/app/webapp /usr/src/app

FROM alpine:3.17

# Copy the binary file from the first image to the second image
COPY --from=builder /usr/bin/app/webapp /usr/bin/app/webapp

# Launch the web application
CMD ["/usr/bin/app/webapp"]