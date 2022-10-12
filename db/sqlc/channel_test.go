package db

import (
	"context"
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
