package store

import (
	"context"
	"github.com/go-webapi-team/task-app/entity"
)

func (r *Repository) CreateTag(ctx context.Context, db Execer, t *entity.Tag) error {
	const sqlStr = `INSERT INTO tags (user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	
	now := r.Clocker.Now()
	result, err := db.ExecContext(ctx, sqlStr, t.UserID, t.Name, now, now)
	
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TagID(id)
	return nil
}

func (r *Repository) ListTags(ctx context.Context, db Queryer) (entity.Tags, error) {
	const sqlStr = `SELECT id, user_id, name, created_at, updated_at FROM tags`

	rows, err := db.QueryContext(ctx, sqlStr)
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