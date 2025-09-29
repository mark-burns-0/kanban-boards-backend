CREATE TABLE IF NOT EXISTS user_challenge_progress(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    challenge_uuid UUID NOT NULL REFERENCES challenges(uuid),
    progress_count INT DEFAULT 0,        -- количество выполненного за указанный день
    progress_date DATE NOT NULL,         -- дата, когда был сделан этот прогресс
    status TEXT DEFAULT 'started',       -- started, in_progress, finished
    streak SMALLINT NOT NULL DEFAULT 0,
    UNIQUE (user_id, challenge_uuid, progress_date)
);