CREATE TABLE IF NOT EXISTS edit_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER,
    edit_type TEXT,
    ts INTEGER,
    remote_file_mod_ts TEXT,
    key TEXT,      
    value TEXT,    
    new_value TEXT,
    old_value TEXT,
    FOREIGN KEY (job_id) REFERENCES job_table(id)
);