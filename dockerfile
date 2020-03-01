# docker build -t kkond/hw1_tree .
FROM golang:1.9.2
COPY . .
RUN go run main.go .