CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE objects (
    id SERIAL PRIMARY KEY,
    header TEXT NOT NULL,
    body TEXT NOT NULL,
    image TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id INTEGER NOT NULL REFERENCES users(id),
    date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_objects_user_id ON objects(user_id);
CREATE INDEX idx_objects_price ON objects(price);
CREATE INDEX idx_objects_date ON objects(date);