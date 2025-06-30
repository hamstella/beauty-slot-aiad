# マルチステージビルド
FROM golang:1.21-alpine AS builder

WORKDIR /app

# モジュールファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 実行用イメージ
FROM alpine:latest

# セキュリティ更新とSSL証明書をインストール
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# タイムゾーン設定
ENV TZ=Asia/Tokyo

# ポート公開
EXPOSE 8080

# アプリケーション実行
CMD ["./main"]