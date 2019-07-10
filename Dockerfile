FROM golang:1.12-alpine3.9

WORKDIR /go/src/github.com/mike-douglas/chaosproxy

COPY . .

RUN apk update && apk add git
RUN wget -O -  https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build -o proxy .

VOLUME ["/chaosproxy.yaml"]
EXPOSE 8080

CMD ["./proxy", "-config", "/chaosproxy.yaml"]