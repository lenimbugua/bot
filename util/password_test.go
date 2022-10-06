package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	require := require.New(t)
	password := RandomString(6)
	hashedPassword1, err := HashPassword(password)
	require.NoError(err)
	require.NotEmpty(hashedPassword1)
	err = CheckPassword(password, hashedPassword1)
	require.NoError(err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(err)
	require.NotEmpty(hashedPassword2)
	require.NotEqual(hashedPassword1, hashedPassword2)
}
