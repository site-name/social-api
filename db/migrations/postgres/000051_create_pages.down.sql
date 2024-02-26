DROP INDEX IF EXISTS idx_pages_metadata;
DROP INDEX IF EXISTS idx_pages_private_metadata;
DROP INDEX IF EXISTS idx_pages_slug;
DROP INDEX IF EXISTS idx_pages_title;
DROP INDEX IF EXISTS idx_pages_title_lower_textpattern;
DROP INDEX IF EXISTS page_search_gin;

DROP TABLE IF EXISTS pages;