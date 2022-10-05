FROM golang

WORKDIR /app

COPY . .

CMD ["/app/main"]
