package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/testutil"
)

func TestUserRepo_CRUD(t *testing.T) {
	db := testutil.OpenTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	const uniqueEmail = "integration_unique@example.com"

	t.Run("insert and authenticate", func(t *testing.T) {
		require.NoError(t, repo.Insert(ctx, "integration_user", "user", uniqueEmail, []byte("hash")))

		user, err := repo.Authenticate(ctx, uniqueEmail)
		require.NoError(t, err)
		assert.Equal(t, "integration_user", user.Username)
		assert.Equal(t, "user", user.Role)
	})

	t.Run("get by id", func(t *testing.T) {
		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM users WHERE email = $1", uniqueEmail))

		user, err := repo.GetById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, uniqueEmail, user.Email)
	})

	t.Run("authenticate unknown email returns ErrNoRecord", func(t *testing.T) {
		_, err := repo.Authenticate(ctx, "nobody@example.com")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})

	t.Run("patch updates fields", func(t *testing.T) {
		var id int
		require.NoError(t, db.GetContext(ctx, &id, "SELECT id FROM users WHERE email = $1", uniqueEmail))

		require.NoError(t, repo.Patch(ctx, "renamed_user", "admin", "", nil, id))

		user, err := repo.GetById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "renamed_user", user.Username)
		assert.Equal(t, "admin", user.Role)
	})

	t.Run("list returns seeded + created users", func(t *testing.T) {
		users, err := repo.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, users)
	})

	t.Run("delete user and cascade tasks atomically", func(t *testing.T) {
		// create a user with a task, then delete via repo.Delete (transactional)
		require.NoError(t, repo.Insert(ctx, "doomed", "user", "doomed@example.com", []byte("hash")))
		var userID int
		require.NoError(t, db.GetContext(ctx, &userID, "SELECT id FROM users WHERE email = 'doomed@example.com'"))
		taskRepo := NewTaskRepo(db)
		require.NoError(t, taskRepo.Insert(ctx, "orphan task", "", 1, 7, userID))

		require.NoError(t, repo.Delete(ctx, userID))

		// user gone
		_, err := repo.GetById(ctx, userID)
		assert.ErrorIs(t, err, domain.ErrNoRecord)

		// and their tasks are gone too
		var count int
		require.NoError(t, db.GetContext(ctx, &count, "SELECT count(*) FROM tasks WHERE user_id = $1", userID))
		assert.Equal(t, 0, count)
	})

	t.Run("delete unknown user returns ErrNoRecord", func(t *testing.T) {
		err := repo.Delete(ctx, 99999999)
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})
}
