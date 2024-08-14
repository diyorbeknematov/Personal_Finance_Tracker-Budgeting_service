FROM golang:1.22.2 AS builder

WORKDIR /budgeting-service

COPY . .
RUN go mod download

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../budgeting .

FROM alpine:latest
WORKDIR /budgeting-service

COPY --from=builder /budgeting-service/budgeting .
COPY --from=builder /budgeting-service/pkg/logs/app.log ./pkg/logs/
COPY --from=builder /budgeting-service/.env .

EXPOSE 1111

CMD [ "./budgeting" ]