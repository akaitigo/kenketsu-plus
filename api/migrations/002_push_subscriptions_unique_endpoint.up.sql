-- push_subscriptions.endpoint に UNIQUE 制約を追加（#18）
-- 同一エンドポイントの重複登録による重複通知を防止する。

-- 1. 既存の重複データをクリーンアップ。
--    同一 endpoint については最新（created_at が最大）の1件のみ残す。
--    created_at が同値の場合は id が最小の1件を残す。
DELETE FROM push_subscriptions a
USING push_subscriptions b
WHERE a.endpoint = b.endpoint
  AND (
        a.created_at < b.created_at
     OR (a.created_at = b.created_at AND a.id > b.id)
  );

-- 2. UNIQUE 制約を追加。
ALTER TABLE push_subscriptions
    ADD CONSTRAINT uq_push_subscriptions_endpoint UNIQUE (endpoint);
