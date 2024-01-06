# 设置基础镜像
FROM 1.22rc1-alpine3.19 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制应用程序源代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# 创建最终的镜像
#FROM alpine:latest
FROM golang:1.22rc1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .

# 设置容器启动命令
CMD ["./app"]
