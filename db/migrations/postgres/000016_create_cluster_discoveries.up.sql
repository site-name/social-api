CREATE TABLE IF NOT EXISTS cluster_discoveries (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    type varchar(64) NOT NULL,
    cluster_name varchar(64) NOT NULL,
    host_name varchar(512) NOT NULL,
    gossip_port integer NOT NULL,
    port integer NOT NULL,
    created_at bigint NOT NULL,
    last_ping_at bigint NOT NULL
);