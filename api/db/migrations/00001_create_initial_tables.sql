-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       user_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                       username TEXT NOT NULL UNIQUE,
                       hashed_password TEXT NOT NULL,
                       email TEXT NOT NULL UNIQUE,
                       full_name TEXT,
                       bio TEXT,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
                        token_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                        user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                        token TEXT NOT NULL UNIQUE,
                        expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE statuses (
                          status_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                          uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                          user_id BIGINT,
                          label TEXT,
                          created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                          UNIQUE (label, user_id)
);

CREATE TABLE items (
                       item_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                       title TEXT NOT NULL,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lists (
                       list_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                       name TEXT NOT NULL,
                       user_id BIGINT NOT NULL,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE list_items (
                            list_item_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                            uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                            list_id BIGINT NOT NULL,
                            item_id BIGINT NOT NULL,
                            position INT NOT NULL,
                            prev_item_id BIGINT,
                            next_item_id BIGINT,
                            status_id BIGINT,
                            created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (list_id) REFERENCES lists (list_id) ON DELETE CASCADE,
                            FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE,
                            FOREIGN KEY (prev_item_id) REFERENCES list_items (list_item_id) ON DELETE SET NULL,
                            FOREIGN KEY (next_item_id) REFERENCES list_items (list_item_id) ON DELETE SET NULL,
                            FOREIGN KEY (status_id) REFERENCES statuses (status_id),
                            UNIQUE (list_id, item_id)
);

CREATE TABLE reviews (
                         review_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                         uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                         content TEXT NOT NULL,
                         user_id BIGINT NOT NULL,
                         item_id BIGINT NOT NULL,
                         created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                         FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE,
                         UNIQUE (user_id, item_id)
);

-- User indexes
CREATE INDEX idx_users_uuid ON users (uuid);

CREATE INDEX idx_users_created_date ON users (created_date);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

-- Review indexes
CREATE INDEX idx_reviews_uuid ON reviews (uuid);

CREATE INDEX idx_reviews_user_id ON reviews (user_id);

CREATE INDEX idx_reviews_item_id ON reviews (item_id);

CREATE INDEX idx_reviews_created_date ON reviews (created_date);

CREATE INDEX idx_reviews_user_id_created_date ON reviews (user_id, created_date);

-- Item indexes
CREATE INDEX idx_items_uuid ON items (uuid);

-- List indexes
CREATE INDEX idx_lists_uuid ON lists (uuid);

CREATE INDEX idx_lists_user_id_created_date ON lists (user_id, created_date);

-- List item indexes
CREATE INDEX idx_list_items_list_id_position ON list_items (list_id, position);

CREATE INDEX idx_list_items_list_id_item_id ON list_items (list_id, item_id);

CREATE INDEX idx_list_items_prev_item_id ON list_items (status_id, prev_item_id);

CREATE INDEX idx_list_items_next_item_id ON list_items (status_id, next_item_id);

CREATE INDEX idx_list_items_status_id ON list_items (list_id, status_id);

-- Status indexes
CREATE INDEX idx_statuses_uuid ON statuses (uuid);

CREATE INDEX idx_statuses_user_id ON statuses (user_id);

CREATE INDEX idx_statuses_created_date ON statuses (created_date);

CREATE INDEX idx_statuses_user_id_created_date ON statuses (user_id, created_date);

-- Refresh token indexes

CREATE INDEX idx_refresh_token_user_id ON refresh_tokens (user_id);

CREATE INDEX idx_refresh_token_token ON refresh_tokens (token);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop all indexes
DROP INDEX idx_refresh_token_user_id;

DROP INDEX idx_refresh_token_token;

DROP INDEX idx_users_uuid;

DROP INDEX idx_users_created_date;

DROP INDEX idx_users_username;

DROP INDEX idx_users_email;

DROP INDEX idx_reviews_uuid;

DROP INDEX idx_reviews_user_id;

DROP INDEX idx_reviews_item_id;

DROP INDEX idx_reviews_created_date;

DROP INDEX idx_reviews_user_id_created_date;

DROP INDEX idx_items_uuid;

DROP INDEX idx_statuses_uuid;

DROP INDEX idx_statuses_user_id;

DROP INDEX idx_statuses_created_date;

DROP INDEX idx_statuses_user_id_created_date;

DROP INDEX idx_lists_uuid;

DROP INDEX idx_list_items_list_id_position;

DROP INDEX idx_list_items_list_id_item_id;

DROP INDEX idx_list_items_prev_item_id;

DROP INDEX idx_list_items_next_item_id;

DROP INDEX idx_list_items_status_id;

-- Drop all tables
DROP TABLE refresh_tokens;

DROP TABLE reviews;

DROP TABLE list_items;

DROP TABLE items;

DROP TABLE statuses;

DROP TABLE lists;

DROP TABLE users;

-- +goose StatementEnd