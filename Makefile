binary=apimonitor.linux
linuxBinary=apimonitor.linux
version=200127
dockerImage=sangil/apimonitor-server

build:
	rm -f $(binary)
	go build -o ${binary} .

build-docker:
	rm -f $(linuxBinary)
	GOOS=linux GOARCH=amd64 go build -o $(linuxBinary) .
	docker build -f ./build/Dockerfile -t $(dockerImage):$(version) .
	docker push $(dockerImage):$(version)