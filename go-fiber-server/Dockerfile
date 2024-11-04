FROM debian:stable-slim

# 设置工作目录
WORKDIR /app

# 添加程序
ADD ./blog-server .

# 需要暴露的端口
EXPOSE 4000

# 执行启动命令
CMD ["/app/blog-server", "-conf", "/configs/config-prod.yaml"]