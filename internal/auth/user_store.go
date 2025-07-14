package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type userData struct {
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

func GetUserDataFromGithub(accessToken string) (int, error) {
	var userID int

	requestURL := "https://api.github.com/user"
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return userID, fmt.Errorf("failed to create github validation request : %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/vnd.github5+json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return userID, fmt.Errorf("failed to send request to github : %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return userID, fmt.Errorf("failed to read response body : %v", err)
	}

	var user userData
	err = json.Unmarshal(body, &user)
	if err != nil {
		return userID, fmt.Errorf("failed to parse user data : %v", err)
	}
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		return userID, fmt.Errorf("failed to open database connection : %v", err)
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO users (github_id,name,avatar_url) VALUES (?,?,?)
	ON CONFLICT(github_id) DO UPDATE SET name=excluded.name,avatar_url=excluded.avatar_url`,
		user.ID, user.Name, user.AvatarURL)
	if err != nil {
		return userID, fmt.Errorf("failed to create/update user : %v", err)
	}
	err = db.QueryRow("SELECT id FROM users WHERE github_id = ?", user.ID).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("failed to fetch user id from database : %v", err)
	}
	return userID, nil
}

func GetUserDataFromDB(userID int) (userData, error) {
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		return userData{}, fmt.Errorf("failed to open database connection : %v", err)
	}
	defer db.Close()

	var user userData
	err = db.QueryRow("SELECT avatar_url,name FROM users WHERE id = ?", userID).Scan(&user.AvatarURL, &user.Name)
	if err != nil {
		return userData{}, fmt.Errorf("failed to fetch user data from database : %v", err)
	}
	user.ID = userID
	return user, nil
}
