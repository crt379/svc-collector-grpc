CREATE TABLE app_svc_relation(
    uuid BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    aid BIGINT REFERENCES application(uuid) NOT NULL,
    sid BIGINT REFERENCES service(uuid) NOT NULL,
    create_time TIMESTAMP(0) NOT NULL,
    update_time TIMESTAMP(0),
    UNIQUE (aid, sid)
);