CREATE TABLE IF NOT EXISTS ad_campaigns (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(160) NOT NULL,
  placement VARCHAR(40) NOT NULL,
  format VARCHAR(20) NOT NULL,
  media_url VARCHAR(2048) NOT NULL,
  click_url VARCHAR(2048) NULL,
  advertiser VARCHAR(160) NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  starts_at DATETIME(3) NOT NULL,
  ends_at DATETIME(3) NOT NULL,
  impression_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
  click_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_ad_campaigns_placement_active (placement,is_active),
  KEY idx_ad_campaigns_schedule (starts_at,ends_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ad_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  campaign_id BIGINT UNSIGNED NOT NULL,
  event_type VARCHAR(16) NOT NULL,
  placement VARCHAR(40) NOT NULL,
  created_at DATETIME(3) NULL,
  PRIMARY KEY (id),
  KEY idx_ad_events_campaign_date (campaign_id,created_at),
  KEY idx_ad_events_type_date (event_type,created_at),
  CONSTRAINT fk_ad_events_campaign FOREIGN KEY (campaign_id) REFERENCES ad_campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
