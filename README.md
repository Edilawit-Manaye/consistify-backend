
# Consistify (Backend)

A Go/Gin backend that tracks daily coding consistency on LeetCode & Codeforces, computes streaks, and sends personalized push reminders via Firebase Cloud Messaging. Built with clean architecture.

## Architecture (Clean Architecture)
- **Domain** — core entities & errors (`User`, `DailyConsistency`, `PlatformActivity`, `StreakInfo`).
- **Usecases** — business logic (`UserUsecase`, `ConsistencyUsecase`, `PlatformUsecase`).
- **Repositories** — data access interfaces/implementations (MongoDB).
- **Infrastructure** — external services: auth (JWT/bcrypt), database (Mongo), platform APIs (LeetCode, Codeforces), notifications (FCM), scheduler (cron), Gemini client.
- **Delivery** — HTTP layer: controllers, routers, middleware.

## Features
- JWT auth (register/login, profile update)
- Daily consistency checks per platform (LeetCode, Codeforces)
- Streak tracking (current & longest)
- Personalized reminder notifications at user-defined time/timezone
- Scheduled daily checks & hourly reminder dispatch (cron)

## Tech Stack
Go, Gin, MongoDB, JWT + bcrypt, Firebase Admin (FCM), Docker

## Run
```bash
go run ./Delivery
```
Requires `.env` with `SERVER_PORT`, `JWT_SECRET`, `MONGO_URI`, `FIREBASE_SERVICE_ACCOUNT_PATH`, `FIREBASE_PROJECT_ID`, `LEETCODE_API_BASE_URL`, `CODEFORCES_API_BASE_URL`.
