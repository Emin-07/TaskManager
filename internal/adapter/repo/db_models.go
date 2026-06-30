package repo

import "time"

type TaskDb struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
	Expires   time.Time `db:"expires"`
	UserId    int       `db:"user_id"`
}

type UserDb struct {
	Id           int       `db:"id"`
	Username     string    `db:"username"`
	Role         string    `db:"role"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
