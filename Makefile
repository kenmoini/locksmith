test:
	go build -v -o dist/locksmith main.go
	./scripts/generate_test_pki.bundle.sh

build:
	go build -v -o dist/locksmith main.go

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -v -o dist/locksmith-linux-amd64 main.go

build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -v -o dist/locksmith-dawrwin-amd64 main.go

run:
	go run main.go

test-bundle:
	./scripts/generate_test_pki.bundle.sh

test-compare:
	./scripts/generate_test_pki.compare.sh
