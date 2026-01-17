FROM golang:1.24-alpine

WORKDIR /app

# まず go.mod / go.sum をコピー
COPY ./test_app/go.mod ./test_app/go.sum ./
RUN go mod download

# 全コードをコピー
COPY ./test_app .

# ビルド（※ main.go は app直下にある）
RUN go build -o main ./main.go

EXPOSE 8000
CMD ["./main"]
