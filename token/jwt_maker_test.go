package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require := require.New(t)
	require.NoError(err)

	username := util.RandomString(6)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(err)
	require.NotEmpty(token)
	require.NotEmpty(payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(err)
	require.NotEmpty(token)

	require.NotZero(payload.ID)
	require.Equal(username, payload.Username)
	require.WithinDuration(issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require := require.New(t)
	require.NoError(err)

	token, payload, err := maker.CreateToken(util.RandomString(6), -time.Minute)
	require.NoError(err)
	require.NotEmpty(token)
	require.NotEmpty(payload)

	payload, err = maker.VerifyToken(token)
	require.Error(err)
	require.EqualError(err, ErrExpiredToken.Error())
	require.Nil(payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomString(6), time.Minute)
	require := require.New(t)
	require.NoError(err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(err)

	payload, err = maker.VerifyToken(token)
	require.Error(err)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(payload)
}
