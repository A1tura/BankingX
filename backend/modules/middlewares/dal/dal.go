package dal

import "db"

func EmailConfirmed(db *db.DB, userId int) (bool, error) {
	var confirmed bool

	row := db.QueryRow("SELECT email_verificated FROM users WHERE id=$1;", userId)

	if err := row.Scan(&confirmed); err != nil {
		return false, err
	}

	return confirmed, nil
}

func KYCStatus(db *db.DB, userId int) (string, error) {
	var status string
	row := db.QueryRow("SELECT status FROM kyc WHERE user_id=$1", userId)

	if err := row.Scan(&status); err != nil {
		return "", err
	}

	return status, nil
}
