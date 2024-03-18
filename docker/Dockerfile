FROM golang:1.22

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 443

# 运行应用程序
CMD ["./main"]
