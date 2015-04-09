all: test build

get:
	go get -d -v

install: get
	go install -v

test: install
	go test -a

build:
	go build -a -o up

#deb: #fpm
