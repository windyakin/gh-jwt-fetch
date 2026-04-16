FROM mirror.gcr.io/library/golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux \
  go build -ldflags="-s -w" -o /gh-jwt-fetch .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /gh-jwt-fetch /gh-jwt-fetch

ENTRYPOINT ["/gh-jwt-fetch"]
