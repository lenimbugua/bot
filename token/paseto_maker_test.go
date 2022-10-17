package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/lenimbugua/bot/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	phone := util.RandomPhoneNumber()
	userID := util.RandInt(1, 1000)
	name := util.RandomString(6)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, payload, err := maker.CreateToken(phone, userID, name,duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, name, payload.Name)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, phone, payload.Phone)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	phone := util.RandomPhoneNumber()
	userID := util.RandInt(1, 1000)
	name := util.RandomString(6)
	token, payload, err := maker.CreateToken(phone, userID, name, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
