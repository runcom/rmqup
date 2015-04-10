all: test build

get:
	go get -d -v

install: get
	go install -v

test: install
	go test -a

build:
	go build -a -o rmqup

shell:
	docker build -t rmqup-dev .
	docker run -ti --rm -v `pwd`:/go/src/github.com/runcom/rmqup rmqup-dev /bin/bash

#deb: #fpm
