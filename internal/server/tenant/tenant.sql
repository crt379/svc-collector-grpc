CREATE TABLE tenant(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    describe VARCHAR(255),
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0)
);
