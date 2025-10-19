FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o consistent_1 ./Delivery

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/consistent_1 .

RUN echo "$FIREBASE_SERVICE_ACCOUNT_JSON_CONTENT" > ./firebase-service-account.json

RUN echo "FIREBASE_SERVICE_ACCOUNT_PATH=./firebase-service-account.json\nFIREBASE_PROJECT_ID=${FIREBASE_PROJECT_ID}" > ./.env

EXPOSE 8080

CMD ["./consistent_1"]