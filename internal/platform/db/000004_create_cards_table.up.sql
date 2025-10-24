
CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    board_id UUID NOT NULL REFERENCES boards(id) ON DELETE RESTRICT,
    column_id INTEGER NOT NULL REFERENCES board_columns(id) ON DELETE RESTRICT,
    text TEXT NOT NULL,
    description TEXT,
    position INTEGER DEFAULT 0,
    properties JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
