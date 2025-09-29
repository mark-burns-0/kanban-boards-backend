CREATE TABLE IF NOT EXISTS challenge_metrics(
    id SERIAL PRIMARY KEY,
    challenge_uuid UUID NOT NULL REFERENCES challenges(uuid) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    unit TEXT,
    description TEXT
);