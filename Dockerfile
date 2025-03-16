# TODO: 各バージョンをgo.modとGitHubActionと合わせる

# Build
FROM golang:1.24.1-bullseye AS builder
WORKDIR /app
COPY . .
# go buildでmain（実行ファイル）を作成
RUN go build -o main ./cmd && \
    # マイグレーションのためにgolang-migrate（実行ファイル）をダウンロード
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run
FROM debian:bullseye-20250224-slim
WORKDIR /app
# builder stageで作成されたmain（実行ファイル）をコピー
COPY --from=builder /app/main .
# builder stageで作成されたmigrate（実行ファイル）をコピー
COPY --from=builder /app/migrate ./migrate
# 環境変数が定義されたファイルをコピー
# TODO: Load env from secret
COPY env/app.env /env/app.env
# スキーマなどが定義されているsqlファイルをコピー
COPY db/migration ./migration
# 起動スクリプトをコピー
COPY start.sh .

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]