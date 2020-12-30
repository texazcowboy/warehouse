CREATE TABLE usr
(
    id   SERIAL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
)
