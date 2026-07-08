-- 002_push_subscriptions_unique_endpoint.up.sql のロールバック（#18）
-- UNIQUE 制約を削除する。クリーンアップで削除した重複データは復元されない。
ALTER TABLE push_subscriptions
    DROP CONSTRAINT IF EXISTS uq_push_subscriptions_endpoint;
