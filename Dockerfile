FROM golang:latest
LABEL maintainer="spamfree@matthieubessat.fr"

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o jobatator-cli main.go
RUN go test -count=1 -v ./pkg/connexion

ENTRYPOINT ["./jobatator-cli"]
CMD []
