all: test vet lint staticcheck

test:
	go test   -v .
	go test   -v ./consumers/
	go test  -v ./producers/

vet:
	go vet ./...

install-lint:
	go get golang.org/x/lint/golint
	go list -f {{.Target}} golang.org/x/lint/golint

lint:
	go list ./... | grep -v /msg/ | xargs -L1 golint -set_exit_status

install-staticcheck:
	cd /tmp && GOPROXY="" go get honnef.co/go/tools/cmd/staticcheck

staticcheck:
	staticcheck ./...