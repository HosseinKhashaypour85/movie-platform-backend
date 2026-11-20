# Movie Platform Backend - API Docs & Bruno Collection

مستندات رسمی، جامع و کالکشن کامل **Bruno** و **OpenAPI 3.0.3** برای وب‌سرویس پلتفرم پخش فیلم و سریال (Movie Platform Backend).
این مستندات شامل **نمونه کامل Request و Response به صورت JSON برای تک‌تک ۶۸ اندپوینت** سامانه همراه با کدهای آماده cURL، هدرها، پارامترها و توضیحات فارسی می‌باشد.

---

## 📁 ساختار ریپازیتوری (Repository Structure)

```text
├── openapi.yaml                  # مستندات کامل استاندارد OpenAPI 3.0.3 (YAML)
├── endpoints.json                # کاتالوگ جامع تمام ۶۸ اندپوینت به صورت یک فایل JSON یکپارچه
├── endpoints/                    # فایل‌های JSON استاندارد تک‌تک اندپوینت‌ها (شامل مقادیر نمونه و نحوه اجرا)
├── bruno.json                    # فایل کانفیگ کالکشن اصلی نرم‌افزار برونو در ریشه پروژه
├── environments/
│   └── Local.bru                 # متغیرهای محیطی برونو (baseUrl, token, workerToken)
├── Auth/                         # پوشه کالکشن برونو برای احراز هویت
├── Catalog/                      # پوشه کالکشن برونو برای کاتالوگ فیلم و سریال
├── Playback/                     # پوشه کالکشن برونو برای استریم و زیرنویس
├── System/                       # پوشه کالکشن برونو برای سلامت سیستم
├── ...                           # تمام ۲۰ پوشه کالکشن به صورت مستقیم در ریشه
├── .gitignore
└── README.md
```

---

## ⚡ راهنمای استفاده و لود کالکشن در Bruno (Quickstart)

این ریپازیتوری به صورت **کالکشن استاندارد Bruno در ریشه پروژه** تنظیم شده است:

### روش اول: کلون مستقیم در Bruno با آدرس گیت (توصیه شده)
1. در نرم‌افزار **Bruno** روی گزینه **Clone Git Repository** کلیک کنید.
2. آدرس گیت ریپازیتوری را وارد نمایید:
   ```text
https://github.com/HosseinKhashaypour85/movie-platform-backend-api-docs.git
```
3. تمام ۲۰ پوشه و ۶۸ ریکوئست مستقیماً در سایدبار Bruno نمایش داده می‌شوند.
4. از گوشه بالا سمت راست، محیط **Local** را انتخاب کنید.

### روش دوم: باز کردن مستقیم پوشه پروژه
1. در نرم‌افزار Bruno روی **Open Collection** کلیک کنید.
2. پوشه دانلود شده ریپازیتوری را انتخاب نمایید (همان پوشه‌ای که فایل `bruno.json` در آن قرار دارد).

---

## ⚙️ متغیرهای محیطی در برونو (Environments)

فایل `environments/Local.bru` شامل متغیرهای زیر است:

| متغیر | مقدار پیش‌فرض | توضیحات |
| :--- | :--- | :--- |
| `baseUrl` | `http://localhost:8080` | آدرس سرور بک‌اند |
| `token` | *(خالی)* | توکن JWT کاربر یا ادمین (Bearer Token) |
| `workerToken` | *(خالی)* | توکن اختصاصی ورکر هوش مصنوعی (`SUBTITLE_WORKER_TOKEN`) |

> **نکته احراز هویت:** پس از فراخوانی اندپوینت `POST /api/v1/auth/login`، مقدار `access_token` دریافتی را در متغیر `token` محیط Local قرار دهید تا تمام روت‌های محافظت‌شده به طور خودکار احراز هویت شوند.

---

## 📑 فهرست دسترسی سریع به بخش‌های API (Table of Contents)

