FROM golang:1.16

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 443

# 运行应用程序
CMD ["./main"]