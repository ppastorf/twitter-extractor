FROM golang:alpine3.13 AS builder

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/app .

FROM alpine:latest

COPY --from=builder /go/bin/app /worker-extractor
CMD ["/worker-extractor"]
