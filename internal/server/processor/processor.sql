CREATE TABLE processor(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    addr VARCHAR(255) NOT NULL,
    weight integer,
    state VARCHAR(255),
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    aid BIGINT REFERENCES application(uuid) NOT NULL,
    tenant_id BIGINT REFERENCES tenant(uuid) NOT NULL,
    UNIQUE (addr, aid)
);