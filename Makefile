all: test vet lint staticcheck

test:
	ls -lhat
	go test client_test.go
	go test consumers/consumer_test.go
	go test producers/producer_test.go

vet:
	go vet ./...

install-lint:
	go get golang.org/x/lint/golint
	go list -f {{.Target}} golang.org/x/lint/golint

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

install-staticcheck:
	cd /tmp && GOPROXY="" go get honnef.co/go/tools/cmd/staticcheck

staticcheck:
	staticcheck ./...