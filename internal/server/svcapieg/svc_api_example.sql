CREATE TABLE svc_api_example(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    aid BIGINT REFERENCES service_api(uuid) NOT NULL,
    tenant_id BIGINT REFERENCES tenant(uuid) NOT NULL,
    jid BIGINT REFERENCES jdata(uuid) NOT NULL
);