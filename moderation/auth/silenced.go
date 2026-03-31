package auth

import (
	"context"
	"database/sql"
)

func IsSilenced(ctx context.Context, db *sql.DB, user User) bool {
	var count int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM silenced_users WHERE uuid=$1;
	`, user.ID).Scan(&count)
	if count == 0 || err == sql.ErrNoRows || err != nil {
		return false
	} else {
		return true
	}
}

func SilenceUser(ctx context.Context, db *sql.DB, user User) error {
	if IsSilenced(ctx, db, user) {
		return nil
	}
	_, err := db.QueryContext(ctx, `
		INSERT INTO silenced_users VALUES ($1);
	`, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func UnSilenceUser(ctx context.Context, db *sql.DB, user User) error {
	if !IsSilenced(ctx, db, user) {
		return nil
	}
	_, err := db.QueryContext(ctx, `
		DELETE FROM silenced_users WHERE uuid=$1;
	`, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func ListSilencedUsers(ctx context.Context, db *sql.DB, ind int, count int) ([]User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uuid FROM silenced_users
		LIMIT $1
		OFFSET $2;
	`, ind, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

