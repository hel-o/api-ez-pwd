CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL ,
    password_hash VARCHAR(500) NOT NULL ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE secret_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_secrets(
    id SERIAL PRIMARY KEY,
    description VARCHAR(250),
    username VARCHAR(250),
    password_json JSONB,
    safe_note_json JSONB,
    url_site VARCHAR(250),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    category_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_categories FOREIGN KEY(category_id) REFERENCES secret_categories(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
