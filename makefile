run-rsi:
	go run services/rsi/rsi.go -f services/rsi/etc/rsi.yaml
run-price:
	go run services/price/price.go -f services/price/etc/price.yaml

run-order:
	go run services/order/order.go -f services/order/etc/order.yaml
run-position:
	go run services/position/position.go -f services/position/etc/position.yaml
update:
	git submodule update --remote

gen-proto: #update
	goctl rpc protoc ./proto/$(name)/$(name).proto --go_out=./pb/ --go-grpc_out=./pb/ --zrpc_out=./services/$(name)/  --style=go_zero