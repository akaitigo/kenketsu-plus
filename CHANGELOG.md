# Changelog

## [1.2.0] - 2026-04-04

### Security
- `NOTIFY_SECRET` 必須化 — 未設定時は503を返す (C-2 強化)
- `RequireAdminKey` ミドルウェア追加: POST /api/centers, PUT /api/inventory に適用 (C-3)

### Fixed
- 通知レスポンスが実送信していないのに「送信しました」と返す問題 → 「キューしました」に修正 (NEW-1)
- リポジトリエラーが全て400で返る問題 → ValidationError→400、その他→500に分類 (M-7)

### Improved
- `NotificationToggle` をインベントリページに接続 (H-4)
- `DonationForm` のレンダリングテスト追加 (H-8)
- `api.ts` にオプションランタイムバリデーターを追加 (M-2)

## [1.1.0] - 2026-04-04

### Security
- CORS origin を環境変数 `CORS_ALLOWED_ORIGIN` で制限 (C-1)
- 通知APIに `X-Notify-Secret` ヘッダー認証追加 (C-2)
- セキュリティヘッダー追加: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy` (H-2)
- 全POST/PUTハンドラに `MaxBytesReader` (1MB) 適用 (H-3)
- HTTP Server に `ReadTimeout`/`WriteTimeout`/`IdleTimeout` 設定 (H-1)

### Fixed
- `notify.go` Targets フィールドの二重JSONエンコードバグ (H-7)
- 成分献血で `volumeMl=0` がバリデーションエラーになる問題 (M-4)
- `NextDonationDate` の useEffect クリーンアップ未実装による競合状態 (H-9)
- `sw.js` の push data パースでmalformed JSONによるSWクラッシュ (M-1)
- `radius` パラメータに0-500km範囲バリデーション追加 (H-5)

### Performance
- `ListByDistance` のHaversine距離計算をキャッシュしソート時の再計算を排除 (M-3)

## [1.0.0] - 2026-04-04

### Added
- **献血ルームマップ** (MVP-2): Leaflet.jsマップ、マーカー表示、距離フィルタ
- **献血記録管理** (MVP-3): CRUD API、次回献血可能日自動計算、献血間隔制限対応
- **在庫ダッシュボード** (MVP-4): 血液型別バーチャート、逼迫度色分け、アラートバナー
- **PWAプッシュ通知** (MVP-5): Service Worker、購読/解除、在庫逼迫通知トリガー
- **プロジェクト基盤** (MVP-1): Go API + Next.js モノレポ、CI、ドメインモデル

### Architecture Decisions
- ADR-001: 献血間隔制限の計算アルゴリズム（スライディングウィンドウ方式）
