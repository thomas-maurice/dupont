all: binary cross_compile

bin:
	if ! [ -d bin ]; then mkdir bin; fi

binary: bin
	go build -o bin/dupont 

# cross compile for raspberry pies and all
cross_compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o bin/dupont-linux-armv6
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o bin/dupont-linux-armv7
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dupont-linux-amd64