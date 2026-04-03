-- 献血ルーム
CREATE TABLE donation_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address VARCHAR(500) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 0,
    available_slots INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 献血記録
CREATE TABLE donations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_type VARCHAR(10) NOT NULL,
    donation_type VARCHAR(20) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    donated_at TIMESTAMPTZ NOT NULL,
    volume_ml INTEGER NOT NULL,
    memo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 血液型別在庫
CREATE TABLE blood_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blood_type VARCHAR(10) NOT NULL UNIQUE,
    level VARCHAR(20) NOT NULL DEFAULT 'normal',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- プッシュ通知購読
CREATE TABLE push_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint TEXT NOT NULL,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 初期在庫データ
INSERT INTO blood_inventory (blood_type, level) VALUES
    ('A+', 'normal'), ('A-', 'normal'),
    ('B+', 'normal'), ('B-', 'normal'),
    ('O+', 'normal'), ('O-', 'normal'),
    ('AB+', 'normal'), ('AB-', 'normal');

CREATE INDEX idx_donation_centers_lat_lng ON donation_centers (lat, lng);
CREATE INDEX idx_donations_donated_at ON donations (donated_at);
CREATE INDEX idx_blood_inventory_blood_type ON blood_inventory (blood_type);
