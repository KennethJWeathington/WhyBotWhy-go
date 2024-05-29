.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build: vet
	go build
build-windows: vet
	env GOOS=windows GOARCH=amd64 go build
cleanup:
	rm whybotwhy_go; rm whybotwhy_go.exe