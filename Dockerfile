# Stage 1: Build the Go application
# Using golang:1.24-alpine as requested for the builder stage.
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container for the builder stage
WORKDIR /app

# Copy go.mod and go.sum files to the working directory.
# This allows Go to download dependencies before copying the rest of the code,
# which optimizes Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source code into the /app directory in the container.
# This includes all your .go files.
COPY . .

# Build the application.
# CGO_ENABLED=0 is important for creating statically linked binaries, reducing image size.
# -ldflags "-s -w" reduces the binary size further by omitting debug info.
# -o consistent_1 specifies the output executable name.
# ./Delivery tells Go to build the main package found in the 'Delivery' directory.
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

# Stage 2: Create the final lean production image
# Use a small base image like alpine for the final production container.
FROM alpine:latest

# Set the working directory for the final image. This is where your app will run from.
# It's important to keep this consistent with where your app expects to find its files (e.g., .env).
WORKDIR /root/

# Copy the built executable from the builder stage into the final image's working directory.
# Ensure 'consistent_1' is the exact name of your compiled binary.
COPY --from=builder /app/consistent_1 .

# --- START OF FIREBASE SERVICE ACCOUNT AND .ENV FIXES ---

# Create the 'firebase-service-account.json' file using the content from a Render environment variable.
# The 'FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT' variable MUST be set on Render with the full JSON string.
# This makes the file available at './firebase-service-account.json' in the container's /root/ directory.
RUN echo "$FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT" > ./firebase-service-account.json

# Create a minimal '.env' file. Your Go application expects this file to exist due to `viper.ReadInConfig()`.
# We specifically set 'FIREBASE_SERVICE_ACCOUNT_PATH' here to point to the file we just created.
# Other variables (like MONGO_URI, JWT_SECRET, etc.) will be automatically overridden by Render's
# environment variables due to `viper.AutomaticEnv()` in your Go code, so they don't need to be
# explicitly set in this .env file unless your app has a hard requirement for them to be in the file.
RUN echo "FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json" > ./.env

# --- END OF FIREBASE SERVICE ACCOUNT AND .ENV FIXES ---

# Expose the port your application listens on. Render will use this.
EXPOSE 8080

# Command to run the executable when the container starts.
# Make sure 'consistent_1' is the exact name of your compiled binary.
CMD ["./consistent_1"]