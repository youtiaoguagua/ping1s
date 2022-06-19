LDFLAGS :=  -w -s

build:
	go build -ldflags "$(LDFLAGS)"

