# Auth & Sessions Module

ماژول احراز هویت، مدیریت توکن‌ها و چرخه حیات نشست‌های کاربری (Session Lifecycle).

## قابلیت‌ها
- ثبت‌نام کاربر با اعتبارسنجی ورودی و هش رمز عبور با **bcrypt** (حداقل ۱۲ و حداکثر ۷۲ بایت).
- ورود امن با صدور جفت توکن:
  - **Access Token:** توکن کوتاه‌مدت JWT با الگوریتم HS256 و اعتبارسنجی صادرکننده (Issuer) و انقضا (پیش‌فرض ۱۵ دقیقه).
  - **Refresh Token:** توکن امن تصادفی با قابلیت چرخش (Token Rotation) و ذخیره هش SHA-256 در دیتابیس.
- تمدید نشست (Token Refresh) همراه با باطل‌سازی توکن قبلی و صدور توکن جدید.
- خروج امن (Logout) از نشست فعلی یا خروج همزمان از تمامی دستگاه‌ها (Logout All).
- شناسایی مشخصات کاربر از جمله `User-Agent` و آدرس IP برای ردیابی نشست‌های فعال.

## مدل دیتابیس (Sessions Table)
```sql
CREATE TABLE sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,
    user_agent VARCHAR(255) NULL,
    ip_address VARCHAR(45) NULL,
    expires_at DATETIME(3) NOT NULL,
    revoked_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL,
    INDEX idx_sessions_user_active (user_id, revoked_at, expires_at)
);
```

## اندپوینت‌ها (API Endpoints)
| متد | مسیر | دسترسی | توضیحات |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | عمومی | ثبت‌نام کاربر جدید با نام، ایمیل و رمز عبور |
| `POST` | `/api/v1/auth/login` | عمومی | ورود و دریافت Access Token و Refresh Token |
| `POST` | `/api/v1/auth/refresh` | عمومی | تمدید Access Token با استفاده از Refresh Token |
| `POST` | `/api/v1/auth/logout` | با توکن | ابطال نشست و رفرش توکن فعلی |
| `POST` | `/api/v1/auth/logout-all` | با توکن | خروج از تمامی دستگاه‌ها و ابطال تمام نشست‌های کاربر |

## متغیرهای محیطی مرتبط
- `JWT_SECRET`: کلید رمزنگاری توکن‌های JWT (حداقل ۳۲ بایت تصادفی).
- `JWT_ISSUER`: صادرکننده معتبر توکن (پیش‌فرض: `movie-platform`).
- `JWT_TTL`: مدت زمان اعتبار Access Token (پیش‌فرض: `15m`).
