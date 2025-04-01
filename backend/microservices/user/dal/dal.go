package dal

import (
	"db"
	"fmt"
)

func EmailInUse(db db.DB, email string) (bool, bool) {
	var emailInUse bool
	row := db.QueryRow(`SELECT CASE
           WHEN EXISTS (SELECT 1 FROM users WHERE email=$1) THEN true
           ELSE false
       END;`, email)

    if err := row.Scan(&emailInUse); err != nil {
        fmt.Println(err)
        return false, false
    }

	return true, emailInUse
}
