FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o go-app
EXPOSE 8080
CMD ["/app/go-app"]
