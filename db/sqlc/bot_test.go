package db

import (
	"context"
	"testing"
	"time"

	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func createRandomBot(t *testing.T) Bot {
	require := require.New(t)
	company := createRandomCompany(t)
	arg := CreateBotParams{
		Title:     util.RandomString(6),
		CompanyID: company.ID,
	}
	bot, err := testQueries.CreateBot(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(bot)
	require.Equal(arg.Title, bot.Title)
	require.Equal(arg.CompanyID, bot.CompanyID)
	require.NotZero(bot.CreatedAt)
	require.NotZero(bot.UpdatedAt)
	return bot
}

func TestCreateBot(t *testing.T) {
	createRandomBot(t)
}

func TestGetBot(t *testing.T) {
	require := require.New(t)
	bot1 := createRandomBot(t)
	bot2, err := testQueries.GetBot(context.Background(), bot1.ID)
	require.NoError(err)
	require.NotEmpty(bot2)
	require.Equal(bot1.Title, bot2.Title)
	require.Equal(bot1.CompanyID, bot2.CompanyID)
	require.WithinDuration(bot1.CreatedAt, bot2.CreatedAt, time.Second)
	require.WithinDuration(bot1.UpdatedAt, bot2.UpdatedAt, time.Second)
}
