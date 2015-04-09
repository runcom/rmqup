FROM golang:1.4.2

VOLUME /go

COPY . /go/src/github.com/runcom/rmqup

WORKDIR /go/src/github.com/runcom/rmqup

RUN make
RUN cp ./rmqup /usr/local/bin

ENTRYPOINT ["rmqup"]
