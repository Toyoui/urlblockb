# 基础镜像
FROM golang:latest

# 使用 Go Modules
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

# 设置工作目录
WORKDIR /app

# 初始化Go modules
RUN go mod init example.com/urlodai

# 将代码复制到容器中
COPY . .

# 编译Go程序
RUN go build -o main .

# 运行Go程序
CMD ["./main"]
