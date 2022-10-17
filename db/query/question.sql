-- name: CreateQuestion :one
INSERT INTO questions (
    question,
    type, 
    parent_id, 
    bot_id, 
    next_question_id, 
    updated_at 
) VALUES (
    $1,$2,$3,$4,$5,$6
) RETURNING *;
