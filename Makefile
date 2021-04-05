test:
	go build -v -o dist/locksmith main.go
	./scripts/generate_test_pki.bundle.sh

build:
	go build -v -o dist/locksmith main.go

run:
	go run main.go

test-bundle:
	./scripts/generate_test_pki.bundle.sh

test-compare:
	./scripts/generate_test_pki.compare.sh