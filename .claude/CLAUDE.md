# kenketsu-plus — Agent Instructions

## アーキテクチャ
モノレポ: frontend/ (Next.js) + api/ (Go) + PostgreSQL

## フロントエンド (frontend/)
- Next.js App Router, TypeScript strict
- biome.json + oxlint でlint/format
- Leaflet.js でマップ表示
- vitest でテスト

## バックエンド (api/)
- Go net/http, 標準ライブラリ優先
- golangci-lint (.golangci.yml)
- テスト: `go test -race ./...`

## DB
- PostgreSQL, マイグレーションは api/migrations/
- テーブル: donation_centers, donations, blood_inventory, push_subscriptions

## 献血間隔ルール
| 種別 | 男性 | 女性 |
|------|------|------|
| 全血400ml | 12週 | 16週 |
| 全血200ml | 4週 | 4週 |
| 成分献血 | 2週 | 2週 |
| 年間回数上限(全血) | 3回 | 2回 |

## CI
- `.github/workflows/ci.yml` — lint + test + build (frontend & api)
