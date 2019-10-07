FROM golang:latest
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go mod download
RUN go build -o main.go
EXPOSE 8080
CMD ["/app/main"]