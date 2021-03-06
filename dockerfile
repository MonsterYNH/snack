
# 得到最新的 golang docker 镜像
FROM golang:latest
# 在容器内部创建一个目录来存储我们的 web 应用
RUN mkdir -p /go/src/snack
#接着使它成为工作目录。
WORKDIR /go/src/snack/
# 复制 go_web 目录到容器中
COPY . /go/src/snack/
#编译，编译成可执行文件
RUN go build /go/src/snack/main.go
# 给主机暴露 80 端口，这样外部网络可以访问你的应用
EXPOSE 5000
# 告诉 Docker 启动容器运行的命令
CMD /go/src/snack/main