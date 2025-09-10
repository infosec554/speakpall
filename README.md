2025-09-14 (Yakshanba) — Project skeleton & infra

Vazifalar:

Go mod init, folder struktura yaratish.

config (siz yozgan), logger (zap) ni ulash.

Dockerfile, docker-compose.yml, Makefile ishlashini tekshirish.

Health endpoint: GET /health (200, build info).

Natija: make up-local bilan app 8080’da turadi, /health OK.

Tekshiruv: container healthcheck green.

2025-09-15 (Dushanba) — DB migratsiyalar (Postgres)

Vazifalar:

Oldingi bergan SQL sxemani migrations/001_init.sql qilib qo‘yish (triggersiz).

pgx repo layer skeleton: db.Connect(), ping, ctx timeout.

Indekslar borligini tekshirish.

Natija: make up-local → Postgres ko‘tariladi, migration qo‘llanadi.

Tekshiruv: sessions, messages, users jadvallari bor.

2025-09-16 (Seshanba) — Auth: Email OTP (verify)

Vazifalar:

auth_email_tokens repo: create, mark used, expire check.

POST /auth/email/request (email qabul qiladi) → OTP yaratadi, SMTP yuboradi.

POST /auth/email/verify (email+code) → users + auth_identities(provider='email', email_verified=true) yaratadi yoki tapadi, JWT qaytaradi.

Natija: Email verify flow ishlaydi (dev’da Gmail app password).

Tekshiruv: token expiry/used case lar.

2025-09-17 (Chorshanba) — Auth: Google OAuth

Vazifalar:

GET /auth/google/login (redirect URL), GET /auth/google/callback → token olib, auth_identities(provider='google') yaratish.

JWT generatsiya.

Middleware: AuthRequired (JWT parse).

Natija: Google bilan kirish ishlaydi.

Tekshiruv: bir user’da 2 provider bo‘lsa ham OK.

2025-09-18 (Payshanba) — Profiles API

Vazifalar:

GET /me (users+profiles join)

PATCH /me (display_name, about, level, timezone, interests_txt)

Avatar upload → S3 presigned URL endpoint (ixtiyoriy).

profiles.updated_at ni app layer’da yangilash.

Natija: Profil sozlanadi.

Tekshiruv: Update → GET bilan tekshir.

2025-09-19 (Juma) — Redis & presence

Vazifalar:

Redis client, Ping().

Presence keys: online:{userID} EX=60s heart-beat (clientdan 30s da bir POST /presence/ping).

Match queue set/list: queue:{level} yoki queue:{interest} modelini tanlash (MVP: queue:{level}).

Natija: Online/presence ko‘rinadi.

Tekshiruv: Redis’da keys paydo bo‘ladi, expire ishlaydi.

2025-09-20 (Shanba) — Matchmaking (MVP, server-pull)

Vazifalar:

POST /match/attempt {desired_level} → Redis queue ga qo‘yish, match_attempts ga audit yozish.

Orqa fon service (goroutine/cron) queue’dan juftlab sessions yaratadi.

GET /match/status → matched bo‘lsa session_id qaytarish.

Natija: Match topiladi, session yaratiladi.

Tekshiruv: 2 ta client bilan sinov.

2025-09-21 (Yakshanba) — Sessions API

Vazifalar:

GET /sessions?cursor=&limit= (my history)

GET /sessions/{id} (meta) — faqat ishtirokchilarga.

POST /sessions/{id}/end {rating, notes} → duration va ended_at to‘ldirish.

Natija: Sessiya boshqaruvi.

Tekshiruv: Oxirida streak triggerini app’da ishlatish (keyingi kun).

2025-09-22 (Dushanba) — Messages API (text)

Vazifalar:

POST /sessions/{id}/messages {type='text', text}

GET /sessions/{id}/messages?cursor=&limit= (DESC paginate + cursor)

pg_trgm qidiruv: GET /messages/search?q=... (MVP ixtiyoriy).

