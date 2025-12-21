FROM golang:1.24-alpine

WORKDIR /app

# まず go.mod / go.sum をコピー
COPY ./app/go.mod ./app/go.sum ./
RUN go mod download

# 全コードをコピー
COPY ./app .

# ビルド（※ main.go は app直下にある）
RUN go build -o main ./entry_point/cmd/main.go

CMD ["tail", "-f", "/dev/null"]
