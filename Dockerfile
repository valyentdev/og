FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build .

FROM chromedp/headless-shell:130.0.6723.59

COPY --from=builder /app/og /og

ENTRYPOINT [ "/og" ]
