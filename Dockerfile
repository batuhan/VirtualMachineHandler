FROM golang:1.13-alpine AS builder
RUN apk update && apk add --no-cache curl

RUN curl -L https://github.com/vmware/govmomi/releases/download/v0.21.0/govc_linux_amd64.gz | gunzip > /govc
RUN chmod +x /govc

RUN mkdir /app
WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o output

FROM golang:1.13-alpine

COPY --from=builder /app/output /app
COPY --from=builder /govc /usr/local/bin/govc
ENTRYPOINT ["/app"]
