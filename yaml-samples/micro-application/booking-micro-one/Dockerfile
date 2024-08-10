FROM golang:buster

WORKDIR /app
ADD . .
RUN go build -o /usr/local/bin/hello-booking

EXPOSE 4000
CMD ["/usr/local/bin/hello-booking"]
