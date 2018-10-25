FROM golang:1.10.0-alpine3.7

WORKDIR /go/src/github.com/mike-douglas/chaosproxy

COPY . .

RUN apk update && apk add git
RUN wget -O -  https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build -o proxy .

VOLUME ["/chaosproxy.yaml", "/server.pem", "/server.key"]
EXPOSE 8080

ENTRYPOINT [ "./proxy", "-config", "/chaosproxy.yaml", "-port", "8080"]