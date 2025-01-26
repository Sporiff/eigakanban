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

CREATE TABLE boards (
                        board_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                        uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                        name TEXT NOT NULL,
                        description TEXT,
                        user_id BIGINT,
                        created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                        UNIQUE (name, user_id)
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
                       status_id BIGINT,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (status_id) REFERENCES statuses (status_id)
);

CREATE TABLE lists (
                       list_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                       name TEXT NOT NULL,
                       board_id BIGINT,
                       user_id BIGINT NOT NULL,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE,
                       FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                       UNIQUE (board_id, name)
);

CREATE TABLE list_items (
                            list_item_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                            uuid UUID DEFAULT gen_random_uuid () UNIQUE,
                            list_id BIGINT NOT NULL,
                            item_id BIGINT NOT NULL,
                            position INT NOT NULL,
                            prev_item_id BIGINT,
                            next_item_id BIGINT,
                            created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (list_id) REFERENCES lists (list_id) ON DELETE CASCADE,
                            FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE,
                            FOREIGN KEY (prev_item_id) REFERENCES list_items (list_item_id) ON DELETE SET NULL,
                            FOREIGN KEY (next_item_id) REFERENCES list_items (list_item_id) ON DELETE SET NULL,
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

-- Board indexes
CREATE INDEX idx_boards_uuid ON boards (uuid);

CREATE INDEX idx_boards_user_id ON boards (user_id);

CREATE INDEX idx_boards_created_date ON boards (created_date);

CREATE INDEX idx_boards_user_id_created_date ON boards (user_id, created_date);

CREATE INDEX idx_users_uuid ON users (uuid);

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

CREATE INDEX idx_lists_board_id_created_date ON lists (board_id, created_date);

-- List item indexes
CREATE INDEX idx_list_items_list_id_position ON list_items (list_id, position);

CREATE INDEX idx_list_items_list_id_item_id ON list_items (list_id, item_id);

CREATE INDEX idx_list_items_prev_item_id ON list_items (prev_item_id);

CREATE INDEX idx_list_items_next_item_id ON list_items (next_item_id);

-- Status indexes
CREATE INDEX idx_statuses_uuid ON statuses (uuid);

CREATE INDEX idx_statuses_user_id ON statuses (user_id);

CREATE INDEX idx_statuses_created_date ON statuses (created_date);

CREATE INDEX idx_statuses_user_id_created_date ON statuses (user_id, created_date);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop all indexes
DROP INDEX idx_users_uuid;

DROP INDEX idx_users_created_date;

DROP INDEX idx_users_username;

DROP INDEX idx_users_email;

DROP INDEX idx_boards_uuid;

DROP INDEX idx_boards_user_id;

DROP INDEX idx_boards_created_date;

DROP INDEX idx_boards_user_id_created_date;

DROP INDEX idx_users_uuid;

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

-- Drop all tables
DROP TABLE reviews;

DROP TABLE list_items;

DROP TABLE items;

DROP TABLE statuses;

DROP TABLE lists;

DROP TABLE boards;

DROP TABLE users;

-- +goose StatementEnd