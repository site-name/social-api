CREATE TABLE IF NOT EXISTS menuitems (
  id character varying(36) NOT NULL PRIMARY KEY,
  menuid character varying(36),
  name character varying(128),
  parentid character varying(36),
  url character varying(256),
  categoryid character varying(36),
  collectionid character varying(36),
  pageid character varying(36),
  metadata jsonb,
  privatemetadata jsonb,
  sortorder integer
);
