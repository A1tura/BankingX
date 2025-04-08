package controllers

import (
	"encoding/json"
	"error"
	"fmt"
	"middlewares"
	"net/http"
	storage "storage/Storage"
	"storage/types"

	"github.com/minio/minio-go/v7"
)

func Document(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		errors := error.NewError(true, w)
		authInfo := middlewares.GetAuth(r.Context())
		minioClient := r.Context().Value("minio").(*minio.Client)
		var req types.Upload

		// verify user
		// TODO: implement email verification logic
		if !authInfo.IsAuth {
			errors.NewError("You must be authenticated to access this resource.")
			errors.ThrowError()
			return
		}

		if authInfo.EmailConfirmed == nil {
			errors.ThrowInternalError()
			return
		}

		if !*authInfo.EmailConfirmed {
			errors.NewError("Your email address is not yet confirmed. Please verify your email before accessing this resource.")
			errors.ThrowError()
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.NewError("Invalid request body.")
			errors.ThrowError()
			return
		}
		defer r.Body.Close()

		if req.Type != "id_front" && req.Type != "id_back" {
			errors.NewError("Invalid document type.")
			errors.ThrowError()
			return
		}

		// save file to minio
		err := storage.UploadDocument(minioClient, req.Type, req.Document, authInfo.UserId)
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		fmt.Fprint(w, "ok")
	} else {
		w.WriteHeader(404)
		return
	}
}
