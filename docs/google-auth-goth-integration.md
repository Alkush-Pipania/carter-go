# Google Authentication (Login + Register) with Goth

This guide is tailored to the current `carter-go` architecture.

## 1) Current auth flow in this project

Today, auth is email/password based and works like this:

- Public auth routes are mounted under `/api/v1/auth`.
- `POST /login` validates credentials and creates a session in Postgres + Redis cache.
- `POST /register` creates a user using the user module.
- Session is stored in `session_id` cookie and validated by middleware on protected routes.

This is good foundation for OAuth because we can reuse **the same session creation path** after Google callback succeeds.

## 2) Why Goth fits this codebase

Goth gives you:

- Provider setup (`google.New(...)`)
- Start flow (`gothic.BeginAuthHandler`)
- Callback completion (`gothic.CompleteUserAuth`)

That means we only need to add project-specific logic for:

1. Finding/creating local users
2. Linking rows in `oauth_accounts`
3. Creating local session and cookie (same as normal login)

## 3) Data model alignment

You already have `oauth_accounts` table and sqlc queries:

- `GetOAuthAccount(provider, provider_user_id)`
- `CreateOAuthAccount(...)`

So no schema change is required for Google login/register.

## 4) Implementation plan (recommended)

### Step A — Add env config

Add these keys to config:

- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_CALLBACK_URL` (example: `http://localhost:8080/api/v1/auth/google/callback`)
- `OAUTH_SESSION_SECRET` (used by Gothic cookie store)
- optional `APP_URL`

### Step B — Initialize Goth once on startup

In `cmd/api/main.go`, after loading config and before serving routes:

1. Create Gothic cookie store (`gorilla/sessions`)
2. Configure provider:
   - `google.New(clientID, clientSecret, callbackURL, "email", "profile")`
3. Register provider with `goth.UseProviders(...)`

### Step C — Add auth routes

In `internal/modules/authentication/routes.go`, add:

- `GET /google/login`
- `GET /google/callback`

under existing `/api/v1/auth` mount.

### Step D — Add handler methods

In `internal/modules/authentication/handler.go`:

1. `GoogleLogin(w, r)`
   - set query `provider=google`
   - call `gothic.BeginAuthHandler(w, r)`

2. `GoogleCallback(w, r)`
   - set query `provider=google`
   - `gothic.CompleteUserAuth(w, r)` to get provider user
   - call service method `LoginWithGoogle(providerUserID, email)`
   - set existing `session_id` cookie and return 200

### Step E — Add service use-case

In `internal/modules/authentication/service.go`, add `LoginWithGoogle(...)`:

1. `GetOAuthAccount(google, providerUserID)`
2. If link exists:
   - use linked local user
3. If link does not exist:
   - find local user by email
   - if missing, create local user (`verified=true`, random password hash)
   - create `oauth_accounts` link row
4. Create local session (same as current login)
5. Cache session in Redis (best effort)
6. Return `LoginResponse` (`user_id`, `session_id`)

This supports both:

- **Register with Google** (new local user)
- **Login with Google** (existing linked user)

### Step F — Repository interface updates

In auth repository/service interfaces, expose:

- `GetOAuthAccount`
- `CreateOAuthAccount`
- `CreateUser` (if service creates users directly)

You already have generated sqlc methods for OAuth queries, so this is mostly plumbing.

### Step G — Frontend contract

Frontend can do:

- redirect browser to `GET /api/v1/auth/google/login`
- backend completes OAuth and sets `session_id` cookie
- frontend then calls protected APIs; auth middleware works unchanged

## 5) Google Cloud Console setup

1. Create OAuth Client ID (Web application).
2. Authorized redirect URI:
   - `http://localhost:8080/api/v1/auth/google/callback` (dev)
   - production callback URL for your deployed API
3. Add Authorized JavaScript origin(s) if needed.
4. Copy client ID/secret into environment.

## 6) Security notes for production

- Use HTTPS and secure cookies (`Secure=true`).
- Keep `OAUTH_SESSION_SECRET` strong and private.
- Restrict CORS origins to known frontend domains.
- Consider allowing only verified Google emails if business requires.
- Add logout cleanup for both local session and gothic session if needed.

## 7) Minimal endpoint behavior summary

- `GET /api/v1/auth/google/login` → redirect to Google consent
- `GET /api/v1/auth/google/callback` → create/link user, create local session, set cookie
- Existing protected routes continue using current middleware + session validation

## 8) Suggested dependency list

- `github.com/markbates/goth`
- `github.com/markbates/goth/gothic`
- `github.com/markbates/goth/providers/google`
- `github.com/gorilla/sessions`

