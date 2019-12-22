FROM golang:alpine as builder

RUN apk add git
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/apidummy
WORKDIR /go/src/apidummy

RUN dep init -v
RUN dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o main

FROM scratch
COPY --from=builder /go/src/apidummy/main /application/main
COPY --from=builder /go/src/apidummy/config.env /application/config.env
WORKDIR /application
CMD ["./main"]
