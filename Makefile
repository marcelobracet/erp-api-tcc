build:
	go build ./... 

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-watch:
	reflex -r '\.go$$' make test


	