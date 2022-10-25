package db

import (
	"context"
	"database/sql"
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

func TestListAllBots(t *testing.T) {
	require := require.New(t)
	for i := 0; i < 10; i++ {
		createRandomBot(t)
	}
	arg := ListAllBotsParams{
		Limit:  5,
		Offset: 5,
	}

	bots, err := testQueries.ListAllBots(context.Background(), arg)
	require.NoError(err)
	require.Equal(len(bots), 5)

	for _, bot := range bots {
		require.NotEmpty(bot)
	}
}

func createCompanyBot(t *testing.T, companyID int64) Bot {
	require := require.New(t)
	arg := CreateBotParams{
		Title:     util.RandomString(6),
		CompanyID: companyID,
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

func TestListCompanyBots(t *testing.T) {
	require := require.New(t)
	company := createRandomCompany(t)
	for i := 0; i < 10; i++ {
		createCompanyBot(t, company.ID)
	}
	for i := 0; i < 10; i++ {
		createRandomBot(t)
	}
	arg := ListCompanyBotsParams{
		CompanyID: company.ID,
		Limit:     5,
		Offset:    5,
	}

	companybots, err := testQueries.ListCompanyBots(context.Background(), arg)
	require.NoError(err)
	require.Equal(len(companybots), 5)

	for _, bot := range companybots {
		require.NotEmpty(bot)
	}

	allBotsArgs := ListAllBotsParams{
		Limit:  20,
		Offset: 20,
	}
	allBots, err := testQueries.ListAllBots(context.Background(), allBotsArgs)
	require.NoError(err)
	require.Equal(len(allBots), 20)

	for _, allBot := range allBots {
		require.NotEmpty(allBot)
	}
}

func TestUpdateBotAllFields(t *testing.T) {
	require := require.New(t)

	bot := createRandomBot(t)
	company := createRandomCompany(t)

	title := util.RandomString(6)
	arg := UpdateBotParams{
		Title:     sql.NullString{String: title, Valid: true},
		CompanyID: sql.NullInt64{Int64: company.ID, Valid: true},
		ID:        bot.ID,
	}
	updatedBot, err := testQueries.UpdateBot(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedBot)
	require.NotEqual(bot.Title, updatedBot.Title)
	require.Equal(title, updatedBot.Title)
	require.NotEqual(bot.CompanyID, updatedBot.CompanyID)
	require.Equal(company.ID, updatedBot.CompanyID)
	require.NotEqual(bot.UpdatedAt, updatedBot.UpdatedAt)
	require.WithinDuration(bot.UpdatedAt, updatedBot.UpdatedAt, 30*time.Second)
	require.WithinDuration(bot.CreatedAt, updatedBot.CreatedAt, time.Second)
}

func TestUpdateBotTitleOnly(t *testing.T) {
	require := require.New(t)

	bot := createRandomBot(t)
	title := util.RandomString(6)

	arg := UpdateBotParams{
		Title: sql.NullString{String: title, Valid: true},
		ID:    bot.ID,
	}
	updatedBot, err := testQueries.UpdateBot(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedBot)

	require.NotEqual(updatedBot.Title, bot.Title)
	require.Equal(title, updatedBot.Title)
	require.Equal(updatedBot.CompanyID, updatedBot.CompanyID)
	require.WithinDuration(bot.UpdatedAt, updatedBot.UpdatedAt, 30*time.Second)
	require.WithinDuration(bot.CreatedAt, updatedBot.CreatedAt, time.Second)
}
func TestUpdateBotCompanyIDOnly(t *testing.T) {
	require := require.New(t)
	company := createRandomCompany(t)

	bot := createRandomBot(t)

	arg := UpdateBotParams{
		CompanyID: sql.NullInt64{Int64: company.ID, Valid: true},
		ID:        bot.ID,
	}
	updatedBot, err := testQueries.UpdateBot(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedBot)

	require.NotEqual(updatedBot.CompanyID, bot.CompanyID)
	require.Equal(company.ID, updatedBot.CompanyID)
	require.Equal(updatedBot.Title, updatedBot.Title)
	require.WithinDuration(bot.UpdatedAt, updatedBot.UpdatedAt, 30*time.Second)
	require.WithinDuration(bot.CreatedAt, updatedBot.CreatedAt, time.Second)
}

func TestDeleteBot(t *testing.T) {
	require := require.New(t)
	bot := createRandomBot(t)

	err := testQueries.DeleteBot(context.Background(), bot.ID)
	require.NoError(err)
	bot1, err := testQueries.GetBot(context.Background(), bot.ID)
	require.Error(err)
	require.EqualError(err, sql.ErrNoRows.Error())
	require.Empty(bot1)
}
