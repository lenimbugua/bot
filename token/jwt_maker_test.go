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

	phone := util.RandomPhoneNumber()
	userID := util.RandInt(1, 1000)
	companyID := util.RandInt(1, 1000)
	name := util.RandomString(6)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(phone, userID, name, companyID, duration)
	require.NoError(err)
	require.NotEmpty(token)
	require.NotEmpty(payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(err)
	require.NotEmpty(token)

	require.NotZero(payload.ID)
	require.Equal(name, payload.Name)
	require.Equal(userID, payload.UserID)
	require.Equal(phone, payload.Phone)
	require.WithinDuration(issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require := require.New(t)
	require.NoError(err)

	phone := util.RandomPhoneNumber()
	userID := util.RandInt(1, 1000)
	companyID := util.RandInt(1, 1000)
	name := util.RandomString(6)

	token, payload, err := maker.CreateToken(phone, userID, name, companyID, -time.Minute)
	require.NoError(err)
	require.NotEmpty(token)
	require.NotEmpty(payload)

	payload, err = maker.VerifyToken(token)
	require.Error(err)
	require.EqualError(err, ErrExpiredToken.Error())
	require.Nil(payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	phone := util.RandomPhoneNumber()
	userID := util.RandInt(1, 1000)
	companyID := util.RandInt(1, 1000)
	name := util.RandomString(6)
	payload, err := NewPayload(phone, userID, name, companyID, time.Minute)
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
