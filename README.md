# TimeServer

## Latest Version: 2.0.0

## Makefile commands
### Install the package
```
	make install
```
### Run timeserver on default port 8080
```
	make run
```
### Run timeserver on a different port
```
	make run FLAGS='--port 8081'
```
### Run timeserver with different logging configuration
```
	make run FLAGS='--log etc/my-log.xml'
```
### Run timeserver with different set of templates
```
	make run FLAGS='--templates new-template-folder'
```
### Check timeserver version
```
	make run FLAGS='-v'
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
1. Start timeserver using command `make run`
2. Access timeserver using the following URL: [http://localhost:8080/]
3. Login timeserver: put your name in the input text box and click "submit" button
4. Check current time using the following URL: [http://localhost:8080/time]
5. Logout timeserver using the following URL: [http://localhost:8080/logout]
6. Logs can be found in out/timeserver.log   