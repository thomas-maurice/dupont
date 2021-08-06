all: binary

bin:
	if ! [ -d bin ]; then mkdir bin; fi

binary: bin
	go build -o bin/dupont 