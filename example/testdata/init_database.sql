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

CREATE TABLE public.users_settings
(
    user_id      UUID PRIMARY KEY NOT NULL,
    startup_page CHARACTER VARYING(255) DEFAULT 'index',
    theme        JSONB,
    FOREIGN KEY (user_id) REFERENCES public.users (id)
        MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
);

