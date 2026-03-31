package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"moderation/config"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID
	Email		string
	Roles		[]string
	CreatedAt	time.Time
}

type Statuses struct {
	Flagged		bool
	Silenced	bool
}

func GetUser(username string) (User, error) {
	url := fmt.Sprintf("%s/users/%s", config.AuthUrl, username)
	resp, err := http.Get(url)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(respBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserStatuses(ctx context.Context, db *sql.DB, user User) Statuses {
	return Statuses{
		Flagged: IsFlagged(ctx, db, user),
		Silenced: IsSilenced(ctx, db, user),
	}
}


