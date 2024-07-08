CREATE TABLE service_api(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    path VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL,
    describe VARCHAR(255),
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    sid BIGINT REFERENCES service(uuid) NOT NULL,
    tenant_id BIGINT REFERENCES tenant(uuid) NOT NULL,
    UNIQUE (sid, path, method)
);