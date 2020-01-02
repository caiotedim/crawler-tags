FROM golang:1.12 as builder
WORKDIR /crawler-tags
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o crawler-tags

FROM debian:stable-slim
RUN apt update && apt upgrade -y
RUN apt install -y ca-certificates apt-transport-https && update-ca-certificates
COPY --from=builder /crawler-tags/crawler-tags .
CMD ["/crawler-tags"]