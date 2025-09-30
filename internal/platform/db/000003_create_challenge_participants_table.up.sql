CREATE TABLE IF NOT EXISTS challenge_participants (
    id SERIAL PRIMARY KEY,
    challenge_uuid UUID NOT NULL REFERENCES challenges(uuid) ON DELETE RESTRICT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    joined_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(challenge_uuid, user_id)
);