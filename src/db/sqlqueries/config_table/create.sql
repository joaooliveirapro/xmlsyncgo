CREATE TABLE IF NOT EXISTS config_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER
    root TEXT,
    job_node TEXT,
    external_reference_key TEXT,
    remote_filepath TEXT,
    remote_filename TEXT,
    download_to TEXT,
    sftp_host TEXT,
    sftp_port TEXT,
    sftp_username TEXT,
    sftp_password TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (client_id) REFERENCES client_table(id)
);