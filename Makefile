


build:
	go build -o go_api_boilerplate github.com/mobile-health/go-api-boilerplate/src
get_realize:
	go get github.com/tockins/realize
run: 
	realize start --path="src" --run --no-config
test:
	go test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=650s ./src/services