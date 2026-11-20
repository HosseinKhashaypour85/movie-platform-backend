# Users Module

ماژول مدیریت کاربران، پروفایل شخصی و امنیت حساب کاربری.

## قابلیت‌ها
- دریافت مشخصات حساب کاربری فعال (پروفایل) شامل نام، ایمیل، نقش و تاریخ عضویت.
- ویرایش اطلاعات هویتی کاربر (نام).
- تغییر کلمه عبور با تأیید پسورد فعلی و ابطال خودکار کلیه نشست‌های فعال در سایر دستگاه‌ها جهت حفظ امنیت.
- عدم افشای رمز عبور هش‌شده در خروجی‌های JSON (`json:"-"`).

## مدل دیتابیس (Users Table)
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user' NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL
);
```

## اندپوینت‌ها (API Endpoints)
| متد | مسیر | دسترسی | توضیحات |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/profile` | با توکن | دریافت اطلاعات پروفایل کاربر جاری |
| `PATCH` | `/api/v1/profile` | با توکن | ویرایش نام کاربر (بین ۲ تا ۱۰۰ کاراکتر) |
| `POST` | `/api/v1/profile/change-password` | با توکن | تغییر پسورد با وارد کردن پسورد فعلی و پسورد جدید |
