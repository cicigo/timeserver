default: build

build:
	go build -o bin/timeserver timeserver.go

fmt:
	go fmt timeserver.go

doc:
	godoc -http=:6060 -index

run: build
	bin/timeserver

clean: 
	rm -rf bin

version: build
	bin/timeserver -v
