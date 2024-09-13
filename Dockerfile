# Use the golang container and give it a friendly alias using AS
FROM golang:1.21.0 AS gobuild
ENV GO111MODULE=off
# Set our working directory
WORKDIR /go/src/bjss.com/academy/hello

# Copy the source file
COPY app.go .

# build the app
RUN go build -a -o app .

# Stage 2
FROM alpine:3.17

WORKDIR /app/

# Copy the app from the first container
COPY --from=gobuild /go/src/bjss.com/academy/hello/app .
CMD ["./app"]
 