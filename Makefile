LDFLAGS :=  -w -s

build:
	go build -ldflags "$(LDFLAGS)"

buildForLinux:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)"