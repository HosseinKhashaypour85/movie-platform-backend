CREATE TABLE IF NOT EXISTS transcode_jobs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  source_asset_id BIGINT UNSIGNED NOT NULL,
  status VARCHAR(20) NOT NULL,
  attempts TINYINT UNSIGNED NOT NULL DEFAULT 0,
  error_message VARCHAR(1000) NULL,
  locked_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_transcode_jobs_source_asset (source_asset_id),
  KEY idx_transcode_jobs_status (status),
  CONSTRAINT fk_transcode_jobs_source_asset FOREIGN KEY (source_asset_id) REFERENCES media_assets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
