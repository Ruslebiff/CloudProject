FROM golang:1.13.4

WORKDIR /go/src/CloudProject
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o main .

EXPOSE 8080
CMD ["cmd"]