Natija: Text chat ishlaydi.

Tekshiruv: 2 user orasida yuborish/olish.

2025-09-23 (Seshanba) — Messages (voice link) & S3

Vazifalar:

POST /uploads/presign → presigned S3 PUT URL qaytarish.

POST /sessions/{id}/messages {type='voice', audio_url} (URL clientdan keladi).

UI/Client test: avval yukla → keyin message create.

Natija: Voice xabarlar link bilan ishlaydi.

Tekshiruv: Audio URL saqlanadi, GET bilan ko‘rinadi.

2025-09-24 (Chorshanba) — Block/Report

Vazifalar:

POST /users/{id}/block, DELETE /users/{id}/block

POST /reports {target_user_id, reason, note?}

Matching/service layer’da block chek (app layer).

Natija: Xavfsizlik nazorati.

Tekshiruv: Block qilingan user match bo‘lmasin.

2025-09-25 (Payshanba) — Streaks & Badges (MVP)

Vazifalar:

Session tugaganda (/end) app darajasida streak update.

Badge logika: FIRST_CHAT, SEVEN_DAYS.

GET /me/streaks, GET /me/badges.

Natija: Gamification ko‘rinadi.

Tekshiruv: Mock holatda streak ↑.

2025-09-26 (Juma) — Admin mini-panel API

Vazifalar:

GET /admin/reports?status=open, POST /admin/reports/{id}/close

GET /admin/stats/overview (MVP: sessions today, dau)

audit_logs ga admin actions yozish.

Natija: Moderatsiya ishlaydi.

Tekshiruv: Report create → admin close.

2025-09-27 (Shanba) — Swagger & Validation

Vazifalar:

swaggo/swag bilan OpenAPI generate (comments).

Global validation (e.g., go-playground/validator) request DTO larda.

Error response formatini birxillashtirish.

Natija: /swagger/index.html ochiladi.

Tekshiruv: Endpoints ko‘rinadi, manual try.

2025-09-28 (Yakshanba) — Logging & Metrics

Vazifalar:

zap logger (siz yozgansiz) → request/response middlewares (redact PII).

Prometheus /metrics (ixtiyoriy), yoki hech bo‘lmasa p99 latency log.

Request ID (trace) qo‘shish.

Natija: Observability minimal.

Tekshiruv: Loglarda namespace, request id chiqadi.

2025-09-29 (Dushanba) — Security & Rate limiting

Vazifalar:

JWT expiry/refresh strategiya (MVP: access 15m, refresh 7d).

Rate limit (Redis) — auth va messages uchun oddiy limiter.

CORS, size limit, allowed content types.

Natija: Xavfsizroq API.

Tekshiruv: Limit bosilganda 429.

2025-09-30 (Seshanba) — E2E happy-path testi

Vazifalar:

Postman collection yoki ghz/k6 bilan kichik load.

Full flow: signup → profile → presence → match → chat → end → streak/badge → report/close.

Natija: Skriptlar bor.

Tekshiruv: 10–20 parallel session mini test.

2025-10-01 (Chorshanba) — Clean-up & polish

Vazifalar:

Dead code, err handling, TODO larni yopish.

README.md (run, env, migrate), .env.example.

Migrations versiya nomlash.

Natija: Repo silliq.

Tekshiruv: make up-local zero-to-run.

2025-10-02 (Payshanba) — Deploy (Render yoki ECS)

Vazifalar:

Render Web Service (8080), env’lar, health.

S3, Redis (managed) ulash.

Domain + HTTPS (agar bor).

Natija: Public URL’da /health, /swagger.

Tekshiruv: Real URL’da auth-flow ishlaydi.

2025-10-03 (Juma) — V1 freeze & smoke test

Vazifalar:

Real smoke test (2 device) → match, chat, end.

Loglarni ko‘rish, top metrics (RPS, error rate).

Release tag: v1.0.0.

Natija: V1 tayyor.

Tekshiruv: Bugfixlar (kritik bo‘lsa) darrov.