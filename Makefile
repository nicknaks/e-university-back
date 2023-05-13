
gg:: # Перегенерирует GQL схему сервиса
	@echo "\n --- 🧬 GraphQL generation --- \n"
	printf '// +build tools\npackage tools\nimport _ "github.com/99designs/gqlgen"' | gofmt > tools.go
	go mod tidy
	go run github.com/99designs/gqlgen generate

run::
	go run server.go


.PHONY: db
db:
	psql -c "drop database if exists e_university;"
	createdb e_university
	goose -allow-missing -dir migrations postgres "dbname=e_university sslmode=disable" up

.PHONY: db2
db2:
	psql -c "drop database if exists e_university;"
	createdb e_university
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=newPassword dbname=e_university sslmode=disable" goose -dir migrations up
