package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (m TaskRepo) List(ctx context.Context, limit, offset, userId int) ([]*repo.TaskDb, error) {
	var tasks []*repo.TaskDb
	query := "SELECT t.id, t.title, t.text, t.priority, t.created, t.expires, t.user_id FROM tasks t WHERE (t.expires > CURRENT_TIMESTAMP AND t.user_id = $1) OR EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = 'admin') ORDER BY t.id LIMIT $2 OFFSET $3;"
	err := m.DB.SelectContext(ctx, &tasks, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m TaskRepo) Get(ctx context.Context, id, userId int, role string) (*repo.TaskDb, error) {
	task := &repo.TaskDb{}
	err := m.DB.GetContext(ctx, task, "SELECT id, title, text, priority, created, expires, user_id FROM tasks WHERE id = $1 AND (user_id = $2 OR $3 = 'admin')", id, userId, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return task, nil
}

func (m TaskRepo) Insert(ctx context.Context, title, text string, priority, expireDays, userId int) error {
	query := `INSERT INTO tasks (title, text, priority, expires, user_id) VALUES ($1, $2, $3, CURRENT_TIMESTAMP + MAKE_INTERVAL(days => $4), $5)`
	_, err := m.DB.ExecContext(ctx, query, title, text, priority, expireDays, userId)

	if err != nil {
		return err
	}

	return nil
}

func queryOrderTracker(query *strings.Builder, isNotFirst *bool) {
	if *isNotFirst {
		query.WriteString(`, `)
	} else {
		*isNotFirst = true
	}
}

func (m TaskRepo) Patch(ctx context.Context, title, text, userRole string, priority, expireDays, id, userId int) error {
	var query strings.Builder
	var args []any
	var isNotFirst bool
	cnt := 1
	query.WriteString("UPDATE tasks SET ")
	if title != "" {
		queryOrderTracker(&query, &isNotFirst)
		fmt.Fprintf(&query, `title = $%d`, cnt)
		cnt++
		args = append(args, title)
	}
	if text != "" {
		queryOrderTracker(&query, &isNotFirst)
		fmt.Fprintf(&query, `text = $%d `, cnt)
		cnt++
		args = append(args, text)
	}
	if priority != 0 {
		queryOrderTracker(&query, &isNotFirst)
		fmt.Fprintf(&query, `priority = $%d `, cnt)
		cnt++
		args = append(args, priority)
	}
	if expireDays != 0 {
		queryOrderTracker(&query, &isNotFirst)
		fmt.Fprintf(&query, `expires = CURRENT_TIMESTAMP + MAKE_INTERVAL(days => $%d) `, cnt)
		cnt++
		args = append(args, expireDays)
	}
	if len(args) == 0 {
		return domain.ErrNoData
	}
	args = append(args, id, userId, userRole)
	fmt.Fprintf(&query, "WHERE id = $%d AND (user_id = $%d OR $%d = 'admin')", cnt, cnt+1, cnt+2)

	result, err := m.DB.ExecContext(ctx, query.String(), args...)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNoRecord
	}

	return nil
}

func (m TaskRepo) Delete(ctx context.Context, id, userId int, role string) error {
	result, err := m.DB.ExecContext(ctx, "DELETE FROM tasks t WHERE t.id = $1 AND (t.user_id = $2 OR $3 = 'admin')", id, userId, role)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNoRecord
	}
	return nil
}
