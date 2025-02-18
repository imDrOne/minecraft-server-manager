CREATE TABLE "node"
(
    id         BIGSERIAL PRIMARY KEY,
    host       TEXT                                NOT NULL,
    port       INTEGER   DEFAULT 80                NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT unique_host_port UNIQUE (host, port)
);

ALTER SEQUENCE node_id_seq RESTART WITH 100;

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON "node"
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
