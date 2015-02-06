PACKAGES=utils timeserver github.com/cihub/seelog

GOPATH=$(CURDIR)
GODOC_PORT=:6060

all: fmt install

install:
	GOPATH=$(GOPATH) go install $(PACKAGES)

fmt:
	GOPATH=$(GOPATH) go fmt $(PACKAGES)

doc:
	GOPATH=$(GOPATH) godoc -v --http=$(GODOC_PORT) --index=true

clean:
	rm -rf bin pkg

run: install
	bin/timeserver

version: install
	bin/timeserver -v
