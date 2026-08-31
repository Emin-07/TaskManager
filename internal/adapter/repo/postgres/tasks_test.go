package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/testutil"
)

func TestTaskRepo_CRUD(t *testing.T) {
	db := testutil.OpenTestDB(t)
	repo := NewTaskRepo(db)
	ctx := context.Background()

	// grab a seeded user id (admin from the seed migrations)
	var userID int
	require.NoError(t, db.GetContext(ctx, &userID, "SELECT id FROM users LIMIT 1"))

	t.Run("insert and get by owner", func(t *testing.T) {
		require.NoError(t, repo.Insert(ctx, "Integration task", "some text", 3, 7, userID))

		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM tasks WHERE title = 'Integration task'"))

		task, err := repo.Get(ctx, id, userID, "user")
		require.NoError(t, err)
		assert.Equal(t, "Integration task", task.Title)
		assert.Equal(t, 3, task.Priority)
		assert.Equal(t, userID, task.UserId)
		assert.False(t, task.Expires.IsZero())
	})

	t.Run("get returns ErrNoRecord for other user", func(t *testing.T) {
		var otherID int
		require.NoError(t, db.GetContext(ctx, &otherID, "SELECT id FROM users WHERE id <> $1 LIMIT 1", userID))

		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM tasks WHERE title = 'Integration task'"))

		_, err := repo.Get(ctx, id, otherID, "user")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})

	t.Run("patch updates fields", func(t *testing.T) {
		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM tasks WHERE title = 'Integration task'"))

		require.NoError(t, repo.Patch(ctx, "Updated title", "", "user", 0, 0, id, userID))

		var title string
		require.NoError(t, db.GetContext(ctx, &title, "SELECT title FROM tasks WHERE id = $1", id))
		assert.Equal(t, "Updated title", title)
	})

	t.Run("list returns owned tasks", func(t *testing.T) {
		tasks, err := repo.List(ctx, 5, 0, userID)
		require.NoError(t, err)
		assert.NotEmpty(t, tasks)
	})

	t.Run("delete owned task", func(t *testing.T) {
		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM tasks WHERE title = 'Updated title'"))

		require.NoError(t, repo.Delete(ctx, id, userID, "user"))

		_, err := repo.Get(ctx, id, userID, "user")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})
}
