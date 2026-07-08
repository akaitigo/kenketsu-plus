package repository_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

// migrationFiles are applied, in order, to bring a clean schema up to date.
var migrationFiles = []string{
	"001_init.sql",
	"002_push_subscriptions_unique_endpoint.up.sql",
}

// setupTestDB connects to the PostgreSQL instance named by DATABASE_URL and
// applies the migrations against a freshly reset schema. The integration tests
// are skipped when DATABASE_URL is unset (e.g. a plain local `go test`) or when
// the database is unreachable, mirroring the CI postgres-service pattern (#22).
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping PostgreSQL integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("failed to open database (not available): %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Skipf("failed to ping database (not available): %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	resetSchema(t, db)
	return db
}

// resetSchema drops the application tables and re-applies all migrations so each
// test starts from a known, clean state. This is destructive and is intended to
// run only against a disposable test database.
func resetSchema(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx := context.Background()
	if _, err := db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS push_subscriptions, blood_inventory, donations, donation_centers CASCADE`,
	); err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}

	for _, name := range migrationFiles {
		path := filepath.Join("..", "..", "migrations", name)
		content, err := os.ReadFile(path) //nolint:gosec // fixed migration paths, not user input
		if err != nil {
			t.Fatalf("failed to read migration %s: %v", name, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			t.Fatalf("failed to apply migration %s: %v", name, err)
		}
	}
}

func TestPgSubscriptionRepository_UpsertByEndpoint(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPgSubscriptionRepository(db)
	ctx := context.Background()

	first, err := repo.Create(ctx, &model.PushSubscription{
		Endpoint: "https://push.example.com/sub-a",
		P256dh:   "key-1",
		Auth:     "auth-1",
	})
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	second, err := repo.Create(ctx, &model.PushSubscription{
		Endpoint: "https://push.example.com/sub-a",
		P256dh:   "key-2",
		Auth:     "auth-2",
	})
	if err != nil {
		t.Fatalf("second create (upsert) failed: %v", err)
	}

	subs := repo.List(ctx)
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription after re-subscribe, got %d", len(subs))
	}
	if second.ID != first.ID {
		t.Errorf("expected stable ID %q across upsert, got %q", first.ID, second.ID)
	}
	if subs[0].P256dh != "key-2" || subs[0].Auth != "auth-2" {
		t.Errorf("expected keys updated to key-2/auth-2, got %s/%s", subs[0].P256dh, subs[0].Auth)
	}
}

func TestPgSubscriptions_UniqueConstraintEnforced(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	insert := func() error {
		_, err := db.ExecContext(
			ctx,
			`INSERT INTO push_subscriptions (endpoint, p256dh, auth) VALUES ($1, $2, $3)`,
			"https://push.example.com/dup", "k", "a",
		)
		return err
	}

	if err := insert(); err != nil {
		t.Fatalf("first raw insert failed: %v", err)
	}
	if err := insert(); err == nil {
		t.Fatal("expected UNIQUE(endpoint) violation on duplicate insert, got nil")
	}
}

func TestPgInventoryRepository_UpdateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPgInventoryRepository(db)
	ctx := context.Background()

	updated, err := repo.Update(ctx, model.BloodTypeAPos, model.InventoryLevelCritical)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Level != model.InventoryLevelCritical {
		t.Errorf("expected level critical, got %s", updated.Level)
	}

	got, err := repo.GetByBloodType(ctx, model.BloodTypeAPos)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Level != model.InventoryLevelCritical {
		t.Errorf("expected persisted level critical, got %s", got.Level)
	}

	if _, err := repo.Update(ctx, model.BloodType("Z+"), model.InventoryLevelNormal); err == nil {
		t.Error("expected error for invalid blood type, got nil")
	}
}

// TestPgInventoryRepository_NotFound exercises the errors.Is(sql.ErrNoRows) path
// (#19): updating a valid blood type whose row is missing must return a clear
// not-found error rather than a raw driver error.
func TestPgInventoryRepository_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPgInventoryRepository(db)
	ctx := context.Background()

	if _, err := db.ExecContext(ctx, `DELETE FROM blood_inventory WHERE blood_type = $1`, string(model.BloodTypeAPos)); err != nil {
		t.Fatalf("failed to delete inventory row: %v", err)
	}

	_, err := repo.Update(ctx, model.BloodTypeAPos, model.InventoryLevelCritical)
	if err == nil {
		t.Fatal("expected not-found error after deleting inventory row, got nil")
	}
}

func TestPgCenterRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPgCenterRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, &model.DonationCenter{
		Name:           "Shibuya Center",
		Address:        "Shibuya, Tokyo",
		Lat:            35.6580,
		Lng:            139.7016,
		Capacity:       50,
		AvailableSlots: 10,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated ID")
	}
	if created.Status != model.CenterStatusOpen {
		t.Errorf("expected default status open, got %s", created.Status)
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if fetched.Name != "Shibuya Center" {
		t.Errorf("expected Shibuya Center, got %s", fetched.Name)
	}

	if all := repo.List(ctx); len(all) != 1 {
		t.Errorf("expected 1 center, got %d", len(all))
	}

	near := repo.ListByDistance(ctx, 35.66, 139.70, 5)
	if len(near) != 1 {
		t.Errorf("expected 1 center within 5km, got %d", len(near))
	}

	far := repo.ListByDistance(ctx, 34.6937, 135.5023, 5)
	if len(far) != 0 {
		t.Errorf("expected 0 centers within 5km of Osaka, got %d", len(far))
	}
}

func TestPgDonationRepository_CreateAndList(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPgDonationRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, &model.Donation{
		BloodType:    model.BloodTypeOPos,
		DonationType: model.DonationTypeWhole400,
		Gender:       model.GenderMale,
		DonatedAt:    time.Now().UTC().Truncate(time.Second),
		VolumeMl:     400,
		Memo:         "integration test",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated ID")
	}

	list := repo.List(ctx)
	if len(list) != 1 {
		t.Fatalf("expected 1 donation, got %d", len(list))
	}
	if list[0].BloodType != model.BloodTypeOPos || list[0].Memo != "integration test" {
		t.Errorf("unexpected donation fields: %+v", list[0])
	}
}
