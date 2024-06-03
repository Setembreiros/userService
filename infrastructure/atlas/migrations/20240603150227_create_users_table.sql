CREATE TABLE userservice.users (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(255) UNIQUE,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE
);
