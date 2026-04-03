# Kenketsu-Plus

献血ルーム空き状況リアルタイム表示・献血記録管理・血液型別在庫逼迫通知アプリ

## セットアップ

```bash
cd frontend && npm install
cd ../api && go mod download
```

## 開発

```bash
# フロントエンド
cd frontend && npm run dev

# API
cd api && go run ./cmd/server

# 全チェック
make check
```

## 技術スタック

- Frontend: Next.js / TypeScript / Leaflet.js / PWA
- Backend: Go (net/http)
- Database: PostgreSQL
