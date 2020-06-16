FROM alpine as alpine
RUN apk add -U --no-cache ca-certificates

FROM golang as builder
RUN mkdir -p /app/no/
WORKDIR /app/no/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/no/no /app/
COPY marks/box.png /app/marks/box.png
WORKDIR /app
ENTRYPOINT ["./no"]
