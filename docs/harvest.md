# Harvest: kenketsu-plus

## メトリクス

| 項目 | 値 |
|------|-----|
| コミット数 | 6 |
| PR数 | 5 (全てマージ済み) |
| Issue数 | 5 (全てクローズ済み) |
| ADR数 | 1 |
| CLAUDE.md行数 | 34 |
| CI | GitHub Actions (Frontend + API) |
| Go テストケース | 58 passed (4パッケージ) |
| Frontend テストケース | 28 passed (6ファイル) |
| 合計テスト | 86 |
| v1.0.0 タグ | 作成済み |

## プロジェクト概要

献血ルーム空き状況・献血記録管理・血液型別在庫逼迫通知のWebアプリ。
Go API + Next.js フロントエンドのモノレポ構成。

## MVP機能完了状況

1. **MVP-1 基盤構築** — Go API、Next.js、DB スキーマ、CI ✅
2. **MVP-2 マップ表示** — Leaflet.js マップ、献血ルーム CRUD ✅
3. **MVP-3 献血記録** — 献血間隔計算、ADR-001 ✅
4. **MVP-4 在庫ダッシュボード** — バーチャート、逼迫アラート ✅
5. **MVP-5 PWA通知** — Service Worker、購読/解除 ✅

## 振り返り

### うまくいったこと
- Go標準ライブラリのnet/httpルーターが十分で、外部依存ゼロ
- biome + oxlint + golangci-lintの厳格なlint設定でコード品質を維持
- 献血間隔計算のテストで性別×種別全組み合わせ + エッジケースをカバー
- CLAUDE.mdが34行に収まり50行制限をクリア

### 改善点
- CIのfieldalignment問題で3回修正pushが必要だった → ローカルにfieldalignmentツールを入れて事前チェック
- フロントエンドのJSXテスト環境（jsdom + esbuild jsx:automatic）のセットアップに手間取った
- MVP-2以降の並列エージェント実行がAPI 500エラーで失敗 → フォールバックとして自分で順次実装

### テンプレート改善提案
1. `.golangci.yml` に `fieldalignment` チェックの推奨設定を追記
2. vitest.config.tsテンプレートに `esbuild: { jsx: "automatic" }` と `jsdom` を初期設定に含める
3. Go モジュールの `go.sum` 生成を初期スキャフォールドに含める（CI cache問題回避）
