CREATE TABLE jdata(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    data JSONB NOT NULL,
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    hash_type VARCHAR(255) NOT NULL,
    hash_value VARCHAR(255) NOT NULL,
    UNIQUE(hash_type, hash_value)
);