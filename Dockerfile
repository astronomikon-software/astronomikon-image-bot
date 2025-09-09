# syntax=docker/dockerfile:1

FROM golang:1.25.0

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY src/go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY src/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /image-uploader-bot


# Run
CMD ["/image-uploader-bot"]