CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
    closed BIGINT DEFAULT 0,
    author_id INTEGER REFERENCES users(id) DEFAULT 0,
    assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
    title TEXT,
    content TEXT
);