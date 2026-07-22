.PHONY: test test-unit test-integration lint vet fmt

test:
	cd kernel && go test ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test ./...) || exit 1; \
	done

test-unit:
	cd kernel && go test -short ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test -short ./...) || exit 1; \
	done

test-integration:
	cd kernel && go test -tags=integration ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go test -tags=integration ./...) || exit 1; \
	done

lint:
	golangci-lint run ./...

vet:
	cd kernel && go vet ./...
	@for dir in $$(find protocols -name go.mod -exec dirname {} \;); do \
		(cd $$dir && go vet ./...) || exit 1; \
	done

fmt:
	go fmt ./...
