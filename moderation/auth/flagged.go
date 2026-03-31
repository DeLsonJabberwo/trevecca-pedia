package auth

import (
	"context"
	"database/sql"
)

func IsFlagged(ctx context.Context, db *sql.DB, user User) bool {
	var count int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM flagged_users WHERE uuid=$1;
	`, user.ID).Scan(&count)
	if count == 0 || err == sql.ErrNoRows || err != nil {
		return false
	} else {
		return true
	}
}

func FlagUser(ctx context.Context, db *sql.DB, user User) error {
	if IsFlagged(ctx, db, user) {
		return nil
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO flagged_users VALUES ($1);
	`, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func UnFlagUser(ctx context.Context, db *sql.DB, user User) error {
	if !IsFlagged(ctx, db, user) {
		return nil
	}
	_, err := db.ExecContext(ctx, `
		DELETE FROM flagged_users WHERE uuid=$1;
	`, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func ListFlaggedUsers(ctx context.Context, db *sql.DB, ind int, count int) ([]User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uuid FROM flagged_users
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

