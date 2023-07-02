.PHONY: go-fix-lint
go-fix-lint:
	find . -print | grep --regex '.*\.go$$' | xargs goimports -w -local "github.com/kyosu-1/reversi/server"

.PHONY: go-check-lint
go-check-lint:
	golangci-lint run

.PHONY: go-test
go-test:
	go test -race -cover -count=1 ./...
