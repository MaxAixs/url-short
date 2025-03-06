CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_url TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_short_url ON urls(short_url);
