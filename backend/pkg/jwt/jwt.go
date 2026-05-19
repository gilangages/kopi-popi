package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat JWT token yang berisi informasi user (Claims).
// Biasanya dipanggil ketika user berhasil Login.
func GenerateToken(userID string, name string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("fatal: JWT_SECRET environment variable is not set")
	}

	// Payload / Isi Token
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expired dalam 24 jam
	}

	// Membuat token dengan metode algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Me-return token berbentuk string (Signed)
	return token.SignedString([]byte(secret))
}

// ValidateToken menerima token string dari HTTP Header, memvalidasi signature-nya,
// dan mengembalikan isi payload (Claims) jika sukses.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi apakah algoritma yang dipakai saat membuat token adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Mengekstrak payload/claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
