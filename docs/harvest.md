# Harvest: kenketsu-plus

## メトリクス (v1.1.0)

| 項目 | v1.0.0 | v1.1.0 |
|------|--------|--------|
| コミット数 | 6 | 8 |
| PR数（マージ済み） | 5 | 7 |
| Issue数（クローズ済み） | 5 | 5 |
| ADR数 | 1 | 1 |
| CLAUDE.md行数 | 34 | 34 |
| Go テストケース | 58 | 60 (+2) |
| Frontend テストケース | 28 | 28 |
| 合計テスト | 86 | 88 |
| タグ | v1.0.0 | v1.1.0 |

## v1.1.0 コードレビュー修正サマリー

### 修正した指摘 (12件)

| 重要度 | 件数 | 内容 |
|--------|------|------|
| Critical | 2 | CORS制限 (C-1), 通知API認証 (C-2) |
| High | 5 | タイムアウト (H-1), セキュリティヘッダー (H-2), MaxBytesReader (H-3), radius検証 (H-5), Targetsバグ (H-7), useEffectクリーンアップ (H-9) |
| Medium | 4 | sw.js try/catch (M-1), Haversineキャッシュ (M-3), 成分献血volume (M-4) |

### 未対応（次フェーズ候補）

| ID | 内容 | 理由 |
|----|------|------|
| H-6 | Repository interface + DB接続 | 工数1-2日。ポートフォリオ評価への最大インパクト |
| H-4 | NotificationToggle のページ組み込み | VAPID鍵生成の環境準備が必要 |
| H-8 | DonationForm の動作テスト強化 | jsdom + render テスト拡充 |
| M-2 | api.ts の型アサーション改善 | zodなどランタイムバリデーション導入 |
| M-5 | ユーザー分離（認証） | 認証基盤の設計が前提 |
| C-3 | 管理者API認可 | 認証基盤の設計が前提 |

## 振り返り

### v1.0.0 → v1.1.0 の改善点
- セキュリティスコアが大幅改善（CORS制限、認証、ヘッダー、DoS対策）
- 成分献血の登録バグが解消（volumeMl=0許可）
- フロントエンドの競合状態が解消（useEffect cleanup）
- Haversine距離計算のパフォーマンス改善

### テンプレート改善提案（v1.0.0分含む）
1. `.golangci.yml` に `fieldalignment` チェック推奨設定を追記
2. vitest.config.ts に `esbuild: { jsx: "automatic" }` + `jsdom` を初期設定
3. Go `go.sum` 生成を初期スキャフォールドに含める（CI cache問題回避）
4. CORSミドルウェアのテンプレートに環境変数制限をデフォルトで含める
5. セキュリティヘッダーをLayer-0ミドルウェアテンプレートに含める
6. `MaxBytesReader` をPOST/PUTハンドラのテンプレートに含める
