REGISTRY_NAME?=docker.io/iejalapeno
IMAGE_VERSION?=latest

.PHONY: all telemetry-processor container push clean test

ifdef V
TESTARGS = -v -args -alsologtostderr -v 5
else
TESTARGS =
endif

all: telemetry-processor

telemetry-processor:
	mkdir -p bin
	$(MAKE) -C ./cmd compile-telemetry-processor

telemetry-processor-container: telemetry-processor
	docker build -t $(REGISTRY_NAME)/telemetry-processor:$(IMAGE_VERSION) -f ./build/Dockerfile.telemetry-processor .

push: telemetry-processor-container
	docker push $(REGISTRY_NAME)/telemetry-processor:$(IMAGE_VERSION)

clean:
	rm -rf bin

test:
	GO111MODULE=on go test `go list ./... | grep -v 'vendor'` $(TESTARGS)
	GO111MODULE=on go vet `go list ./... | grep -v vendor`
