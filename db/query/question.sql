-- name: CreateQuestion :one
INSERT INTO questions (
    question,
    type, 
    parent_id, 
    bot_id, 
    next_question_id
) VALUES (
    $1,$2,$3,$4,$5
) RETURNING *;