- [۱. سلامت و وضعیت سیستم (System & Health)](#۱-سلامت-و-وضعیت-سیستم-system--health) (3 اندپوینت)
- [۲. احراز هویت و مدیریت نشست‌ها (Auth & Sessions)](#۲-احراز-هویت-و-مدیریت-نشست‌ها-auth--sessions) (5 اندپوینت)
- [۳. پروفایل کاربر و مدیریت حساب (User Profile)](#۳-پروفایل-کاربر-و-مدیریت-حساب-user-profile) (3 اندپوینت)
- [۴. کاتالوگ عمومی فیلم و سریال (Catalog - Public)](#۴-کاتالوگ-عمومی-فیلم-و-سریال-catalog---public) (6 اندپوینت)
- [۵. پخش ویدیو و زیرنویس‌ها (Playback & Subtitles)](#۵-پخش-ویدیو-و-زیرنویس‌ها-playback--subtitles) (2 اندپوینت)
- [۶. تبلیغات عمومی و رهگیری رویدادها (Ads - Public)](#۶-تبلیغات-عمومی-و-رهگیری-رویدادها-ads---public) (2 اندپوینت)
- [۷. تعاملات: امتیازها و برگزیده‌ها (Engagement)](#۷-تعاملات:-امتیازها-و-برگزیده‌ها-engagement) (5 اندپوینت)
- [۸. تاریخچه و پیشرفت تماشا (Watch History)](#۸-تاریخچه-و-پیشرفت-تماشا-watch-history) (7 اندپوینت)
- [۹. اعلان‌های کاربری (User Notifications)](#۹-اعلان‌های-کاربری-user-notifications) (3 اندپوینت)
- [۱۰. داشبورد و آمار مدیریتی (Admin Dashboard & Analytics)](#۱۰-داشبورد-و-آمار-مدیریتی-admin-dashboard--analytics) (3 اندپوینت)
- [۱۱. مدیریت کاربران (Admin Users)](#۱۱-مدیریت-کاربران-admin-users) (3 اندپوینت)
- [۱۲. مدیریت ژانرها (Admin Genres)](#۱۲-مدیریت-ژانرها-admin-genres) (3 اندپوینت)
- [۱۳. مدیریت فیلم‌ها (Admin Movies)](#۱۳-مدیریت-فیلم‌ها-admin-movies) (3 اندپوینت)
- [۱۴. مدیریت سریال‌ها (Admin Series)](#۱۴-مدیریت-سریال‌ها-admin-series) (3 اندپوینت)
- [۱۵. مدیریت فصل‌ها (Admin Seasons)](#۱۵-مدیریت-فصل‌ها-admin-seasons) (3 اندپوینت)
- [۱۶. مدیریت قسمت‌ها (Admin Episodes)](#۱۶-مدیریت-قسمت‌ها-admin-episodes) (3 اندپوینت)
- [۱۷. آپلود مستقیم به S3 و است‌های رسانه‌ای (Admin Media & S3)](#۱۷-آپلود-مستقیم-به-s3-و-است‌های-رسانه‌ای-admin-media--s3) (5 اندپوینت)
- [۱۸. مدیریت کمپین‌های تبلیغاتی (Admin Ads)](#۱۸-مدیریت-کمپین‌های-تبلیغاتی-admin-ads) (3 اندپوینت)
- [۱۹. ارسال اعلان همگانی (Admin Broadcast Notifications)](#۱۹-ارسال-اعلان-همگانی-admin-broadcast-notifications) (1 اندپوینت)
- [۲۰. ورکر داخلی پردازش زیرنویس (Internal Subtitle Worker)](#۲۰-ورکر-داخلی-پردازش-زیرنویس-internal-subtitle-worker) (2 اندپوینت)

---

## ۱. سلامت و وضعیت سیستم (System & Health)

### 1.1. Root Service Status
- **متد و آدرس:** `GET /`
- **توضیحات:** بررسی وضعیت آنلاین بودن سرویس اصلی بک‌اند (Root Diagnostics)
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "service": "movie-platform-api",
  "status": "running"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/"
```

---

### 1.2. Liveness Probe
- **متد و آدرس:** `GET /health/live`
- **توضیحات:** بررسی زنده بودن پروسس سرور برای کوبرنتیز و داکر (Liveness Probe)
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "status": "alive"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/health/live"
```

---

### 1.3. Readiness Probe
- **متد و آدرس:** `GET /health/ready`
- **توضیحات:** بررسی آماده به کار بودن و اتصال به دیتابیس MySQL (Readiness Probe)
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "database": "connected",
  "status": "ready"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/health/ready"
```

---

## ۲. احراز هویت و مدیریت نشست‌ها (Auth & Sessions)

### 2.1. Register User
- **متد و آدرس:** `POST /api/v1/auth/register`
- **توضیحات:** ثبت‌نام کاربر جدید با نام، ایمیل یکتا و رمز عبور قوی و صدور توکن‌های دسترسی
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **هدرهای ارسالی:**
```http
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "email": "hosein@example.com",
  "name": "Hosein Khashayarpour",
  "password": "SuperSecurePass123!"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `email` | `string` | بله | ایمیل معتبر و یکتا |
| `name` | `string` | بله | نام کامل کاربر (حداکثر ۱۰۰ کاراکتر) |
| `password` | `string` | بله | رمز عبور با حداقل ۸ کاراکتر |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "refresh_token": "b472e3a1-7c89-4d82-95f0-6c9f28d8b412",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ...",
  "user": {
    "email": "hosein@example.com",
    "id": 1,
    "is_active": true,
    "name": "Hosein Khashayarpour",
    "role": "user"
  }
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hosein Khashayarpour","email":"hosein@example.com","password":"SuperSecurePass123!"}'
```

---

### 2.2. Login User
- **متد و آدرس:** `POST /api/v1/auth/login`
- **توضیحات:** ورود کاربر به سامانه با ایمیل و رمز عبور و دریافت JWT Access Token و Refresh Token
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **هدرهای ارسالی:**
```http
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "email": "hosein@example.com",
  "password": "SuperSecurePass123!"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `email` | `string` | بله | ایمیل ثبت‌شده |
| `password` | `string` | بله | کلمه عبور حساب |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "refresh_token": "b472e3a1-7c89-4d82-95f0-6c9f28d8b412",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ...",
  "user": {
    "email": "hosein@example.com",
    "id": 1,
    "is_active": true,
    "name": "Hosein Khashayarpour",
    "role": "user"
  }
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"hosein@example.com","password":"SuperSecurePass123!"}'
```

---

### 2.3. Refresh Access Token
- **متد و آدرس:** `POST /api/v1/auth/refresh`
- **توضیحات:** تمدید توکن دسترسی منقضی شده با استفاده از Refresh Token معتبر
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **هدرهای ارسالی:**
```http
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "refresh_token": "b472e3a1-7c89-4d82-95f0-6c9f28d8b412"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `refresh_token` | `string` | بله | شناسه معتبر نشست کاربر (UUID) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.newly_issued_access_token..."
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"b472e3a1-7c89-4d82-95f0-6c9f28d8b412"}'
```

---

### 2.4. Logout Current Session
- **متد و آدرس:** `POST /api/v1/auth/logout`
- **توضیحات:** خروج از حساب در دستگاه فعلی و ابطال نشست و Refresh Token مربوطه
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند JWT Bearer Token کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "logged out successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/logout" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 2.5. Logout All Sessions
- **متد و آدرس:** `POST /api/v1/auth/logout-all`
- **توضیحات:** خروج همزمان از تمامی دستگاه‌ها و ابطال کامل تمامی توکن‌های کاربر
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند JWT Bearer Token کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "logged out of all sessions successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/logout-all" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

## ۳. پروفایل کاربر و مدیریت حساب (User Profile)

### 3.1. Get User Profile
- **متد و آدرس:** `GET /api/v1/profile`
- **توضیحات:** دریافت مشخصات کامل حساب کاربری احراز هویت شده جاری
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند JWT Bearer Token کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "created_at": "2026-11-20T12:00:00Z",
  "email": "hosein@example.com",
  "id": 1,
  "is_active": true,
  "name": "Hosein Khashayarpour",
  "role": "user"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/profile" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 3.2. Update User Profile
- **متد و آدرس:** `PATCH /api/v1/profile`
- **توضیحات:** ویرایش اطلاعات قابل تغییر پروفایل (مانند نام کاربر)
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند JWT Bearer Token کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "name": "Hosein Khashayarpour (Updated)"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `name` | `string` | بله | نام جدید کاربر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "email": "hosein@example.com",
  "id": 1,
  "name": "Hosein Khashayarpour (Updated)",
  "role": "user"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PATCH "http://localhost:8080/api/v1/profile" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hosein Khashayarpour (Updated)"}'
```

---

### 3.3. Change Password
- **متد و آدرس:** `POST /api/v1/profile/change-password`
- **توضیحات:** تغییر کلمه عبور کاربر و ابطال همزمان تمامی نشست‌های دیگر در سایر دستگاه‌ها
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند JWT Bearer Token کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "current_password": "SuperSecurePass123!",
  "new_password": "NewStrongPassword2026!"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `current_password` | `string` | بله | رمز عبور فعلی کاربر |
| `new_password` | `string` | بله | رمز عبور جدید (حداقل ۸ کاراکتر) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "password updated successfully; other sessions revoked"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/profile/change-password" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"current_password":"SuperSecurePass123!","new_password":"NewStrongPassword2026!"}'
```

---

## ۴. کاتالوگ عمومی فیلم و سریال (Catalog - Public)

### 4.1. List Published Movies
- **متد و آدرس:** `GET /api/v1/catalog/movies`
- **توضیحات:** دریافت لیست فیلم‌های منتشر شده با قابلیت جستجوی متنی، فیلتر ژانر، سال ساخت و صفحه‌بندی
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `limit` | `integer` | خیر | `20` | تعداد آیتم‌ها در هر صفحه |
| `offset` | `integer` | خیر | `0` | تعداد آفست برای صفحه‌بندی |
| `q` | `string` | خیر | - | کلمه کلیدی جهت جستجو در عنوان یا خلاصه |
| `year` | `integer` | خیر | - | فیلتر بر اساس سال انتشار |
| `genre_id` | `integer` | خیر | - | فیلتر بر اساس شناسه ژانر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "duration_minutes": 148,
      "genre_id": 1,
      "id": 1,
      "is_published": true,
      "poster_url": "https://storage.example.com/movies/1/poster.webp",
      "release_year": 2010,
      "slug": "inception",
      "synopsis": "A thief who steals corporate secrets...",
      "title": "Inception"
    }
  ],
  "limit": 20,
  "offset": 0,
  "total": 1
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/movies?limit=20&offset=0"
```

---

### 4.2. Get Movie By ID
- **متد و آدرس:** `GET /api/v1/catalog/movies/{id}`
- **توضیحات:** مشاهده مشخصات کامل یک فیلم با شناسه ID
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه یکتای فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "duration_minutes": 148,
  "genre_id": 1,
  "id": 1,
  "is_published": true,
  "poster_url": "https://storage.example.com/movies/1/poster.webp",
  "release_year": 2010,
  "slug": "inception",
  "synopsis": "A thief who steals corporate secrets...",
  "title": "Inception"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/movies/1"
```

---

### 4.3. List Published Series
- **متد و آدرس:** `GET /api/v1/catalog/series`
- **توضیحات:** دریافت لیست سریال‌های منتشر شده با امکان جستجو و صفحه‌بندی
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `limit` | `integer` | خیر | `20` | تعداد در صفحه |
| `offset` | `integer` | خیر | `0` | آفست صفحه‌بندی |
| `q` | `string` | خیر | - | جستجو در عنوان یا توضیحات |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "description": "A chemistry teacher diagnosed with cancer...",
      "id": 1,
      "is_published": true,
      "release_year": 2008,
      "slug": "breaking-bad",
      "title": "Breaking Bad"
    }
  ],
  "limit": 20,
  "offset": 0,
  "total": 1
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/series?limit=20&offset=0"
```

---

### 4.4. Get Series By ID
- **متد و آدرس:** `GET /api/v1/catalog/series/{id}`
- **توضیحات:** مشاهده جزییات کامل یک سریال با شناسه ID
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه سریال |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "description": "A chemistry teacher diagnosed with cancer...",
  "id": 1,
  "is_published": true,
  "release_year": 2008,
  "slug": "breaking-bad",
  "title": "Breaking Bad"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/series/1"
```

---

### 4.5. List Seasons By Series ID
- **متد و آدرس:** `GET /api/v1/catalog/series/{id}/seasons`
- **توضیحات:** دریافت لیست فصل‌های یک سریال به ترتیب شماره فصل
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه سریال |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "id": 1,
      "season_number": 1,
      "series_id": 1,
      "synopsis": "Walter White enters the drug trade",
      "title": "Season 1"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/series/1/seasons"
```

---

### 4.6. List Episodes By Season ID
- **متد و آدرس:** `GET /api/v1/catalog/seasons/{id}/episodes`
- **توضیحات:** دریافت لیست قسمت‌های یک فصل به ترتیب شماره قسمت
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه فصل |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "duration_minutes": 58,
      "episode_number": 1,
      "id": 1,
      "is_published": true,
      "season_id": 1,
      "synopsis": "Walter turns to cooking meth",
      "title": "Pilot"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/catalog/seasons/1/episodes"
```

---

## ۵. پخش ویدیو و زیرنویس‌ها (Playback & Subtitles)

### 5.1. Get Signed HLS Manifest
- **متد و آدرس:** `GET /api/v1/playback/{kind}/{id}`
- **توضیحات:** صدور URL امضا شده و دارای انقضای زمانی برای پخش استریم HLS ویدیوی فیلم یا اپیزود
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `kind` | `string` | بله | `movies` | نوع رسانه ('movies' یا 'episodes') |
| `id` | `integer` | بله | `1` | شناسه فیلم یا قسمت سریال |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "expires_in": 7200,
  "manifest_url": "https://gateway.example.com/hls/movies/1/master.m3u8?token=sig\u0026expires=1700000000"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/playback/movies/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 5.2. List Published Subtitles
- **متد و آدرس:** `GET /api/v1/media/{kind}/{id}/subtitles`
- **توضیحات:** دریافت زیرنویس‌های منتشر شده و آماده دانلود برای یک فیلم یا اپیزود
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه فیلم یا قسمت |
| `kind` | `string` | بله | `movies` | 'movies' یا 'episodes' |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "download_url": "https://storage.example.com/subtitles/1_fa.vtt",
      "format": "vtt",
      "id": 1,
      "is_generated": true,
      "label": "Persian",
      "language": "fa"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/media/movies/1/subtitles"
```

---

## ۶. تبلیغات عمومی و رهگیری رویدادها (Ads - Public)

### 6.1. Serve Active Ads
- **متد و آدرس:** `GET /api/v1/ads`
- **توضیحات:** دریافت کمپین‌های تبلیغاتی فعال متناسب با موقعیت نمایش (مانند اول ویدیو preroll)
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `placement` | `string` | خیر | `preroll` | جایگاه تبلیغ (preroll, midroll, banner) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "click_url": "https://sponsor.test/deal",
      "id": 1,
      "media_url": "https://cdn.example.com/ads/promo.mp4",
      "placement": "preroll",
      "title": "Summer Promo"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/ads?placement=preroll"
```

---

### 6.2. Track Ad Event
- **متد و آدرس:** `POST /api/v1/ads/{id}/events`
- **توضیحات:** ثبت رویداد نمایش (impression) یا کلیک (click) برای یک تبلیغ
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **هدرهای ارسالی:**
```http
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه کمپین تبلیغاتی |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "event_type": "impression",
  "metadata": "{\"device\":\"mobile\"}"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `event_type` | `string` | بله | نوع رویداد: 'impression' یا 'click' |
| `metadata` | `string` | خیر | متادیتا یا اطلاعات دستگاه به صورت رشته/JSON |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "message": "event tracked successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/ads/1/events" \
  -H "Content-Type: application/json" \
  -d '{"event_type":"impression","metadata":"{\"device\":\"mobile\"}"}'
```

---

## ۷. تعاملات: امتیازها و برگزیده‌ها (Engagement)

### 7.1. Get Movie Rating Summary
- **متد و آدرس:** `GET /api/v1/movies/{movieID}/ratings/summary`
- **توضیحات:** مشاهده میانگین امتیاز و تعداد کل آرای ثبت شده کاربران برای یک فیلم
- **سطح دسترسی:** عمومی (بدون نیاز به توکن)
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "average_score": 8.7,
  "movie_id": 1,
  "total_ratings": 1420
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/movies/1/ratings/summary"
```

---

### 7.2. Rate Movie
- **متد و آدرس:** `POST /api/v1/movies/{movieID}/ratings`
- **توضیحات:** ثبت یا تغییر امتیاز کاربر (بین ۱ تا ۱۰) برای یک فیلم
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Content-Type: application/json
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "score": 9
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `score` | `integer` | بله | نمره امتیاز بین ۱ تا ۱۰ |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "rating saved successfully",
  "movie_id": 1,
  "score": 9
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/movies/1/ratings" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"score":9}'
```

---

### 7.3. List Favorite Movies
- **متد و آدرس:** `GET /api/v1/favorites`
- **توضیحات:** دریافت لیست فیلم‌های برگزیده و نشان‌شده کاربر
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "id": 1,
      "poster_url": "https://storage.example.com/movies/1/poster.webp",
      "slug": "inception",
      "title": "Inception"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/favorites" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 7.4. Add Movie To Favorites
- **متد و آدرس:** `PUT /api/v1/favorites/{movieID}`
- **توضیحات:** افزودن یک فیلم به لیست علاقه‌مندی‌های کاربر
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "movie added to favorites"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/favorites/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 7.5. Remove Movie From Favorites
- **متد و آدرس:** `DELETE /api/v1/favorites/{movieID}`
- **توضیحات:** حذف فیلم از فهرست برگزیده‌های کاربر
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "movie removed from favorites"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/favorites/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

## ۸. تاریخچه و پیشرفت تماشا (Watch History)

### 8.1. List Movie Watch History
- **متد و آدرس:** `GET /api/v1/watch-history`
- **توضیحات:** لیست فیلم‌های در حال تماشای کاربر به همراه ثانیه آخرین موقعیت پخش
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "movie": {
        "duration_minutes": 148,
        "id": 1,
        "title": "Inception"
      },
      "movie_id": 1,
      "position_seconds": 1540,
      "updated_at": "2026-11-20T12:00:00Z"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/watch-history" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 8.2. Get Movie Progress
- **متد و آدرس:** `GET /api/v1/watch-history/{movieID}`
- **توضیحات:** دریافت ثانیه دقیق پیشرفت پخش برای یک فیلم معین
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "movie_id": 1,
  "position_seconds": 1540,
  "updated_at": "2026-11-20T12:00:00Z"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/watch-history/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 8.3. Save Movie Progress
- **متد و آدرس:** `PUT /api/v1/watch-history/{movieID}`
- **توضیحات:** ثبت یا به‌روزرسانی ثانیه تماشا برای یک فیلم (Auto-save / Resume)
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "position_seconds": 1620
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `position_seconds` | `integer` | بله | ثانیه جاری پخش فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "progress saved",
  "movie_id": 1,
  "position_seconds": 1620
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/watch-history/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"position_seconds":1620}'
```

---

### 8.4. Delete Movie Progress
- **متد و آدرس:** `DELETE /api/v1/watch-history/{movieID}`
- **توضیحات:** حذف سابقه تماشای یک فیلم خاص از تاریخچه
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `movieID` | `integer` | بله | `1` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "movie watch history deleted"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/watch-history/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 8.5. List Episode Watch History
- **متد و آدرس:** `GET /api/v1/watch-history/episodes`
- **توضیحات:** لیست قسمت‌های سریال در حال تماشا با ثانیه موقعیت
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "episode": {
        "duration_minutes": 58,
        "id": 1,
        "title": "Pilot"
      },
      "episode_id": 1,
      "position_seconds": 920,
      "updated_at": "2026-11-20T12:00:00Z"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/watch-history/episodes" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 8.6. Save Episode Progress
- **متد و آدرس:** `PUT /api/v1/watch-history/episodes/{episodeID}`
- **توضیحات:** ثبت یا به‌روزرسانی ثانیه تماشای یک قسمت سریال
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `episodeID` | `integer` | بله | `1` | شناسه قسمت سریال |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "position_seconds": 980
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `position_seconds` | `integer` | بله | ثانیه جاری تماشای اپیزود |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "episode_id": 1,
  "message": "progress saved",
  "position_seconds": 980
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/watch-history/episodes/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"position_seconds":980}'
```

---

### 8.7. Delete Episode Progress
- **متد و آدرس:** `DELETE /api/v1/watch-history/episodes/{episodeID}`
- **توضیحات:** حذف سابقه تماشای یک قسمت سریال
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `episodeID` | `integer` | بله | `1` | شناسه اپیزود |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "episode watch history deleted"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/watch-history/episodes/1" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

## ۹. اعلان‌های کاربری (User Notifications)

### 9.1. List User Notifications
- **متد و آدرس:** `GET /api/v1/notifications`
- **توضیحات:** دریافت لیست اعلان‌های کاربر با وضعیت خوانده شده/نشده و صفحه‌بندی
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `limit` | `integer` | خیر | `20` | تعداد در صفحه |
| `offset` | `integer` | خیر | `0` | آفست صفحه‌بندی |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "created_at": "2026-11-20T12:00:00Z",
      "id": 1,
      "is_read": false,
      "message": "Welcome to our streaming platform.",
      "title": "Welcome!"
    }
  ],
  "total": 1
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/notifications?limit=20&offset=0" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 9.2. Mark Notification As Read
- **متد و آدرس:** `PATCH /api/v1/notifications/{id}/read`
- **توضیحات:** علامت‌گذاری یک اعلان مشخص به عنوان خوانده شده
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه اعلان |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "notification marked as read"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PATCH "http://localhost:8080/api/v1/notifications/1/read" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

### 9.3. Mark All Notifications As Read
- **متد و آدرس:** `POST /api/v1/notifications/read-all`
- **توضیحات:** علامت‌گذاری همزمان تمام اعلان‌های خوانده نشده به عنوان خوانده شده
- **سطح دسترسی:** الزامی (`bearer_user`) — نیازمند توکن کاربر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "all notifications marked as read"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/notifications/read-all" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

---

## ۱۰. داشبورد و آمار مدیریتی (Admin Dashboard & Analytics)

### 10.1. Admin Dashboard Counts
- **متد و آدرس:** `GET /api/v1/admin/dashboard`
- **توضیحات:** دریافت آمار کلی پنل مدیریت: تعداد کل کاربران، فیلم‌ها، سریال‌ها، قسمت‌ها و تبلیغات
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر (Admin/Super Admin)
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "active_ad_campaigns": 8,
  "total_episodes": 480,
  "total_movies": 320,
  "total_series": 45,
  "total_users": 1500
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/dashboard" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

### 10.2. Admin System Analytics
- **متد و آدرس:** `GET /api/v1/admin/analytics`
- **توضیحات:** تحلیل‌های آماری پلتفرم شامل تعداد کل تماشاها، فیلم‌های پرمخاطب بر اساس بازه روز
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `window_days` | `integer` | خیر | `30` | بازه زمانی تحلیل (۷، ۳۰، یا ۹۰ روز) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "top_movies": [
    {
      "movie_id": 1,
      "title": "Inception",
      "views": 12000
    }
  ],
  "total_views": 45000,
  "window_days": 30
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/analytics?window_days=30" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

### 10.3. List Admin Audit Logs
- **متد و آدرس:** `GET /api/v1/admin/audit-logs`
- **توضیحات:** مشاهده لاگ‌های امنیتی تغییرات اعمال شده توسط ادمین‌ها (Audit Trail)
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `action` | `string` | خیر | - | فیلتر بر اساس نام عملیات |
| `limit` | `integer` | خیر | `50` | تعداد لاگ در صفحه |
| `offset` | `integer` | خیر | `0` | آفست |
| `user_id` | `integer` | خیر | - | فیلتر بر اساس شناسه مدیر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "action": "CREATE_MOVIE",
      "created_at": "2026-11-20T12:00:00Z",
      "entity_id": 1,
      "entity_type": "movie",
      "id": 1,
      "ip_address": "127.0.0.1",
      "user_id": 1
    }
  ],
  "total": 1
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/audit-logs?limit=50&offset=0" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۱. مدیریت کاربران (Admin Users)

### 11.1. List All Users
- **متد و آدرس:** `GET /api/v1/admin/users`
- **توضیحات:** مشاهده لیست کامل کاربران به همراه فیلتر وضعیت حساب و جستجو
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `is_active` | `boolean` | خیر | - | فیلتر کاربران فعال/مسدود شده |
| `limit` | `integer` | خیر | `20` | تعداد در صفحه |
| `offset` | `integer` | خیر | `0` | آفست |
| `q` | `string` | خیر | - | جستجو در نام یا ایمیل |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "email": "hosein@example.com",
      "id": 1,
      "is_active": true,
      "name": "Hosein Khashayarpour",
      "role": "super_admin"
    }
  ],
  "total": 1
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/users?limit=20&offset=0" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

### 11.2. Set User Active Status
- **متد و آدرس:** `PATCH /api/v1/admin/users/{id}/status`
- **توضیحات:** فعال‌سازی یا مسدودسازی حساب کاربری توسط مدیر
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه کاربر هدف |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "is_active": false
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `is_active` | `boolean` | بله | وضعیت فعال (true) یا غیرفعال (false) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 1,
  "is_active": false,
  "message": "user status updated"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PATCH "http://localhost:8080/api/v1/admin/users/1/status" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"is_active":false}'
```

---

### 11.3. Set User Role
- **متد و آدرس:** `PATCH /api/v1/admin/users/{id}/role`
- **توضیحات:** تغییر نقش کاربر به 'user'، 'admin' یا 'super_admin' (فقط توسط سوپر ادمین)
- **سطح دسترسی:** الزامی (`bearer_super_admin`) — نیازمند نقش Super Admin
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_SUPER_ADMIN_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه کاربر هدف |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "role": "admin"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `role` | `string` | بله | نقش جدید: 'user', 'admin', 'super_admin' |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 1,
  "message": "user role updated",
  "role": "admin"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PATCH "http://localhost:8080/api/v1/admin/users/1/role" \
  -H "Authorization: Bearer <YOUR_SUPER_ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}'
```

---

## ۱۲. مدیریت ژانرها (Admin Genres)

### 12.1. Create Genre
- **متد و آدرس:** `POST /api/v1/admin/genres`
- **توضیحات:** ایجاد یک دسته‌بندی یا ژانر جدید برای فیلم‌ها و سریال‌ها
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "name": "Sci-Fi",
  "slug": "sci-fi"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `name` | `string` | بله | نام ژانر |
| `slug` | `string` | بله | اسلاگ انگلیسی یکتا |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "id": 1,
  "name": "Sci-Fi",
  "slug": "sci-fi"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/genres" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Sci-Fi","slug":"sci-fi"}'
```

---

### 12.2. Update Genre
- **متد و آدرس:** `PUT /api/v1/admin/genres/{id}`
- **توضیحات:** ویرایش نام یا اسلاگ یک ژانر موجود
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Content-Type: application/json
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه ژانر |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "name": "Science Fiction",
  "slug": "science-fiction"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `name` | `string` | بله | نام جدید ژانر |
| `slug` | `string` | بله | اسلاگ جدید ژانر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 1,
  "name": "Science Fiction",
  "slug": "science-fiction"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/genres/1" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Science Fiction","slug":"science-fiction"}'
```

---

### 12.3. Delete Genre
- **متد و آدرس:** `DELETE /api/v1/admin/genres/{id}`
- **توضیحات:** حذف ژانر (در صورتی که فیلم وابسته داشته باشد رد می‌شود)
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه ژانر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "genre deleted successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/genres/1" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۳. مدیریت فیلم‌ها (Admin Movies)

### 13.1. Create Movie
- **متد و آدرس:** `POST /api/v1/admin/movies`
- **توضیحات:** ایجاد یک فیلم جدید در کاتالوگ با عنوان، اسلاگ، سال و مدت زمان
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "duration_minutes": 169,
  "genre_id": 1,
  "is_published": true,
  "release_year": 2014,
  "slug": "interstellar",
  "synopsis": "A team of explorers travel through a wormhole in space...",
  "title": "Interstellar"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `synopsis` | `string` | خیر | خلاصه داستان |
| `title` | `string` | بله | عنوان فیلم |
| `duration_minutes` | `integer` | بله | مدت زمان به دقیقه |
| `genre_id` | `integer` | بله | شناسه ژانر اصلی |
| `is_published` | `boolean` | بله | منتشر شده برای کاربران یا پیش‌نویس |
| `release_year` | `integer` | بله | سال ساخت |
| `slug` | `string` | بله | اسلاگ یکتا |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "duration_minutes": 169,
  "genre_id": 1,
  "id": 2,
  "is_published": true,
  "release_year": 2014,
  "slug": "interstellar",
  "title": "Interstellar"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/movies" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Interstellar","slug":"interstellar","synopsis":"A team of explorers travel through a wormhole...","release_year":2014,"duration_minutes":169,"genre_id":1,"is_published":true}'
```

---

### 13.2. Update Movie
- **متد و آدرس:** `PUT /api/v1/admin/movies/{id}`
- **توضیحات:** ویرایش مشخصات فیلم موجود با شناسه ID
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه فیلم |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "duration_minutes": 169,
  "genre_id": 1,
  "is_published": true,
  "release_year": 2014,
  "slug": "interstellar-imax",
  "synopsis": "Updated synopsis...",
  "title": "Interstellar (IMAX)"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `is_published` | `boolean` | بله | وضعیت انتشار |
| `release_year` | `integer` | بله | سال ساخت |
| `slug` | `string` | بله | اسلاگ |
| `synopsis` | `string` | خیر | خلاصه داستان |
| `title` | `string` | بله | عنوان فیلم |
| `duration_minutes` | `integer` | بله | مدت زمان به دقیقه |
| `genre_id` | `integer` | بله | شناسه ژانر |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 2,
  "slug": "interstellar-imax",
  "title": "Interstellar (IMAX)"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/movies/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Interstellar (IMAX)","slug":"interstellar-imax","release_year":2014,"duration_minutes":169,"genre_id":1,"is_published":true}'
```

---

### 13.3. Delete Movie
- **متد و آدرس:** `DELETE /api/v1/admin/movies/{id}`
- **توضیحات:** حذف فیلم و وابستگی‌های آن از سامانه
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه فیلم |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "movie deleted successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/movies/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۴. مدیریت سریال‌ها (Admin Series)

### 14.1. Create Series
- **متد و آدرس:** `POST /api/v1/admin/series`
- **توضیحات:** ایجاد یک سریال جدید در کاتالوگ
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "description": "When a young boy vanishes, a small town uncovers a mystery...",
  "is_published": true,
  "release_year": 2016,
  "slug": "stranger-things",
  "title": "Stranger Things"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `slug` | `string` | بله | اسلاگ یکتا |
| `title` | `string` | بله | عنوان سریال |
| `description` | `string` | خیر | توضیحات سریال |
| `is_published` | `boolean` | بله | وضعیت انتشار |
| `release_year` | `integer` | بله | سال آغاز انتشار |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "id": 2,
  "is_published": true,
  "release_year": 2016,
  "slug": "stranger-things",
  "title": "Stranger Things"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/series" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Stranger Things","slug":"stranger-things","release_year":2016,"is_published":true}'
```

---

### 14.2. Update Series
- **متد و آدرس:** `PUT /api/v1/admin/series/{id}`
- **توضیحات:** ویرایش مشخصات سریال
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه سریال |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "description": "Updated description...",
  "is_published": true,
  "release_year": 2016,
  "slug": "stranger-things",
  "title": "Stranger Things Updated"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `description` | `string` | خیر | توضیحات |
| `is_published` | `boolean` | بله | وضعیت انتشار |
| `release_year` | `integer` | بله | سال انتشار |
| `slug` | `string` | بله | اسلاگ |
| `title` | `string` | بله | عنوان جدید |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 2,
  "title": "Stranger Things Updated"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/series/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Stranger Things Updated","slug":"stranger-things","release_year":2016,"is_published":true}'
```

---

### 14.3. Delete Series
- **متد و آدرس:** `DELETE /api/v1/admin/series/{id}`
- **توضیحات:** حذف کامل یک سریال و فصل‌ها و قسمت‌های وابسته به آن
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه سریال |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "series deleted successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/series/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۵. مدیریت فصل‌ها (Admin Seasons)

### 15.1. Create Season
- **متد و آدرس:** `POST /api/v1/admin/seasons`
- **توضیحات:** افزودن یک فصل به یک سریال موجود
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "season_number": 2,
  "series_id": 1,
  "synopsis": "The story continues with higher stakes",
  "title": "Season 2"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `season_number` | `integer` | بله | شماره فصل |
| `series_id` | `integer` | بله | شناسه سریال والد |
| `synopsis` | `string` | خیر | خلاصه فصل |
| `title` | `string` | خیر | عنوان اختیاری فصل |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "id": 2,
  "season_number": 2,
  "series_id": 1,
  "title": "Season 2"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/seasons" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"series_id":1,"season_number":2,"title":"Season 2","synopsis":"The story continues"}'
```

---

### 15.2. Update Season
- **متد و آدرس:** `PUT /api/v1/admin/seasons/{id}`
- **توضیحات:** ویرایش عنوان و شماره فصل
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه فصل |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "season_number": 2,
  "synopsis": "Updated synopsis",
  "title": "Season 2 (Special)"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `season_number` | `integer` | بله | شماره فصل |
| `synopsis` | `string` | خیر | خلاصه |
| `title` | `string` | خیر | عنوان |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 2,
  "season_number": 2,
  "title": "Season 2 (Special)"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/seasons/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"season_number":2,"title":"Season 2 (Special)"}'
```

---

### 15.3. Delete Season
- **متد و آدرس:** `DELETE /api/v1/admin/seasons/{id}`
- **توضیحات:** حذف یک فصل و قسمت‌های وابسته به آن
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه فصل |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "season deleted successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/seasons/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۶. مدیریت قسمت‌ها (Admin Episodes)

### 16.1. Create Episode
- **متد و آدرس:** `POST /api/v1/admin/episodes`
- **توضیحات:** افزودن یک قسمت جدید به یک فصل
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "duration_minutes": 48,
  "episode_number": 2,
  "is_published": true,
  "season_id": 1,
  "synopsis": "Walt and Jesse attempt to dispose of two bodies",
  "title": "Cat's in the Bag..."
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `synopsis` | `string` | خیر | خلاصه اپیزود |
| `title` | `string` | بله | عنوان اپیزود |
| `duration_minutes` | `integer` | بله | مدت زمان به دقیقه |
| `episode_number` | `integer` | بله | شماره قسمت |
| `is_published` | `boolean` | بله | وضعیت انتشار |
| `season_id` | `integer` | بله | شناسه فصل والد |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "duration_minutes": 48,
  "episode_number": 2,
  "id": 2,
  "is_published": true,
  "season_id": 1,
  "title": "Cat's in the Bag..."
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/episodes" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"season_id":1,"episode_number":2,"title":"Cat'\''s in the Bag...","duration_minutes":48,"is_published":true}'
```

---

### 16.2. Update Episode
- **متد و آدرس:** `PUT /api/v1/admin/episodes/{id}`
- **توضیحات:** ویرایش اطلاعات یک اپیزود
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه اپیزود |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "duration_minutes": 48,
  "episode_number": 2,
  "is_published": true,
  "synopsis": "Updated episode synopsis",
  "title": "Cat's in the Bag (Remastered)"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `episode_number` | `integer` | بله | شماره قسمت |
| `is_published` | `boolean` | بله | وضعیت انتشار |
| `synopsis` | `string` | خیر | خلاصه |
| `title` | `string` | بله | عنوان |
| `duration_minutes` | `integer` | بله | مدت زمان |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 2,
  "title": "Cat's in the Bag (Remastered)"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/admin/episodes/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"episode_number":2,"title":"Cat'\''s in the Bag (Remastered)","duration_minutes":48,"is_published":true}'
```

---

### 16.3. Delete Episode
- **متد و آدرس:** `DELETE /api/v1/admin/episodes/{id}`
- **توضیحات:** حذف یک قسمت سریال
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه اپیزود |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "episode deleted successfully"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X DELETE "http://localhost:8080/api/v1/admin/episodes/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

## ۱۷. آپلود مستقیم به S3 و است‌های رسانه‌ای (Admin Media & S3)

### 17.1. Initiate Presigned S3 Upload
- **متد و آدرس:** `POST /api/v1/admin/media/uploads`
- **توضیحات:** ایجاد لینک امضا شده موقت S3 (Presigned PUT) برای آپلود مستقیم ویدیوی باکیفیت یا پوستر از فرانت‌اند
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "content_type": "video/mp4",
  "kind": "video",
  "size_bytes": 1048576000,
  "target_id": 1,
  "target_type": "movie"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `size_bytes` | `integer` | بله | حجم فایل به بایت |
| `target_id` | `integer` | بله | شناسه فیلم یا اپیزود هدف |
| `target_type` | `string` | بله | 'movie' یا 'episode' |
| `content_type` | `string` | بله | MIME type فایل (مثلاً video/mp4) |
| `kind` | `string` | بله | نوع است: 'video', 'trailer', 'poster' |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "expires_in": 3600,
  "storage_key": "movies/1/video/source.mp4",
  "upload_url": "https://s3.example.com/bucket/movies/1/video/source.mp4?X-Amz-Signature=..."
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/media/uploads" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"target_type":"movie","target_id":1,"kind":"video","content_type":"video/mp4","size_bytes":1048576000}'
```

---

### 17.2. Complete And Verify S3 Upload
- **متد و آدرس:** `POST /api/v1/admin/media/uploads/complete`
- **توضیحات:** راستی‌آزمایی وجود فایل در باکت S3 و ثبت در دیتابیس و قرارگیری خودکار در صف ترنسکود HLS
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Content-Type: application/json
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "kind": "video",
  "storage_key": "movies/1/video/source.mp4",
  "target_id": 1,
  "target_type": "movie"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `target_id` | `integer` | بله | شناسه هدف |
| `target_type` | `string` | بله | 'movie' یا 'episode' |
| `kind` | `string` | بله | 'video', 'trailer', 'poster' |
| `storage_key` | `string` | بله | کلید آبجکت استوریج ذخیره‌شده |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "asset": {
    "id": 1,
    "is_active": false,
    "kind": "video",
    "mime_type": "video/mp4",
    "movie_id": 1,
    "size_bytes": 1048576000,
    "storage_key": "movies/1/video/source.mp4"
  },
  "message": "source upload verified; asset remains inactive until HLS processing completes",
  "status": "uploaded_pending_processing"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/media/uploads/complete" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"target_type":"movie","target_id":1,"kind":"video","storage_key":"movies/1/video/source.mp4"}'
```

---

### 17.3. Register Asset Manually
- **متد و آدرس:** `POST /api/v1/admin/media/assets`
- **توضیحات:** ثبت دستی یک فایل رسانه‌ای آماده در سیستم (مثلاً پوستر در CDN یا لینک مستقیم)
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "kind": "poster",
  "mime_type": "image/webp",
  "movie_id": 1,
  "playback_url": "https://storage.example.com/movies/1/poster.webp",
  "size_bytes": 245000,
  "storage_key": "movies/1/poster.webp"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `storage_key` | `string` | بله | کلید ذخیره‌سازی در باکت |
| `episode_id` | `integer` | خیر | شناسه قسمت (اگر برای قسمت است) |
| `kind` | `string` | بله | 'video', 'trailer', 'poster' |
| `mime_type` | `string` | بله | MIME type |
| `movie_id` | `integer` | خیر | شناسه فیلم (اگر برای فیلم است) |
| `playback_url` | `string` | خیر | URL مستقیم در صورت وجود |
| `size_bytes` | `integer` | بله | حجم فایل به بایت |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "id": 2,
  "kind": "poster",
  "movie_id": 1,
  "playback_url": "https://storage.example.com/movies/1/poster.webp",
  "storage_key": "movies/1/poster.webp"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/media/assets" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"movie_id":1,"kind":"poster","storage_key":"movies/1/poster.webp","mime_type":"image/webp","size_bytes":245000,"playback_url":"https://storage.example.com/movies/1/poster.webp"}'
```

---

### 17.4. Queue Subtitle Generation Job
- **متد و آدرس:** `POST /api/v1/admin/media/subtitle-jobs`
- **توضیحات:** ارسال کار استخراج صوت و تولید زیرنویس هوش مصنوعی به صف پردازش ورکر
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "source_asset_id": 1,
  "source_language": "en",
  "target_language": "fa"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `source_asset_id` | `integer` | بله | شناسه فایل ویدیوی مرجع |
| `source_language` | `string` | بله | زبان صوت فیلم (مثلاً 'en') |
| `target_language` | `string` | بله | زبان مقصد زیرنویس (مثلاً 'fa') |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "job_id": 1,
  "source_asset_id": 1,
  "status": "queued",
  "target_language": "fa"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/admin/media/subtitle-jobs" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"source_asset_id":1,"source_language":"en","target_language":"fa"}'
```

---

### 17.5. Publish Subtitle
- **متد و آدرس:** `PATCH /api/v1/admin/media/subtitles/{id}/publish`
- **توضیحات:** تغییر وضعیت انتشار زیرنویس (قابل دسترسی عمومی برای بینندگان یا مخفی)
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه زیرنویس |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "is_published": true
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `is_published` | `boolean` | بله | وضعیت انتشار (true/false) |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 1,
  "is_published": true,
  "message": "subtitle publication status updated"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PATCH "http://localhost:8080/api/v1/admin/media/subtitles/1/publish" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"is_published":true}'
```

---

## ۱۸. مدیریت کمپین‌های تبلیغاتی (Admin Ads)

### 18.1. List All Ad Campaigns
- **متد و آدرس:** `GET /api/v1/ads`
- **توضیحات:** لیست کامل کمپین‌های تبلیغاتی شامل کمپین‌های فعال، منقضی شده یا پیش‌نویس
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **پارامترهای کوئری (Query Parameters):**

| پارامتر | نوع | الزامی | پیش‌فرض | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `placement` | `string` | خیر | - | فیلتر بر اساس جایگاه |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "items": [
    {
      "click_url": "https://brand.test",
      "end_at": "2026-08-31T23:59:59Z",
      "id": 1,
      "is_active": true,
      "media_url": "https://cdn.example.com/ad.mp4",
      "placement": "preroll",
      "start_at": "2026-06-01T00:00:00Z",
      "title": "Summer Campaign"
    }
  ]
}
```
- **دستور اجرا (cURL):**
```bash
curl -X GET "http://localhost:8080/api/v1/ads" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>"
```

---

### 18.2. Create Ad Campaign
- **متد و آدرس:** `POST /api/v1/ads`
- **توضیحات:** ایجاد یک کمپین تبلیغاتی جدید با زمان‌بندی و لینک فرود
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "click_url": "https://sponsor.test/spring",
  "end_at": "2026-04-01T00:00:00Z",
  "is_active": true,
  "media_url": "https://cdn.example.com/spring.mp4",
  "placement": "preroll",
  "start_at": "2026-03-01T00:00:00Z",
  "title": "Spring Promo"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `is_active` | `boolean` | بله | فعال بودن کمپین |
| `media_url` | `string` | بله | آدرس فایل ویدیو یا بنر |
| `placement` | `string` | بله | 'preroll', 'midroll', 'banner' |
| `start_at` | `string` | بله | تاریخ شروع |
| `title` | `string` | بله | عنوان کمپین |
| `click_url` | `string` | بله | آدرس سایت هدف هنگام کلیک |
| `end_at` | `string` | بله | تاریخ پایان |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "id": 2,
  "is_active": true,
  "placement": "preroll",
  "title": "Spring Promo"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/ads" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Spring Promo","placement":"preroll","media_url":"https://cdn.example.com/spring.mp4","click_url":"https://sponsor.test/spring","is_active":true,"start_at":"2026-03-01T00:00:00Z","end_at":"2026-04-01T00:00:00Z"}'
```

---

### 18.3. Update Ad Campaign
- **متد و آدرس:** `PUT /api/v1/ads/{id}`
- **توضیحات:** ویرایش مشخصات یا تمدید زمان یک کمپین تبلیغاتی
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `2` | شناسه کمپین |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "click_url": "https://sponsor.test/spring",
  "end_at": "2026-04-15T00:00:00Z",
  "is_active": true,
  "media_url": "https://cdn.example.com/spring.mp4",
  "placement": "preroll",
  "start_at": "2026-03-01T00:00:00Z",
  "title": "Spring Promo Extended"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `click_url` | `string` | بله | لینک فرود |
| `end_at` | `string` | بله | تاریخ پایان |
| `is_active` | `boolean` | بله | وضعیت |
| `media_url` | `string` | بله | فایل مدیا |
| `placement` | `string` | بله | جایگاه |
| `start_at` | `string` | بله | تاریخ شروع |
| `title` | `string` | بله | عنوان |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "id": 2,
  "title": "Spring Promo Extended"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X PUT "http://localhost:8080/api/v1/ads/2" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Spring Promo Extended","placement":"preroll","media_url":"https://cdn.example.com/spring.mp4","click_url":"https://sponsor.test/spring","is_active":true,"start_at":"2026-03-01T00:00:00Z","end_at":"2026-04-15T00:00:00Z"}'
```

---

## ۱۹. ارسال اعلان همگانی (Admin Broadcast Notifications)

### 19.1. Broadcast Notification
- **متد و آدرس:** `POST /api/v1/notifications/broadcast`
- **توضیحات:** ارسال اعلان همگانی (Broadcast) به تمام کاربران فعال سامانه
- **سطح دسترسی:** الزامی (`bearer_admin`) — نیازمند توکن مدیر
- **هدرهای ارسالی:**
```http
Content-Type: application/json
Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>
```
- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "message": "Scheduled maintenance tonight at 02:00 AM UTC.",
  "title": "System Maintenance"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `message` | `string` | بله | متن کامل پیام |
| `title` | `string` | بله | عنوان اعلان همگانی |

- **نمونه پاسخ سرور (Response 201):**
```json
{
  "message": "broadcast notification dispatched to all active users"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/api/v1/notifications/broadcast" \
  -H "Authorization: Bearer <YOUR_ADMIN_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"System Maintenance","message":"Scheduled maintenance tonight at 02:00 AM UTC."}'
```

---

## ۲۰. ورکر داخلی پردازش زیرنویس (Internal Subtitle Worker)

### 20.1. Claim Next Subtitle Job
- **متد و آدرس:** `POST /internal/v1/worker/subtitle-jobs/claim`
- **توضیحات:** تحویل گرفتن و قفل کردن جاب بعدی تولید زیرنویس از صف پردازش توسط ورکر
- **سطح دسترسی:** الزامی (`worker_token`) — نیازمند SUBTITLE_WORKER_TOKEN در هدر Bearer
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_WORKER_TOKEN>
```
- **نمونه پاسخ سرور (Response 200):**
```json
{
  "job": {
    "id": 1,
    "source_asset_id": 1,
    "source_language": "en",
    "status": "processing",
    "target_language": "fa"
  },
  "source_video_url": "https://s3.example.com/download/source.mp4?..."
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/internal/v1/worker/subtitle-jobs/claim" \
  -H "Authorization: Bearer <YOUR_WORKER_TOKEN>"
```

---

### 20.2. Complete Subtitle Job
- **متد و آدرس:** `POST /internal/v1/worker/subtitle-jobs/{id}/complete`
- **توضیحات:** ثبت اتمام موفق پردازش جاب همراه با لینک فایل زیرنویس تولید شده یا ثبت خطا
- **سطح دسترسی:** الزامی (`worker_token`) — نیازمند SUBTITLE_WORKER_TOKEN
- **هدرهای ارسالی:**
```http
Authorization: Bearer <YOUR_WORKER_TOKEN>
Content-Type: application/json
```
- **پارامترهای آدرس (Path Parameters):**

| پارامتر | نوع | الزامی | مقدار نمونه | توضیحات |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `integer` | بله | `1` | شناسه جاب |

- **نمونه بادی ارسالی (Request Body JSON):**
```json
{
  "download_url": "https://storage.example.com/subtitles/job1_fa.vtt",
  "error_message": "",
  "status": "completed"
}
```
- **توضیحات فیلدهای ورودی:**

| فیلد | نوع | الزامی | توضیحات |
| :--- | :--- | :--- | :--- |
| `download_url` | `string` | خیر | آدرس دانلود فایل نهایی زیرنویس در صورت موفقیت |
| `error_message` | `string` | خیر | پیام خطا در صورت شکست |
| `status` | `string` | بله | 'completed' یا 'failed' |

- **نمونه پاسخ سرور (Response 200):**
```json
{
  "message": "job completed and subtitle registered"
}
```
- **دستور اجرا (cURL):**
```bash
curl -X POST "http://localhost:8080/internal/v1/worker/subtitle-jobs/1/complete" \
  -H "Authorization: Bearer <YOUR_WORKER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status":"completed","download_url":"https://storage.example.com/subtitles/job1_fa.vtt","error_message":""}'
```

---

