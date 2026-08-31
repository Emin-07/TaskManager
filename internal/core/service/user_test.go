package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	dbmodel "github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/testutil"
)

func newUserService(repo *testutil.MockUserRepo) UserServ {
	return NewUserService(repo)
}

func TestUserServ_Authenticate(t *testing.T) {
	t.Run("returns user on valid creds", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("Authenticate", mock.Anything, "bob@example.com").Return(
			&dbmodel.UserDb{
				Id:           2,
				Username:     "bob",
				Role:         "user",
				Email:        "bob@example.com",
				PasswordHash: bcryptHash(t, "password1"),
			}, nil)
		svc := newUserService(repo)

		user, err := svc.Authenticate(context.Background(), "bob@example.com", "password1")
		require.NoError(t, err)
		assert.Equal(t, "bob", user.Username)
	})

	t.Run("rejects wrong password", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("Authenticate", mock.Anything, "bob@example.com").Return(
			&dbmodel.UserDb{Email: "bob@example.com", PasswordHash: bcryptHash(t, "password1")}, nil)
		svc := newUserService(repo)

		_, err := svc.Authenticate(context.Background(), "bob@example.com", "wrongpass")
		assert.Error(t, err)
	})

	t.Run("propagates ErrNoRecord for unknown email", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("Authenticate", mock.Anything, "nobody@example.com").Return(nil, domain.ErrNoRecord)
		svc := newUserService(repo)

		_, err := svc.Authenticate(context.Background(), "nobody@example.com", "password1")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})
}

func TestUserServ_GetById(t *testing.T) {
	t.Run("returns mapped user", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("GetById", mock.Anything, 5).Return(
			&dbmodel.UserDb{Id: 5, Username: "carol", Role: "user", Email: "carol@example.com"}, nil)
		svc := newUserService(repo)

		user, err := svc.GetById(context.Background(), "5")
		require.NoError(t, err)
		assert.Equal(t, "carol", user.Username)
	})

	t.Run("rejects non-numeric id", func(t *testing.T) {
		svc := newUserService(&testutil.MockUserRepo{})
		_, err := svc.GetById(context.Background(), "abc")
		assert.Error(t, err)
	})
}

func TestUserServ_Insert(t *testing.T) {
	t.Run("defaults role to user", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("Insert", mock.Anything, mock.Anything, "user", "dave@example.com", mock.Anything).Return(nil)
		svc := newUserService(repo)

		err := svc.Insert(context.Background(), &domain.SignupUser{Email: "dave@example.com", Password: "password1", Username: "dave"})
		assert.NoError(t, err)
	})

	t.Run("generates username when empty", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		svc := newUserService(repo)
		repo.On("Insert", mock.Anything, mock.MatchedBy(func(u string) bool { return u != "" }), "user", "eve@example.com", mock.Anything).Return(nil)

		err := svc.Insert(context.Background(), &domain.SignupUser{Email: "eve@example.com", Password: "password1"})
		assert.NoError(t, err)
	})

	t.Run("rejects invalid email", func(t *testing.T) {
		svc := newUserService(&testutil.MockUserRepo{})
		err := svc.Insert(context.Background(), &domain.SignupUser{Email: "bademail", Password: "password1"})
		assert.Error(t, err)
	})

	t.Run("rejects short password", func(t *testing.T) {
		svc := newUserService(&testutil.MockUserRepo{})
		err := svc.Insert(context.Background(), &domain.SignupUser{Email: "a@b.com", Password: "short"})
		assert.Error(t, err)
	})
}

func TestUserServ_Delete(t *testing.T) {
	t.Run("deletes user", func(t *testing.T) {
		repo := &testutil.MockUserRepo{}
		repo.On("Delete", mock.Anything, 4).Return(nil)
		svc := newUserService(repo)

		err := svc.Delete(context.Background(), "4")
		assert.NoError(t, err)
	})

	t.Run("rejects non-numeric id", func(t *testing.T) {
		svc := newUserService(&testutil.MockUserRepo{})
		err := svc.Delete(context.Background(), "xyz")
		assert.Error(t, err)
	})
}

func bcryptHash(t *testing.T, password string) []byte {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return hash
}
