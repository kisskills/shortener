.PHONY: test-cover gen-grpc inmemory postgres

gen-grpc:
	protoc -I ./api/proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. ./api/proto/shortener/shortener.proto
inmemory:
	docker-compose --profile inmemory up -d --build
postgres:
	docker-compose --profile postgres up -d --build
clean:
	docker-compose --profile postgres --profile inmemory down --volumes
test-cover:
	go test -race ./... -tags unit -coverprofile .testCoverage.out
	go tool cover -html=.testCoverage.out
	rm .testCoverage.out