CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    location TEXT,
    alias TEXT NOT NULL UNIQUE,
    expiresat TIMESTAMP DEFAULT NULL,
    isexpired BOOLEAN DEFAULT FALSE,
    userid BIGINT DEFAULT NULL,
    createdat TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT,
    gmail VARCHAR(100),
    password VARCHAR(100),
    authkey  VARCHAR(100),
    updateat TIMESTAMP,
    createdat TIMESTAMP DEFAULT NOW()
)