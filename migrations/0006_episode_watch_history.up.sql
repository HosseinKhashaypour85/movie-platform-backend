CREATE TABLE IF NOT EXISTS episode_watch_histories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  episode_id BIGINT UNSIGNED NOT NULL,
  progress_seconds BIGINT NOT NULL DEFAULT 0,
  duration_seconds BIGINT NOT NULL DEFAULT 0,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  watched_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY ux_user_episode_history (user_id,episode_id),
  KEY idx_episode_history_watched (user_id,watched_at),
  CONSTRAINT fk_episode_history_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_episode_history_episode FOREIGN KEY (episode_id) REFERENCES episodes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
