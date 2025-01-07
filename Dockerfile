FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build .

FROM chromedp/headless-shell:latest

COPY --from=builder /app/og /og

ENTRYPOINT [ "/og" ]
