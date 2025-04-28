package store

import (
	"context"
	"errors"
	"strings"

	"github.com/go-webapi-team/task-app/entity"
)

var ErrTaskNotFound = errors.New("Task not found")

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	const q = `INSERT INTO tasks
		(user_id,title,description,deadline,priority,is_done,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?,?)`
	now := r.Clocker.Now()
	res, err := db.ExecContext(ctx, q,
		t.UserID, t.Title, t.Description, t.Deadline, t.Priority, t.IsDone, now, now)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID, t.CreatedAt, t.UpdatedAt = entity.TaskID(id), now, now
	return nil
}

// GetTask returns a single task by id (user 自身の物のみ取得する想定)
func (r *Repository) GetTask(ctx context.Context, db Queryer, userID int64, id entity.TaskID) (*entity.Task, error) {
	const q = `SELECT id,user_id,title,description,deadline,priority,is_done,created_at,updated_at
			  FROM tasks WHERE id=? AND user_id=?`
	var t entity.Task
	if err := db.QueryRowContext(ctx, q, id, userID).
		Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Deadline,
			&t.Priority, &t.IsDone, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

// dynamic filter list
type ListTaskFilter struct {
	IsDone *bool
	TagID  *int64
	DueAsc *bool // nil = no sort
}

// ListTasks with optional filters.
func (r *Repository) ListTasks(ctx context.Context, db Queryer, userID int64, f ListTaskFilter) (entity.Tasks, error) {
	sb := strings.Builder{}
	args := make([]any, 0, 4)

	sb.WriteString(`SELECT DISTINCT t.id,t.user_id,t.title,t.description,t.deadline,t.priority,
	                t.is_done,t.created_at,t.updated_at
	                FROM tasks t`)
	if f.TagID != nil {
		sb.WriteString(` INNER JOIN tasks_tags tt ON tt.task_id=t.id`)
	}
	sb.WriteString(` WHERE t.user_id=?`)
	args = append(args, userID)

	if f.IsDone != nil {
		sb.WriteString(` AND t.is_done=?`)
		args = append(args, *f.IsDone)
	}
	if f.TagID != nil {
		sb.WriteString(` AND tt.tag_id=?`)
		args = append(args, *f.TagID)
	}
	if f.DueAsc != nil {
		if *f.DueAsc {
			sb.WriteString(` ORDER BY t.deadline ASC`)
		} else {
			sb.WriteString(` ORDER BY t.deadline DESC`)
		}
	}

	rows, err := db.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts entity.Tasks
	for rows.Next() {
		var t entity.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Deadline,
			&t.Priority, &t.IsDone, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		ts = append(ts, &t)
	}
	return ts, rows.Err()
}

// UpdateTask (全カラム更新; PATCH を細かく分けたい場合は別関数を切る)
func (r *Repository) UpdateTask(ctx context.Context, db Execer, t *entity.Task) error {
	const q = `UPDATE tasks SET title=?,description=?,deadline=?,priority=?,is_done=?,updated_at=?
			   WHERE id=? AND user_id=?`
	_, err := db.ExecContext(ctx, q, t.Title, t.Description, t.Deadline, t.Priority,
		t.IsDone, r.Clocker.Now(), t.ID, t.UserID)
	return err
}

// DeleteTask hard‑delete
func (r *Repository) DeleteTask(ctx context.Context, db Execer, userID int64, id entity.TaskID) error {
	_, err := db.ExecContext(ctx, `DELETE FROM tasks WHERE id=? AND user_id=?`, id, userID)
	return err
}

// ToggleTaskDone flips is_done.
func (r *Repository) ToggleTaskDone(ctx context.Context, db Execer, userID int64, id entity.TaskID) error {
	const q = `UPDATE tasks SET is_done = NOT is_done, updated_at=? WHERE id=? AND user_id=?`
	res, err := db.ExecContext(ctx, q, r.Clocker.Now(), id, userID)
	if err != nil {
        return err
    }

	n, err := res.RowsAffected()
	if err != nil {
        return err
    }
    
    if n == 0 {
        // 更新対象の行が無かったら「タスクが見つからない」とみなす
        return ErrTaskNotFound
    }
	return nil
}