CREATE TABLE IF NOT EXISTS "user" (
    id uuid NOT NULL,
    names character varying(250) NOT NULL,
    last_names character varying(250) NOT NULL,
    email character varying(150) NOT NULL UNIQUE,
    password character varying(300) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    PRIMARY KEY (id)
);

ALTER TABLE "user" OWNER to postgres;

CREATE TABLE IF NOT EXISTS "food" (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    title character varying(150) NOT NULL,
    description character varying(250) NOT NULL,
    food_image character varying(300) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at time with time zone NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES "user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE "food" OWNER to postgres;