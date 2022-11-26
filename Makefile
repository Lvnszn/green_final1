.phony: build docker zip
build:
	rm -rf server
	GOOS=linux GOARCH=amd64 /usr/local/Cellar/go/1.18.2/libexec/bin/go build -o server

docker:
	docker build -t green3 .

zip:
	docker save green3 -o ~/Downloads/green3.tar

run:
	docker run --name green3 --rm -it green3 /bin/bash
