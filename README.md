# Kenketsu-Plus 🩸

[![CI](https://github.com/akaitigo/kenketsu-plus/actions/workflows/ci.yml/badge.svg)](https://github.com/akaitigo/kenketsu-plus/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

献血ルーム空き状況リアルタイム表示・献血記録管理・血液型別在庫逼迫通知アプリ

## 特徴

- **献血ルームマップ** — Leaflet.jsで献血ルーム・献血バスの位置と空き状況をリアルタイム表示
- **献血記録管理** — 献血履歴の記録、次回献血可能日の自動計算（全血/成分、性別別間隔制限対応）
- **在庫ダッシュボード** — 血液型別（A+/A-/B+/B-/O+/O-/AB+/AB-）の在庫逼迫度を色分け表示
- **プッシュ通知** — 在庫逼迫時にPWAプッシュ通知でアラート

## クイックスタート

```bash
# フロントエンド
cd frontend && npm install && npm run dev

# API
cd api && go run ./cmd/server

# 全チェック（lint + test + build）
make check
```

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| Frontend | Next.js 15 / TypeScript / Leaflet.js / PWA |
| Backend | Go 1.23 (net/http, 標準ライブラリ) |
| Database | PostgreSQL |
| Lint | biome + oxlint (TS) / golangci-lint (Go) |
| Test | vitest (TS) / go test (Go) |
| CI | GitHub Actions |

## アーキテクチャ

```
kenketsu-plus/
├── frontend/          # Next.js App Router
│   ├── src/app/       # ページ (/, /map, /donations, /inventory)
│   ├── src/components/# UIコンポーネント
│   ├── src/lib/       # APIクライアント, PWAヘルパー
│   └── src/types/     # 共通型定義
├── api/               # Go HTTP API
│   ├── cmd/server/    # エントリーポイント
│   ├── internal/
│   │   ├── handler/   # HTTPハンドラ + ルーター
│   │   ├── model/     # ドメインモデル
│   │   ├── repository/# インメモリリポジトリ
│   │   └── service/   # ビジネスロジック（献血間隔計算）
│   └── migrations/    # PostgreSQL DDL
└── docs/adr/          # Architecture Decision Records
```

## API エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | /health | ヘルスチェック |
| GET | /api/centers | 献血ルーム一覧 (?lat=&lng=&radius= 対応) |
| GET | /api/centers/{id} | 献血ルーム詳細 |
| POST | /api/centers | 献血ルーム登録 |
| GET | /api/donations | 献血記録一覧 |
| POST | /api/donations | 献血記録登録 |
| GET | /api/donations/next-available | 次回献血可能日計算 (?gender=) |
| GET | /api/inventory | 血液型別在庫一覧 |
| PUT | /api/inventory/{bloodType} | 在庫レベル更新 |
| POST | /api/subscriptions | プッシュ通知購読 |
| DELETE | /api/subscriptions/{id} | 購読解除 |
| POST | /api/notify/inventory-alert | 在庫逼迫通知トリガー |

## 献血間隔ルール

| 種別 | 男性 | 女性 |
|------|------|------|
| 全血400ml | 12週 | 16週 |
| 全血200ml | 4週 | 4週 |
| 成分献血 | 2週 | 2週 |
| 年間回数上限(全血) | 3回 | 2回 |

詳細: [ADR-001](docs/adr/001-donation-interval-calculation.md)

## デモ

<!-- デモGIF/スクリーンショットをここに追加 -->
*Coming soon*

## ライセンス

[MIT](LICENSE)
