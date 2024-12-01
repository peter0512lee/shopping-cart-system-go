FROM golang:1.21-alpine

# 設置工作目錄
WORKDIR /app

# 安裝基本工具
RUN apk add --no-cache git

# 複製 go.mod 和 go.sum
COPY go.mod ./
COPY go.sum ./

# 下載相依套件
RUN go mod download

# 複製整個專案
COPY . .

# 移動到 cmd/api 目錄編譯
WORKDIR /app/cmd/api

# 編譯應用
RUN go build -o main .

# 移回工作目錄
WORKDIR /app

# 移動編譯後的檔案到工作目錄
RUN mv /app/cmd/api/main .

EXPOSE 8080

# 使用編譯後的執行檔
CMD ["./main"]