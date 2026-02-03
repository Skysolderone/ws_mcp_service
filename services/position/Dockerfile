from golang:1.25.6-alpine as builder
workdir /app
copy go.mod go.sum ./
run go mod download
copy internal/ ./internal/
copy model/ ./model/
copy pb/ ./pb/
copy pkg/ ./pkg/
copy services/position/ ./services/position/
run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o position ./services/position/position.go
from alpine:latest
workdir /app
copy --from=builder /app/position /app/position
copy services/position/etc/ /app/etc/
expose 8080
cmd ["/app/position", "-f", "/app/etc/position.yaml"]