CREATE TABLE "connection"
(
    id            BIGSERIAL PRIMARY KEY,
    node_id       BIGINT                              NOT NULL,
    encrypted_key TEXT,
    pub_key       TEXT,
    "user"        VARCHAR   DEFAULT 'root'            NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT uniq_user_per_node UNIQUE ("user", node_id),
    CONSTRAINT fk_node FOREIGN KEY (node_id) REFERENCES node (id) ON DELETE CASCADE
);

ALTER SEQUENCE connection_id_seq RESTART WITH 100;

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON "connection"
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
