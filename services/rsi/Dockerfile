from golang:1.25.6-alpine as builder

workdir /app 


copy go.mod go.sum ./

run go mod download

copy . ./
run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o rsi ./services/rsi/rsi.go


from alpine:latest

workdir /app
copy --from=builder /app/rsi /app/rsi

expose 8080

cmd ["/app/rsi", "-f", "/app/etc/rsi.yaml"]