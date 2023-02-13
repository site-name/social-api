CREATE TABLE IF NOT EXISTS customerevents (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  orderid character varying(36),
  userid character varying(36),
  parameters text
);

ALTER TABLE ONLY customerevents
    ADD CONSTRAINT fk_customerevents_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY customerevents
    ADD CONSTRAINT fk_customerevents_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
