CREATE TABLE service(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    describe VARCHAR(255),
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    tenant_id BIGINT REFERENCES tenant(uuid) NOT NULL,
    UNIQUE (name, tenant_id)
);