# TimeServer

## Latest Version: 2.0.0

## Makefile commands
### Install the package
```
	make install
```
### Run timeserver on default port 8080
```
	make timeserver
```
### Run timeserver on a different port
```
	make timeserver FLAGS='--port 8081'
```
### Run timeserver with different logging configuration
```
	make timeserver FLAGS='--log etc/my-log.xml'
```
### Run timeserver with different set of templates
```
	make timeserver FLAGS='--templates new-template-folder'
```

### Run timeserver with authhost 
```
	make timeserver FLAGS="--authhost http://localhost --max-inflight 1"
```

### Check timeserver version
```
	make timeserver FLAGS='-v'
```

### Run authserver on default port 7070
```
	make authserver
```
    
### Run authserver on default port 7070 with different checkpoint interval and dump file location
```
	make authserver FLAGS="--checkpoint-interval 5 --dumpfile out/auth.json"
```

### Clean build
```
	make clean
```
### Format source code
```
	make fmt
```

## How to use timeserver
1. Start timeserver using command `make timeserver FLAGS="--authhost http://localhost --max-inflight 100"`
1. Start authserver using command `make authserver`
2. Access timeserver using the following URL: [http://localhost:8080/]
3. Login timeserver: put your name in the input text box and click "submit" button
4. Check current time using the following URL: [http://localhost:8080/time]
5. Logout timeserver using the following URL: [http://localhost:8080/logout]
6. Logs can be found in out/timeserver.log   