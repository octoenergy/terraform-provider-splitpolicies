default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	golangci-lint run
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
