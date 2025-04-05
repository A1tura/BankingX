package dal

import (
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

func CreateKYC(db *db.DB, userId int, firstName, middleName, lastName string, dateOfBirth time.Time, phoneNumber, idNumber, idFront, idBack, selfie, country, state, city, address, postalCode string) error {
	row := db.QueryRow("INSERT INTO kyc (user_id, first_name, middle_name, last_name, date_of_birth, phone_number, id_number, id_front, id_back, selfie, country, state, city, address, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);", userId, firstName, middleName, lastName, dateOfBirth, phoneNumber, idNumber, idFront, idBack, selfie, country, state, city, address, postalCode)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
