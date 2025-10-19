# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container for the builder stage
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
# This allows Go to download dependencies before copying the rest of the code,
# which optimizes Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source code into the /app directory in the container
COPY . .

# Build the application
# CGO_ENABLED=0 is important for creating statically linked binaries, reducing image size.
# -ldflags "-s -w" reduces the binary size further by omitting debug info.
# -o consistent_1 specifies the output executable name.
# ./Delivery tells Go to build the main package found in the 'Delivery' directory.
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

# Stage 2: Create the final lean production image
FROM alpine:latest

# Set the working directory for the final image. This is where your app will run from.
WORKDIR /root/

# Copy the built executable from the builder stage into the final image's working directory.
COPY --from=builder /app/consistent_1 .

# Copy the Firebase service account file from the builder stage.
# IMPORTANT: This assumes 'firebase-service-account.json' is in your project root.
COPY --from=builder /app/firebase-service-account.json .

# Expose the port your application listens on. Render will use this.
EXPOSE 8080

# Command to run the executable when the container starts.
CMD ["./consistent_1"]