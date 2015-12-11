FROM golang:1.5.2


WORKDIR /go

ADD . /go/src/github.com/cytoscape-ci/service-go

WORKDIR /go/src/github.com/cytoscape-ci/service-go

RUN go get github.com/rs/cors

RUN go build app.go

EXPOSE 3000

CMD ./app
