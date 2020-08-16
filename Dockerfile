FROM golang:1.15 as builder

WORKDIR /go/src/app

COPY . .

RUN cd cmd && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main

FROM scratch

COPY --from=builder /go/src/app/cmd/main /app/main

EXPOSE 8080

CMD ["/app/main"]