# Harvest: kenketsu-plus

## メトリクス

| 項目 | v1.0.0 | v1.1.0 | v1.2.0 |
|------|--------|--------|--------|
| コミット数 | 6 | 8 | 10 |
| PR数（マージ済み） | 5 | 7 | 9 |
| Issue数（クローズ済み） | 5 | 5 | 5 |
| ADR数 | 1 | 1 | 1 |
| Go テストケース | 58 | 60 | 60 |
| Frontend テストケース | 28 | 28 | 29 |
| 合計テスト | 86 | 88 | 89 |
| レビュー指摘対応 | — | 12件 | +7件 (計19件) |

## レビュー修正履歴

### v1.1.0 (初回レビュー対応: 12件)
- C-1: CORS環境変数制限
- C-2: 通知API認証 (X-Notify-Secret)
- H-1: ReadTimeout/WriteTimeout/IdleTimeout
- H-2: セキュリティヘッダー
- H-3: MaxBytesReader (1MB)
- H-5: radius バリデーション
- H-7: Targets二重エンコードバグ
- H-9: useEffect クリーンアップ
- M-1: sw.js try/catch
- M-3: Haversineキャッシュ
- M-4: 成分献血volume修正

### v1.2.0 (再レビュー対応: 7件)
- C-2強化: NOTIFY_SECRET 必須化（未設定→503）
- C-3: RequireAdminKey ミドルウェア（POST centers, PUT inventory）
- H-4: NotificationToggle をインベントリページに接続
- H-8: DonationForm レンダリングテスト（.ts→.tsx、render+fireEvent）
- M-2: api.ts にオプションランタイムバリデーター
- M-7: writeRepoError でValidationError→400、その他→500に分類
- NEW-1: 通知レスポンスを「キューしました」に修正

### 残存（次フェーズ候補）
| ID | 内容 | 工数 |
|----|------|------|
| H-6 | Repository interface + PostgreSQL接続 | 1-2日 |
| M-5 | ユーザー分離（認証基盤が前提） | 2-3日 |

## テンプレート改善提案
1. `.golangci.yml` に `fieldalignment` チェック推奨設定
2. vitest.config.ts に `esbuild: { jsx: "automatic" }` + `jsdom`
3. Go `go.sum` を初期スキャフォールドに含める
4. CORSミドルウェアに環境変数制限をデフォルトで含める
5. セキュリティヘッダーをLayer-0テンプレートに含める
6. `MaxBytesReader` + `RequireAdminKey` をPOST/PUTテンプレートに含める
7. `writeRepoError` パターン（ValidationError分類）をハンドラテンプレートに含める
