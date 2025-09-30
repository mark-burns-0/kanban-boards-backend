CREATE TABLE IF NOT EXISTS challenges (
    uuid UUID       PRIMARY KEY,
    name            TEXT,
    user_id         BIGINT NOT NULL,
    type            TEXT,
    every_interval  INTERVAL,            -- период, например '15 minutes'
    started_at      TIMESTAMPTZ,
    end_at          TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT fk_challenges_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE NO ACTION
);
