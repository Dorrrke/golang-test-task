FROM golang:1.16-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o test-task-app ./cmd/main.go

CMD [ "./test-task-app" ]