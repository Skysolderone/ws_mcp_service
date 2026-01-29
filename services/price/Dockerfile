from golang:1.25.6-alpine as builder

workdir /app 


copy go.mod go.sum ./

run go mod download


copy .
run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o  price ./price.go


from alpine:latest

workdir /app
copy --from=builder /app/price /app/price

expose 8080

cmd ["/app/price", "-f", "/app/etc/price.yaml"]