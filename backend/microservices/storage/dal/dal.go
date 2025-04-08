package dal

import (
	"db"
	"log"
	"time"
)

func UploadDocumentMetadata(db *db.DB, userId int, documentType, path string) error {
	var recordId *int
	row := db.QueryRow(`SELECT
    CASE
        WHEN EXISTS (SELECT 1 FROM KYC_documents WHERE user_id=$1 AND type=$2)
        THEN (SELECT id FROM KYC_documents WHERE user_id=$1 AND type=$2 LIMIT 1)
        ELSE NULL
    END AS record_id;`, userId, documentType)

	if err := row.Scan(&recordId); err != nil {
		log.Fatal(err)
		return err
	}

	if recordId == nil {
		row := db.QueryRow(`INSERT INTO KYC_documents (user_id, path, type) VALUES ($1, $2, $3)`, userId, path, documentType)

		if row.Err() != nil {
			log.Fatal(row.Err())
			return row.Err()
		}
	} else {
		row := db.QueryRow(`UPDATE KYC_documents SET created_at=$1 WHERE id=$2`, time.Now(), *recordId)

		if row.Err() != nil {
			log.Fatal(row.Err())
			return row.Err()
		}
	}

	return nil
}
