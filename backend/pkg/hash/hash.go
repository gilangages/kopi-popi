package hash

import "golang.org/x/crypto/bcrypt"

// MakeHash meng-enkripsi string (misal: password murni) menjadi hash menggunakan algoritma bcrypt.
func MakeHash(password string) (string, error) {
	// bcrypt.DefaultCost bernilai 10. Angka ini cukup aman dan cepat.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckHash membandingkan password yang diinput user dengan hash yang tersimpan di Database.
// Akan mengembalikan true jika password cocok.
func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
