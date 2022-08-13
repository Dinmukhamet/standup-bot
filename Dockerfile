FROM golang:1.18-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main main.go

FROM alpine:3.16
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Bishkek /etc/localtime
RUN echo "Asia/Bishkek" > /etc/timezone
COPY .env .
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]
