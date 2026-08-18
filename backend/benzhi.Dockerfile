# 评测用镜像：构建并直接启动项目真实入口。
FROM golang:1.22
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/service ./cmd/server
CMD ["/usr/local/bin/service"]

# 多架构交叉构建示例（如需交付双架构镜像）：
# docker buildx build --platform linux/arm64,linux/amd64 -f benzhi.Dockerfile -t <image> .
