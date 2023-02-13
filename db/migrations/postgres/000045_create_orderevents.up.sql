CREATE TABLE IF NOT EXISTS orderevents (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  type character varying(255),
  orderid character varying(36),
  parameters text,
  userid character varying(36)
);

ALTER TABLE ONLY orderevents
    ADD CONSTRAINT fk_orderevents_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY orderevents
    ADD CONSTRAINT fk_orderevents_users FOREIGN KEY (userid) REFERENCES users(id);
