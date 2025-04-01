CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username text UNIQUE NOT NULL,
    email text UNIQUE NOT NULL,
    email_verificated bool DEFAULT false,
    password_hash text NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE kyc (
    id SERIAL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,

    date_of_birth DATE NOT NULL,

    id_number INT NOT NULL,
    id_document TEXT NOT NULL,

    country TEXT NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL,

    status TEXT CHECK (status IN ('pending', 'verified', 'rejected')) DEFAULT 'pending',
    verified_at TIMESTAMP
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id),

    account_number TEXT UNIQUE NOT NULL,
    account_type TEXT CHECK (account_type IN ('checking', 'saving', 'buisness')) NOT NULL,

    balance DECIMAL(18, 2) DEFAULT 0.00 NOT NULL,
    currency TEXT CHECK (currency IN ('USD', 'EUR')) NOT NULL,

    status TEXT CHECK (status IN ('active', 'frozen', 'closed')) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT now()
);

INSERT INTO users (username, email, password_hash) VALUES ('test', 'test@test.com', 'test');
