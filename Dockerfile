# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

ADD . /go/src/beers-api
WORKDIR /go/src/beers-api

RUN cd /go/src/beers-api

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go install cmd/api

EXPOSE 8080
ENTRYPOINT ["/go/src/beers-api"]

CMD ["/beers-api"]

RUN echo "Server running on port 8080"