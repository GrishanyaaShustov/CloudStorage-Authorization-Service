-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id              BIGSERIAL PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    login           TEXT,                -- обычный логин (можно не уникальный)
    password_hash   TEXT,                -- хэш пароля, может быть NULL для чисто соц.аккаунтов
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,

    -- поля для социальной авторизации (nullable, но при наличии должны быть уникальны)
    github_id       TEXT UNIQUE,
    google_id       TEXT UNIQUE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
    )

-- +goose StatementEnd