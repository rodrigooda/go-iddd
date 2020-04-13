GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)

generate_proto:
	@protoc \
		-I service/customer/infrastructure/adapter/primary/grpc \
		-I /usr/local/include \
		-I $(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=plugins=grpc:service/customer/infrastructure/adapter/primary/grpc \
		--grpc-gateway_out=logtostderr=true:service/customer/infrastructure/adapter/primary/grpc \
		--swagger_out=logtostderr=true:service/customer/infrastructure/adapter/primary/grpc \
		service/customer/infrastructure/adapter/primary/grpc/customer.proto

generate_mocked_EventStore:
	@mockery \
		-name EventStore \
		-dir service/lib/es \
		-outpkg mocks \
		-output service/lib/eventstore/mocks \
		-note "+build test"

generate_mocked_ForAssertingUniqueEmailAddresses:
	@mockery \
		-name ForAssertingUniqueEmailAddresses \
		-dir service/customer/application/command \
		-outpkg mocks \
		-output service/customer/infrastructure/adapter/secondary/mocks \
		-note "+build test"

generate_all_mocks: \
	generate_mocked_EventStore \
	generate_mocked_ForAssertingUniqueEmailAddresses

lint:
	golangci-lint run --build-tags test ./...

# https://github.com/golangci/golangci-lint
install-golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.23.8


# https://github.com/psampaz/go-mod-outdated
outdated-list:
	go list -u -m -json all | go-mod-outdated -update -direct