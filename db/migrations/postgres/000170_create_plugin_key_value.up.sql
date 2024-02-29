CREATE TABLE IF NOT EXISTS plugin_key_values (
    plugin_id VARCHAR(190) NOT NULL,
    pkey VARCHAR(150) NOT NULL,
    pvalue BYTEA,
    expire_at BIGINT DEFAULT 0,
    PRIMARY KEY (plugin_id, pkey)
);
