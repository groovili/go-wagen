FROM golang:alpine

RUN mkdir /go/src/{{.Name}}

WORKDIR /go/src/{{.Name}}

COPY . ./

RUN go env

RUN go build -o ./cmd/{{.Name}}/{{.Name}} ./cmd/{{.Name}}/{{.Name}}.go

CMD ./cmd/{{.Name}}/{{.Name}}