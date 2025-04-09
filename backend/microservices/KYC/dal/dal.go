package dal

import (
	"database/sql"
	"db"
	"time"
)

func EmailConfirmed(db *db.DB, userId int) (bool, error) {
	var confirmed bool

	row := db.QueryRow("SELECT email_verificated FROM users WHERE id=$1;", userId)

	if err := row.Scan(&confirmed); err != nil {
		return false, err
	}

	return confirmed, nil
}

func DocumentsUploaded(db *db.DB, userId int) (bool, error) {
	var count int

	rows, err := db.Query(`SELECT id FROM KYC_documents WHERE user_id=$1`, userId)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		count++
	}

	if count == 3 {
		return true, nil
	}

	return false, nil
}

func AlreadyVerificated(db *db.DB, userId int) (bool, error) {
	var verificated bool

	row := db.QueryRow(`SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM kyc WHERE user_id=$1) THEN TRUE
        ELSE FALSE
    END AS record_exists;`, userId)

	if err := row.Scan(&verificated); err != nil {
		return false, err
	}

	return verificated, nil
}

func CreateKYC(db *db.DB, userId int, firstName, middleName, lastName string, dateOfBirth time.Time, phoneNumber, idNumber, country, state, city, address, postalCode string) error {
	row := db.QueryRow("INSERT INTO kyc (user_id, first_name, middle_name, last_name, date_of_birth, phone_number, id_number, country, state, city, address, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);", userId, firstName, middleName, lastName, dateOfBirth, phoneNumber, idNumber, country, state, city, address, postalCode)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func KYCStatus(db *db.DB, userId int) (string, error) {
	var status string
	row := db.QueryRow("SELECT status FROM kyc WHERE user_id=$1", userId)

	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "NE", nil
		}
		return "", err
	}

	return status, nil
}
