package store

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-webapi-team/task-app/clock"
	"github.com/go-webapi-team/task-app/entity"
	"github.com/google/go-cmp/cmp"
)

func TestTagRepository_ListTags(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	c := clock.FixedClocker{}

	// --- モック DB & 期待行を用意（sqlmock は database/sql 互換） ---
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	now := c.Now()
	want := entity.Tags{
		{ID: 1, UserID: 1, Name: "tag1", CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 2, Name: "tag2", CreatedAt: now, UpdatedAt: now},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
		AddRow(want[0].ID, want[0].UserID, want[0].Name, want[0].CreatedAt, want[0].UpdatedAt).
		AddRow(want[1].ID, want[1].UserID, want[1].Name, want[1].CreatedAt, want[1].UpdatedAt)

	mock.ExpectQuery(`SELECT id, user_id, name, created_at, updated_at FROM tags`).
		WillReturnRows(rows)

	sut := &Repository{Clocker: c}
	got, err := sut.ListTags(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("mismatch (-got +want):\n%s", diff)
	}
}

func TestTagRepository_CreateTag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	c := clock.FixedClocker{}
	now := c.Now()

	newTag := &entity.Tag{UserID: 99, Name: "new", CreatedAt: now, UpdatedAt: now}
	wantID := int64(123)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(`INSERT INTO tags`).
		WithArgs(newTag.UserID, newTag.Name, driver.Value(now), driver.Value(now)).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	r := &Repository{Clocker: c}
	if err := r.CreateTag(ctx, db, newTag); err != nil {
		t.Fatalf("error: %v", err)
	}
	if newTag.ID != entity.TagID(wantID) {
		t.Errorf("ID not set, got %d want %d", newTag.ID, wantID)
	}
}

func TestTagRepository_DeleteTag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	// 期待: id=5, user_id=1 の行を 1 行削除
	mock.ExpectExec(`DELETE FROM tags`).
		WithArgs(int64(5), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := &Repository{}
	if err := r.DeleteTag(ctx, db, 1, entity.TagID(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations not met: %v", err)
	}
}

