CREATE TABLE IF NOT EXISTS job_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_reference_key TEXT UNIQUE NOT NULL,
    content TEXT, -- JSON-Blob
    hash TEXT,
    deleted BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);