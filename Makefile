build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o output/sidecar-mac-amd64 ./cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o output/sidecar-mac-arm64 ./cmd
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/sidecar-linux-amd64 ./cmd
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o output/sidecar-linux-arm64 ./cmd
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o output/sidecar-win-amd64.exe ./cmd
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o output/sidecar-win-arm64.exe ./cmd
