from golang:1.25.6-alpine as builder

workdir /app

# 只拷贝构建 rsi 服务需要的代码，而不是整个仓库
copy go.mod go.sum ./
run go mod download

copy internal/ ./internal/
copy model/ ./model/
copy pb/ ./pb/
copy pkg/ ./pkg/
copy services/order/ ./services/order/

run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o order ./services/order/order.go  && go clean -cache


from alpine:latest

workdir /app
copy --from=builder /app/order /app/order
copy services/order/etc/ /app/etc/

expose 8080

cmd ["/app/order", "-f", "/app/etc/order.yaml"]