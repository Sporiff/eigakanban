-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       user_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       username TEXT NOT NULL,
                       hashed_password TEXT NOT NULL,
                       email TEXT NOT NULL,
                       full_name TEXT,
                       bio TEXT,
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE boards (
                        board_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                        name TEXT NOT NULL,
                        description TEXT,
                        user_id BIGINT,
                        FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                        created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE statuses (
                          status_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                          user_id BIGINT,
                          FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                          label TEXT,
                          created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
                       item_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                       title TEXT NOT NULL,
                       status_id BIGINT,
                       FOREIGN KEY (status_id) REFERENCES statuses (status_id),
                       created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE kbcolumns (
                           column_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                           name TEXT NOT NULL,
                           board_id BIGINT,
                           FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE,
                           user_id BIGINT,
                           FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                           position INT NOT NULL,
                           created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE column_items (
                              column_item_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                              column_id BIGINT,
                              FOREIGN KEY (column_id) REFERENCES kbcolumns (column_id) ON DELETE CASCADE,
                              item_id BIGINT,
                              FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE,
                              user_id BIGINT,
                              FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                              created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
                         review_id BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
                         content TEXT NOT NULL,
                         user_id BIGINT,
                         FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
                         item_id BIGINT,
                         FOREIGN KEY (item_id) REFERENCES items (item_id) ON DELETE CASCADE,
                         created_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE reviews;

DROP TABLE column_items;

DROP TABLE items;

DROP TABLE statuses;

DROP TABLE kbcolumns;

DROP TABLE boards;

DROP TABLE users;

-- +goose StatementEnd