FROM golang:1.13

WORKDIR /go/src/eagraf/cloudworker
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["synchronizer"]