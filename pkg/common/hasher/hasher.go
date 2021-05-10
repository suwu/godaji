package hasher

import "golang.org/x/crypto/bcrypt"

func Make(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

func Check(value string, hashedValue string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value)); err != nil {
		return false, err
	}
	return true, nil
}
