go-tool := "GOPROXY=file://$(go env GOMODCACHE)/cache/download go run github.com/kazhuravlev/toolset/cmd/toolset@v0.41.0 run"

test-all:
    @# @go test -race -timeout 5m -count=1 -coverprofile=coverage.out ./...
    @{{go-tool}} gotestsum --format pkgname-and-test-fails --format-icons default -- -race -timeout 5m -count=1 -coverprofile=coverage.out ./...
    @go tool cover -html=coverage.out -o coverage.html


test PACKAGE:
    @{{go-tool}} gotestsum --format testname --format-icons default -- -race -timeout 5m -count=1  ./{{PACKAGE}}/...
# go test -race -timeout 5m -count=1 -coverprofile=coverage.out -covermode=atomic -v -bench=. -benchmem ./...
# go test -race -timeout 5m ./... -tags=goleak

lint:
    @{{go-tool}} golangci-lint run
