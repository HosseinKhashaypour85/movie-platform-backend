# Internal Core Architecture

لایه زیرساخت، تنظیمات و میدلورهای مرکزی سامانه (`internal`).

## اجزای اصلی

### ۱. تنظیمات و محیط (`internal/config`)
- اعتبارسنجی دقیق و سخت‌گیرانه متغیرهای محیطی در زمان اجرای اولیه برنامه.
- کنترل فرمت پورت، اعتبارسنجی مقادیر استخر دیتابیس (Pool)، بررسی انقضای توکن‌ها و طول کلید `JWT_SECRET` (حداقل ۳۲ بایت).
- اعتبارسنجی دامنه و لیست سفید CORS (`APP_ALLOWED_ORIGINS`).

### ۲. پایگاه داده و استخر اتصالات (`internal/database`)
- راه‌اندازی اتصال به MySQL با درایور استاندارد و نگاشت رابط GORM.
- پیکربندی پارامترهای Pool جهت جلوگیری از Connection Leak:
  - `SetMaxOpenConns`
  - `SetMaxIdleConns`
  - `SetConnMaxLifetime`
  - `SetConnMaxIdleTime`
- اجرای تست `PingContext` با مهلت زمانی ۵ ثانیه‌ای در شروع به کار سرور و اندپوینت `/health/ready`.

### ۳. قالب استاندارد پاسخ‌های HTTP (`internal/httpx`)
- متد `httpx.JSON(w, status, data)`: ارسال خروجی با هدر `application/json; charset=utf-8`.
- متد `httpx.Error(w, status, code, message)`: ارسال پاسخ خطای ساختاریافته به شکل:
```json
{
  "error": {
    "code": "validation_error",
    "message": "title, slug and IMDb rating are invalid"
  }
}
```

### ۴. میدلورهای سراسری برنامه (`internal/middleware/app.go`)
- **RequestID:** تخصیص شناسه یکتای ۱۶ بایتی تصادفی به هر درخواست و هدایت آن در پاسخ و کانتکست (`X-Request-ID`).
- **Recover:** رهگیری پنیک‌ها و ثبت خطا در لاگر ساختاریافته بدون کرش کردن پروسه سرور.
- **SecurityHeaders:** تنظیم هدرهای امنیتی `CSP`, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`.
- **BodyLimit:** محدودسازی اندازه بدنه درخواست (پیش‌فرض ۱ مگابایت) برای جلوگیری از حملات DoS.
- **Timeout:** مدیریت حداکثر زمان اجرای درخواست‌های ورودی (۱۵ ثانیه).
- **CORS:** کنترل دسترسی دامنه‌های مجاز و پاسخ به درخواست‌های Preflight.
- **RateLimit:** محدودساز درخواست‌ها مبتنی بر آدرس IP با الگوریتم Token Bucket در حافظه (پیش‌فرض ۱۲۰ درخواست در دقیقه).

### ۵. میدلورهای احراز هویت و سطوح دسترسی (`internal/middleware/auth.go`)
- **Auth:** اعتبارسنجی امضای JWT Bearer، تاریخ انقضا، صادرکننده و تزریق `UserID` و `Role` به `context.Context`.
- **RequireRole:** بررسی سطح دسترسی کاربر بر اساس نقش‌های مجاز (`admin`, `super_admin`).
