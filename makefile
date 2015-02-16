PACKAGES=utils command/authserver github.com/cihub/seelog

GOPATH=$(CURDIR)
GODOC_PORT=:6060
FLAGS=

all: fmt install

install:
	GOPATH=$(GOPATH) go install $(PACKAGES)

fmt:
	GOPATH=$(GOPATH) go fmt $(PACKAGES)

doc:
	GOPATH=$(GOPATH) godoc -v --http=$(GODOC_PORT) --index=true

clean:
	rm -rf bin pkg out

timeserver: install
	bin/timeserver $(FLAGS)

authserver: install
	bin/authserver $(FLAGS)
