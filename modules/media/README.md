# Media & Storage Module

ماژول آپلود مستقیم به Object Storage (S3/MinIO/Cloudflare R2)، ایجاد لینک‌های امن استریم HLS و مدیریت زیرنویس‌ها.

## قابلیت‌ها
- **آپلود مستقیم کلاینت به فضای ابری (Direct Presigned PUT):** تولید URL امن تک‌بارمصرف با انقضای ۱۰ دقیقه‌ای جهت جلوگیری از درگیر شدن پهنای باند سرور اصلی.
- **تأییدیه نهایی آپلود (Upload Finalize):** راستی‌آزمایی وجود فایل و تطبیق سایز و فرمت (MIME Type) در استوریج قبل از ایجاد رکورد رسانه.
- **امضای اختصاصی استریم (Signed Playback URLs):** تولید آدرس‌های کوتاه‌مدت مانیفست HLS (`master.m3u8`) با امضای HMAC-SHA256 و اعتبارسنجی در لایه CDN.
- **سیستم مدیریت زیرنویس و ترجمه هوش مصنوعی:** ثبت و مدیریت فایل‌های VTT/SRT و ایجاد جاب‌های تولید زیرنویس فارسی برای ورکر هوش مصنوعی.

## ساختار مسیر ذخیره‌سازی S3
```
movies/{movie_id}/{kind}/{year}/{month}/{random_hex}.{ext}
episodes/{episode_id}/{kind}/{year}/{month}/{random_hex}.{ext}
```

## اندپوینت‌ها (API Endpoints)
| متد | مسیر | دسترسی | توضیحات |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/playback/{kind}/{id}` | با توکن | دریافت مانیفست امضاشده HLS با لینک انقضادار |
| `GET` | `/api/v1/media/{kind}/{id}/subtitles` | عمومی | لیست زیرنویس‌های منتشرشده برای فیلم یا اپیزود |
| `POST` | `/api/v1/admin/media/uploads` | مدیر | صدور لینک موقت آپلود مستقیم فایل ویدیویی یا پوستر به S3 |
| `POST` | `/api/v1/admin/media/uploads/complete` | مدیر | اعتبارسنجی نهایی و ثبت فایل آپلودشده در دیتابیس |
| `POST` | `/api/v1/admin/media/assets` | مدیر | ثبت دستی است ویدیویی یا پوستر |
| `POST` | `/api/v1/admin/media/subtitle-jobs` | مدیر | ثبت جاب جدید در صف هوش مصنوعی جهت استخراج و ترجمه زیرنویس |
| `PATCH` | `/api/v1/admin/media/subtitles/{id}/publish` | مدیر | تغییر وضعیت انتشار زیرنویس |
| `POST` | `/internal/v1/worker/subtitle-jobs/claim` | ورکر داخلی | قفل و تحویل گرفتن جاب زیرنویس توسط ورکر پردازشی |
| `POST` | `/internal/v1/worker/subtitle-jobs/{id}/complete` | ورکر داخلی | ثبت خروجی زیرنویس یا گزارش خطای ورکر |

## متغیرهای محیطی مرتبط
- `S3_ENDPOINT`: آدرس سرور S3 (مثلاً MinIO، ArvanCloud یا AWS).
- `S3_REGION`, `S3_BUCKET`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`.
- `PLAYBACK_CDN_BASE_URL`: آدرس پایه دامنه CDN پخش مدیا (HTTPS).
- `PLAYBACK_SIGNING_KEY`: کلید اختصاصی امضای لینک‌های HLS (حداقل ۳۲ بایت).
- `SUBTITLE_WORKER_TOKEN`: توکن اختصاصی ورکر ترجمه زیرنویس.
