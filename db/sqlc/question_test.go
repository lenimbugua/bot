package db

import (
	"context"
	"testing"

	"github.com/lenimbugua/bot/util"
	"github.com/stretchr/testify/require"
)

func createRandomQuestion(t *testing.T) Question {
	require := require.New(t)
	bot := createRandomBot(t)

	arg := CreateQuestionParams{
		Question: util.RandomString(6),
		Type:     util.RandomString(6),
		BotID:    bot.ID,
	}

	question, err := testQueries.CreateQuestion(context.Background(), arg)
	require.NoError(err)
	require.NotEmpty(question)
	require.Equal(arg.Question, question.Question)
	require.Equal(arg.Type, question.Type)
	require.Equal(arg.ParentID, question.ParentID)
	require.Equal(arg.BotID, question.BotID)
	require.Equal(arg.NextQuestionID, question.NextQuestionID)
	require.NotZero(question.CreatedAt)
	require.NotZero(question.UpdatedAt)
	return question
}

func TestCreateQuestion(t *testing.T) {
	createRandomQuestion(t)
}
