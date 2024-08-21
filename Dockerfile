FROM golang:1.22.4 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp2 .

FROM alpine:latest

WORKDIR /root/

RUN mkdir -p /root/logs
RUN touch /root/logs/app.log
    
COPY --from=builder /app/myapp2 .
COPY --from=builder /app/config/model.conf ./config/
COPY --from=builder /app/.env .

EXPOSE 8081


CMD [ "./myapp2" ]