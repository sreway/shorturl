CREATE TABLE IF NOT EXISTS urls
(
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    original_url VARCHAR(255) NOT NULL,
    CONSTRAINT uniq_original_url UNIQUE (user_id, original_url)
);