package tools

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(pwd string) (string, error) {
	pwdBytes := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	plainPwdBytes := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdBytes)
	if err != nil {
		return false
	}
	return true
}
