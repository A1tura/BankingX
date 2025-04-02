package utils

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func HashPassword(password string) (string, error) {
	hasher := crypto.SHA512.New()
	if _, err := hasher.Write([]byte(password)); err != nil {
		return "", err
	}

	passwordHash := hex.EncodeToString(hasher.Sum(nil))
	return passwordHash, nil
}

func PasswordLeaked(password string) (bool, int) {
	hasher := crypto.SHA1.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	prefix := strings.ToUpper(hash[:5])
	suffix := strings.ToUpper(hash[5:])

	req, err := http.Get("https://api.pwnedpasswords.com/range/" + prefix)
	if err != nil {
		return false, 0
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return false, 0
	}
	defer req.Body.Close()

	for _, line := range strings.Split(string(res), "\n") {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if parts[0] == suffix {
			leaked, err := strconv.Atoi(parts[1])
			if err != nil {
				return true, 0
			}
			return true, leaked
		}
	}

	return false, 0
}

func VerifyPasswordStrength(password string) bool {
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	return true
}

func GenerateEmailVerificationToken(email string) string {
	timestamp := time.Now().Unix()
	data := email + ":" + string(timestamp)

	h := hmac.New(sha512.New, []byte(os.Getenv("EMAIL_VERIF_SECRET")))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
