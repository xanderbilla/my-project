CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    user_id bigint NOT NULL,
    tags varchar(255)[] NOT NULL DEFAULT '{}',
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);