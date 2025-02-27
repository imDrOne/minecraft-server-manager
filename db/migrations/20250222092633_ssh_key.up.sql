CREATE TABLE "connection"
(
    id         BIGSERIAL PRIMARY KEY,
    node_id    BIGINT                              NOT NULL,
    key        TEXT                                NOT NULL,
    checksum   VARCHAR                             NOT NULL,
    "user"     VARCHAR   DEFAULT 'root',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (node_id) REFERENCES node (id) ON DELETE CASCADE
);

ALTER SEQUENCE connection_id_seq RESTART WITH 100;

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON "connection"
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
