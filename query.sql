CREATE TYPE otp_type AS ENUM ('SIGNUP', 'FORGOT_PASSWORD', 'LOGIN');

CREATE TABLE IF NOT EXISTS otp (
    id SERIAL PRIMARY KEY,
    key VARCHAR(50),
    value VARCHAR(10),
    type otp_type,
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE,
    name VARCHAR(50),
    gmail VARCHAR(50) UNIQUE,
    phone VARCHAR(20),
    password VARCHAR(100),
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS credits (
    id SERIAL PRIMARY KEY,
    available INT DEFAULT 100,        
    user_id INT REFERENCES users(id) UNIQUE,
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TYPE transcation_for AS ENUM ('MAIL', 'PHONE', 'SHORTEN', 'PROFILE_BIO', 'ALL');
CREATE TYPE transcation_type AS ENUM ('CREDITED', 'DEBITED', 'DECLINED', 'PENDING', 'PROCESSING');

CREATE TABLE IF NOT EXISTS transcations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    tfor transcation_for,
    type transcation_type,
    value INT,
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    shorthandname VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS authkey (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    authkey VARCHAR(100) UNIQUE NOT NULL,
    workspace_id INT REFERENCES workspaces(id),
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TYPE workspace_role AS ENUM ('ADMIN', 'MEMBER', 'EDITOR');

CREATE TABLE user_workspace (
    user_id INT REFERENCES users(id),
    workspace_id INT REFERENCES workspaces(id),
    role workspace_role,
    PRIMARY KEY (user_id, workspace_id),
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TYPE clickaction_type AS ENUM ('MAIL', 'PHONE', 'SHORTEN', 'PROFILE_BIO');

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    type otp_type,
    location TEXT,
    alias TEXT NOT NULL UNIQUE,
    isqrcode BOOLEAN DEFAULT FALSE,
    isrecurring BOOLEAN DEFAULT TRUE,
    expiresat TIMESTAMP DEFAULT NULL,
    isexpired BOOLEAN DEFAULT FALSE,
    istempblock BOOLEAN DEFAULT FALSE,
    workspace_id INT REFERENCES workspaces(id),
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS clicks (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id) NOT NULL,
    clickedfrom TEXT DEFAULT NULL,
    updatedat TIMESTAMP DEFAULT NOW(),
    createdat TIMESTAMP DEFAULT NOW(),
    deletedat TIMESTAMP DEFAULT NULL
);