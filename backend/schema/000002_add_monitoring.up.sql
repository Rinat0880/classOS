CREATE TABLE device_status (
    device_name VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    last_heartbeat TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_device_last_heartbeat ON device_status(last_heartbeat);
CREATE INDEX idx_device_username ON device_status(username);

CREATE TABLE user_logs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    log_type VARCHAR(50),
    program VARCHAR(255),
    action TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_username ON user_logs(username);
CREATE INDEX idx_logs_device ON user_logs(device_name);
CREATE INDEX idx_logs_timestamp ON user_logs(timestamp);
CREATE INDEX idx_logs_created_at ON user_logs(created_at);
