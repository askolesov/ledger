.PHONY: run

run:
	go run ./cmd/finances report data.yaml

s:
	go run ./cmd/finances report data.yaml -s