# 使用Go官方镜像作为基础镜像
FROM golang:1.22.0

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" &&\
    chmod +x /go/kubectl && mv /go/kubectl  /usr/bin/kubectl && \
    mkdir /root/.kube
COPY config /root/.kube/config
COPY kubectl /usr/bin/kubectl
COPY kubebuilder /usr/bin/kubebuilder
# 设置工作目录
WORKDIR /workspace