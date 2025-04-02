package dal

import (
	"db"
)

func EmailInUse(db *db.DB, email string) (bool, error) {
	var emailInUse bool
	row := db.QueryRow(`SELECT CASE
           WHEN EXISTS (SELECT 1 FROM users WHERE email=$1) THEN true
           ELSE false
       END;`, email)

	if err := row.Scan(&emailInUse); err != nil {
		return false, err
	}

	return emailInUse, nil
}

func UsernameInUse(db *db.DB, email string) (bool, error) {
	var usernameInUse bool
	row := db.QueryRow(`SELECT CASE
           WHEN EXISTS (SELECT 1 FROM users WHERE username=$1) THEN true
           ELSE false
       END;`, email)

	if err := row.Scan(&usernameInUse); err != nil {
		return false, err
	}

	return usernameInUse, nil
}

func CreateUser(db *db.DB, username, password_hash, email string) (int, error) {
	var id int
	row := db.QueryRow(`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id;`, username, email, password_hash)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func UserExist(db *db.DB, email, password_hash string) (bool, error) {
	var exist bool
	row := db.QueryRow(`SELECT CASE
           WHEN EXISTS (SELECT 1 FROM users WHERE email=$1 AND password_hash=$2) THEN true
           ELSE false
       END;`, email, password_hash)

	if err := row.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func GetUserId(db *db.DB, email string) (int, error) {
	var id int

	row := db.QueryRow(`SELECT id FROM users WHERE email=$1`, email)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func CreateEmailVerificationToken(db *db.DB, userId int, token string) error {
	row := db.QueryRow(`INSERT INTO email_tokens (user_id, token) VALUES ($1, $2)`, userId, token)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
