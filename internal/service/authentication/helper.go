package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
)
import "google.golang.org/grpc/codes"

func statusUnimplemented(method string) error {
	return status.Errorf(codes.Unimplemented, "%s is not implemented yet", method)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
