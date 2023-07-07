CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    location TEXT,
    alias TEXT NOT NULL UNIQUE,
    expiresat TIMESTAMP DEFAULT NULL,
    isexpired BOOLEAN DEFAULT FALSE,
    userid BIGINT DEFAULT NULL,
    createdat TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    gmail VARCHAR(50) UNIQUE,
    password VARCHAR(100),
    authkey  VARCHAR(100) UNIQUE,
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW()
)