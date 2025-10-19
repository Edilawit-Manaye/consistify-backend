# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

# FROM alpine:latest

# WORKDIR /root/

# COPY --from=builder /app/consistent_1 .

# RUN echo "$FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT" > ./firebase-service-account.json

# RUN echo "FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json\nFIREBASE_PROJECT_ID=${FIREBASE_PROJECT_ID}" > ./.env

# EXPOSE 8080

# CMD ["./consistent_1"]










# # Stage 1: Build the Go application
# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

# # Stage 2: Create the final lean production image
# FROM alpine:latest

# WORKDIR /root/

# COPY --from=builder /app/consistent_1 .

# # --- SECURE FIREBASE SERVICE ACCOUNT HANDLING ---
# # Create a temporary .env file for FIREBASE_SERVICE_ACCOUNT_PATH
# # This file will tell your Go app where to find the service account JSON
# RUN echo "FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json" > ./.env

# # Create the firebase-service-account.json file from the environment variable
# # This reads the *content* of the JSON from an environment variable (set on Render)
# # and writes it to the file your app expects.
# RUN apk add --no-cache bash # Install bash for 'printenv' command, if not already present in alpine
# RUN printenv FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT > firebase-service-account.json
# # --- END SECURE FIREBASE HANDLING ---

# EXPOSE 8080

# CMD ["./consistent_1"]



# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

# Stage 2: Create the final lean production image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/consistent_1 .

# Install bash (needed if using bash syntax)
RUN apk add --no-cache bash

# --- SECURE FIREBASE SERVICE ACCOUNT HANDLING ---
# Create .env file for your Go app
RUN echo "FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json" > .env

# At runtime, create the service account JSON file from env variable
ENTRYPOINT ["/bin/sh", "-c", "echo \"$FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT\" > firebase-service-account.json && ./consistent_1"]
# --- END SECURE FIREBASE HANDLING ---

EXPOSE 8080







