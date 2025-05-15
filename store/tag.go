package store

import (
	"context"
	"errors"

	"github.com/go-webapi-team/task-app/entity"
)

var ErrNotFound = errors.New("not found")

func (r *Repository) CreateTag(ctx context.Context, db Execer, t *entity.Tag) error {
	const sqlStr = `INSERT INTO tags (user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`

	now := r.Clocker.Now()
	res, err := db.ExecContext(ctx, sqlStr, t.UserID, t.Name, now, now)

	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TagID(id)
	return nil
}

func (r *Repository) ListTags(ctx context.Context, db Queryer, userID int64) (entity.Tags, error) {
	const query = `SELECT id, user_id, name, created_at, updated_at FROM tags WHERE user_id = ?`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags entity.Tags
	for rows.Next() {
		var t entity.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) DeleteTag(ctx context.Context, db Execer, userID int64, id entity.TagID) error {
	res, err := db.ExecContext(ctx, `DELETE FROM tags WHERE id=? AND user_id=?`, id, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) AddTagToTask(ctx context.Context, db Execer, userID int64, t *entity.TaskTag) error {
	now := r.Clocker.Now()
	const q = `
      INSERT INTO tasks_tags (task_id, tag_id, created_at, updated_at)
      SELECT t.id, tg.id, ?, ?
        FROM tasks AS t
        JOIN tags  AS tg ON tg.id = ?
       WHERE t.id = ? AND t.user_id = ? AND tg.user_id = ?
    `
	res, err := db.ExecContext(ctx, q,
		now, now,
		t.TagID,
		t.TaskID, userID, userID,
	)
	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskTagID(id)
	return nil
}

func (r *Repository) DeleteTagFromTask(ctx context.Context, db Execer, userID int64, t *entity.TaskTag) error {
	const q = `
      DELETE tt
        FROM tasks_tags AS tt
        JOIN tasks AS t  ON tt.task_id = t.id
        JOIN tags  AS tg ON tt.tag_id  = tg.id
       WHERE tt.task_id = ? AND tt.tag_id = ?
         AND t.user_id  = ? AND tg.user_id = ?
    `
	res, err := db.ExecContext(ctx, q, t.TaskID, t.TagID, userID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
