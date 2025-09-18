CREATE TABLE public.users
(
    id         UUID PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    email      CHARACTER VARYING(128)   NOT NULL,
    name       CHARACTER VARYING(128),
    password   CHARACTER VARYING(128)   NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX users_email_unique ON users USING BTREE (email);

