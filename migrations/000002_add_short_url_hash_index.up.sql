CREATE INDEX IF NOT EXISTS idx_links_short_url_hash ON links USING HASH (short_url);
