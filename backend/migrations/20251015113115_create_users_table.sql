-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    role VARCHAR(50) DEFAULT 'user',
    region VARCHAR(100),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (
    login,
    password,
    is_admin,
    role,
    region,
    first_name,
    last_name,
    email,
    created_at,
    updated_at
) VALUES (
             'admin',
             'admin-foton-dealer-dev',
             TRUE,
             'admin',
             'all-russia',
             'Admin',
             'User',
             'admin@dealer-platform.com',
             NOW(),
             NOW()
         ) ON CONFLICT (login) DO UPDATE SET
                                             password = EXCLUDED.password,
                                             is_admin = EXCLUDED.is_admin,
                                             role = EXCLUDED.role,
                                             region = EXCLUDED.region,
                                             first_name = EXCLUDED.first_name,
                                             last_name = EXCLUDED.last_name,
                                             email = EXCLUDED.email,
                                             updated_at = NOW();

-- Вставка обычного пользователя
INSERT INTO users (
    login,
    password,
    is_admin,
    role,
    region,
    first_name,
    last_name,
    email,
    created_at,
    updated_at
) VALUES (
             'user',
             'user123',
             FALSE,
             'user',
             'Central',
             'Regular',
             'User',
             'user@dealer-platform.com',
             NOW(),
             NOW()
         ) ON CONFLICT (login) DO UPDATE SET
                                             password = EXCLUDED.password,
                                             is_admin = EXCLUDED.is_admin,
                                             role = EXCLUDED.role,
                                             region = EXCLUDED.region,
                                             first_name = EXCLUDED.first_name,
                                             last_name = EXCLUDED.last_name,
                                             email = EXCLUDED.email,
                                             updated_at = NOW();

-- +goose Down
DROP TABLE IF EXISTS users;