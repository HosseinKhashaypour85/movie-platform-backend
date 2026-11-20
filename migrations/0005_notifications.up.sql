CREATE TABLE IF NOT EXISTS notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(160) NOT NULL,
  body TEXT NOT NULL,
  type VARCHAR(40) NOT NULL DEFAULT 'general',
  action_url VARCHAR(2048) NULL,
  read_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_notifications_user_read (user_id,read_at),
  KEY idx_notifications_created (created_at),
  CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
