INSERT INTO links (short_url, original_url) VALUES ($1, $2) ON CONFLICT DO NOTHING
