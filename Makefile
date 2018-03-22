
BUILD_ARGS= go build -o scheduler github.com/mobile-health/scheduler-service/src


build:
	$(BUILD_ARGS)

build-linux:
	env GOOS=linux GOARCH=amd64 $(BUILD_ARGS)

get_realize:
	go get github.com/tockins/realize

run: 
	realize start --path="src" --run --no-config

test:
	go test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=650s ./src/services
	go test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=650s ./src/stores
	go test $(GOFLAGS) -run=$(TESTS) -test.v -test.timeout=650s ./src/api1