# 使用与开发时版本匹配的基础镜像
FROM golang:1.17 AS Bulder

# 设置工作路径
WORKDIR /go/src/app

# 将go项目copy到 /go/sec/app目录下
COPY . .
#原始方式：直接镜像内打包编译
RUN export GO111MODULE=auto && export GOPROXY=https://goproxy.cn && go mod tidy

# 编译go源代码 并命名为go_gateway 存放到 /bin/go_gateway
RUN go build -o ./bin/go_gateway

# 运行命令 通过命令行参数的方式指定配置文件和启动类型
CMD ./bin/go_gateway -config=./conf/dev/ -endpoint=dashboard