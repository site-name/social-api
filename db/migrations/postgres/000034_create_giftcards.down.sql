DROP INDEX IF EXISTS idx_giftcards_code;
DROP INDEX IF EXISTS idx_giftcards_metadata;
DROP INDEX IF EXISTS idx_giftcards_private_metadata;
DROP INDEX IF EXISTS idx_giftcards_tag;
DROP INDEX IF EXISTS unique_giftcard_id_tag_id;
DROP INDEX IF EXISTS unique_giftcard_tags_name;
DROP TABLE IF EXISTS giftcard_tag_giftcards;
DROP TABLE IF EXISTS giftcard_tags;

DROP TABLE IF EXISTS giftcards;