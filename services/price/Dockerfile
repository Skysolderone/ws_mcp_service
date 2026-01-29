from golang:1.25.6-alpine as builder

workdir /app

# 只拷贝构建 price 服务需要的代码，而不是整个仓库
copy go.mod go.sum ./
run go mod download

copy internal/ ./internal/
copy pb/ ./pb/
copy pkg/ ./pkg/
copy services/price/ ./services/price/

run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o price ./services/price/price.go


from alpine:latest

workdir /app
copy --from=builder /app/price /app/price

expose 8080

cmd ["/app/price", "-f", "/app/etc/price.yaml"]