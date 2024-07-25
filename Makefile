APP_NAME := bloomberg-rss

linux:
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o bin/linux_amd64/$(APP_NAME) $(APP_NAME).go

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o bin/darwin_amd64/$(APP_NAME) $(APP_NAME).go
