package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func createRandomChannel(t *testing.T) Channel {
	require := require.New(t)
	name := util.RandomString(6)
	channel, err := testQueries.CreateChannel(context.Background(), name)
	require.NoError(err)
	require.NotEmpty(channel)
	require.Equal(name, channel.Name)
	require.NotZero(channel.CreatedAt)
	require.NotZero(channel.UpdatedAt)
	return channel
}

func TestCreateChannel(t *testing.T) {
	createRandomChannel(t)
}

func TestGetChannel(t *testing.T) {
	require := require.New(t)
	channel1 := createRandomChannel(t)
	channel2, err := testQueries.GetChannel(context.Background(), channel1.Name)
	require.NoError(err)
	require.NotEmpty(channel2)
	require.Equal(channel1.Name, channel2.Name)
	require.WithinDuration(channel1.CreatedAt, channel2.CreatedAt, time.Second)
	require.WithinDuration(channel1.UpdatedAt, channel2.UpdatedAt, time.Second)
}

func TestUpdateChannel(t *testing.T) {
	require := require.New(t)

	channel := createRandomChannel(t)

	name := util.RandomString(6)
	arg := UpdateChannelParams{
		Name: sql.NullString{String: name, Valid: true},
		ID:   channel.ID,
	}
	updatedChannel, err := testQueries.UpdateChannel(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(updatedChannel)
	require.NotEqual(channel.Name, updatedChannel.Name)
	require.Equal(name, updatedChannel.Name)
	require.NotEqual(channel.UpdatedAt, updatedChannel.UpdatedAt)
	require.WithinDuration(channel.UpdatedAt, updatedChannel.UpdatedAt, 30*time.Second)
	require.WithinDuration(channel.CreatedAt, updatedChannel.CreatedAt, time.Second)
}

func TestDeleteChannel(t *testing.T) {
	require := require.New(t)
	channel := createRandomChannel(t)

	err := testQueries.DeleteChannel(context.Background(), channel.ID)
	require.NoError(err)
	channel1, err := testQueries.GetChannel(context.Background(), channel.Name)
	require.Error(err)
	require.EqualError(err, sql.ErrNoRows.Error())
	require.Empty(channel1)
}

func TestListChannels(t *testing.T) {
	require := require.New(t)
	for i := 0; i < 10; i++ {
		createRandomChannel(t)
	}
	arg := ListChannelsParams{
		Limit:  5,
		Offset: 5,
	}

	channels, err := testQueries.ListChannels(context.Background(), arg)
	require.NoError(err)
	require.Equal(len(channels), 5)

	for _, channel := range channels {
		require.NotEmpty(channel)
	}
}
