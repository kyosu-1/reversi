.PHONY: gen-api
gen-api:
	mkdir -p ./gen/api
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0
	oapi-codegen --config config/oapi-codegen/server.yaml ./api/openapi.yaml

.PHONY: terraform-fix-lint
terraform-fix-lint:
	for file in $$(find deployments/ -type f -name '*.tf'); do terraform fmt $$file; done

.PHONY: go-fix-lint
go-fix-lint:
	find . -print | grep --regex '.*\.go$$' | xargs goimports -w -local "github.com/kyosu-1/reversi"

.PHONY: go-check-lint
go-check-lint:
	golangci-lint run

.PHONY: go-test
go-test:
	go test -race -cover -count=1 ./...

.PHONY: __init-db-args
__init-db-args:
ifndef DB_HOST
	$(warning DB_HOST was not set; localhost is used)
	$(eval DB_HOST := localhost)
endif
ifndef DB_PORT
	$(warning DB_PORT was not set; 3306 is used)
	$(eval DB_PORT := 3306)
endif
ifndef DB_USER
	$(warning DB_USER was not set; root is used)
	$(eval DB_USER := root)
endif
ifndef DB_PASS
	$(warning DB_PASS was not set; passw0rd is used)
	$(eval DB_PASS := pass)
endif
ifndef DB_NAME
	$(warning DB_NAME was not set; lime is used)
	$(eval DB_NAME := reversi)
endif

.PHONY: db-migrate
db-migrate: __init-db-args
	go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	migrate -source "file://ddl" -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

.PHONY: gen-xo
gen-xo: __init-db-args
	mkdir -p ./gen/dbschema
	go install github.com/xo/xo@42b11c7999bc6ac5be620949723f44bd0ec63e02
	xo schema --out "gen/dbschema" -t json "mysql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"

.PHONY: gen-db
gen-db:
	mkdir -p ./gen/daocore
	go run ./script/dbgen ./gen/dbschema/xo.xo.json

.PHONY: seed
seed: __init-db-args
	for file in $$(find test/data -type f -not -name '*cleanup.sql' -name '*.sql' | sort); do $(MAKE) seedfile FILE=$$file; done

.PHONY: clean-db
clean-db: __init-db-args
	mysql -h$(HOST) -u$(USERNAME) -p$(PASS) -P$(PORT) --protocol=tcp $(DBNAME) < test/data/cleanup.sql