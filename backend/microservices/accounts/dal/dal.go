package dal

import (
	"database/sql"
	"db"
	"encoding/json"
)

func GetAccounts(db *db.DB, userId int) ([]Account, error) {
	var accounts []Account

	rows, err := db.Query("SELECT account_number, account_type, currency, balance, limits, status, created_at FROM accounts WHERE user_id=$1", userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return accounts, nil
		}
		return accounts, err
	}

	for rows.Next() {
		var account Account
		var limits []byte

		if err := rows.Scan(&account.AccountNumber, &account.AccountType, &account.Currency, &account.Balance, &limits, &account.Status, &account.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				return accounts, nil
			}
			return accounts, err
		}

		if err := json.Unmarshal(limits, &account.Limits); err != nil {
			return accounts, nil
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func AccountNumberExist(db *db.DB, accountNumber string) (bool, error) {
	var exist bool;
	row := db.QueryRow(`SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM kyc WHERE account_number=$1) THEN TRUE
        ELSE FALSE
    END AS record_exists;`, accountNumber)

	if err := row.Scan(exist); err != nil {
		return exist, err
	}

	return exist, nil
}

func CreateAccount(db *db.DB, userId int, accountType, accountNumber, currency string, isPrimary bool, limits AccountLimits) (bool, error) {
	jsonLimits, err := json.Marshal(limits)
	if err != nil {
		return false, err
	}

	row := db.QueryRow(`INSERT INTO accounts (user_id, is_primary, limits, account_number, account_type, currency) VALUES ($1, $2, $3, $4, $5, $6)`, userId, isPrimary, jsonLimits, accountNumber, accountType, currency)
	if row.Err() != nil {
		return false, row.Err()
	}

	return true, nil
}
