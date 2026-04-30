kron/todo.md#L1-200
# Kron — TODO

This file lists remaining tasks, priorities, and actionable steps for the Kron project. You will be doing the implementations; this document is a roadmap and checklist to help you pick up where you left off.

---

## High-level goals
- Allow users to create "jobs" (HTTP requests with expected responses and schedules).
- Persist jobs in PocketBase `jobs` collection and tie them to authenticated users.
- Provide a dashboard UI for creating, listing, editing, and deleting jobs (HTMX-driven).
- Implement a scheduler/worker to execute jobs, compare responses, and notify users on failures.
- Add tests and logging to make the app maintainable and observable.

---

## Priorities
- P0 (Critical): Make job creation work end-to-end from the dashboard form to PocketBase collection with correct `user_id`. Wire POST route.
- P1 (High): Implement job listing per-user and ensure the jobs list refreshes after creation. Provide success/error feedback in the UI (HTMX fragments).
- P2 (Medium): Implement scheduling/execution engine to run jobs on schedule and compare responses.
- P3 (Low): Add edit/delete UI and endpoints, notifications (email/webhook), tests, and robust error handling.

---

## App surface & current status (short)
- Routes:
  - `GET /` — public homepage (implemented)
  - `GET /login` + `POST /login` — login/register (implemented)
  - `GET /dash` — dashboard page (implemented)
  - `GET /dash/jobs` — returns jobs list (implemented)
  - `POST /dash/jobs` — create job (stub exists but not implemented, route currently commented out)
- DB: migration `migrations/1770728248_jobs.go` defines `jobs` collection with fields: `user_id` (relation), `name`, `target`, `request` (JSON), `expected_response` (JSON), `schedule` (JSON)
- Helpers: `utils.jobRecordToStruct` implemented to parse DB record into `models.Job`.

---

## Concrete tasks

### P0 — Implement job creation (blocker)
1. Wire the route
   - File: `kron/main.go`
   - Action: Register `POST /dash/jobs` to the handler `pJob` (or equivalent). Currently commented out.
   - Acceptance: Submitting the dashboard form triggers the POST handler.

2. Implement `pJob` handler
   - File: `kron/dashboard.go`
   - Responsibilities:
     - Parse and validate form input. Required fields in the dashboard form are:
       - `name` (optional but recommended)
       - `target` (required, URL)
       - `method` (required)
       - `request` (optional body)
       - `expected_response` (optional)
     - Construct the `request` JSON object: `{ "method": "<METHOD>", "body": "<body>", "headers": { ... } }`
     - Set `user_id` on the new record to the currently authenticated user's id (use `r.Auth()` or `r.Record()` equivalent from PocketBase API to retrieve auth user).
     - Save a new record into the `jobs` collection using the PocketBase API (create core.NewRecord(collection) or `r.App.NewRecord(...)` as appropriate).
     - Return an HTML fragment (small snippet) suitable for HTMX that indicates success or errors. The dashboard form expects a response to be placed in `#newjobres`.
     - On success, the fragment should either:
       - Trigger a jobs reload using the existing `htmx` trigger `reload-jobs` (e.g. include `<div hx-trigger="reload-jobs from:body"></div>`), or
       - Directly include the new job in the jobs list partial.
   - Edge cases:
     - Validate `target` is a valid URL; sanitize inputs.
     - Ensure JSON fields are saved in the format the migration expects (`JSONField`s).
   - Acceptance:
     - Form submits result in a new record in the `jobs` collection with `user_id` set.
     - UI shows success and jobs list updates.

### P1 — Jobs listing and per-user visibility
1. Adjust `gJobs` to only return jobs visible to the current user
   - File: `kron/dashboard.go`
   - Action: Use PocketBase query to filter `jobs` by `user_id` equal to the current user's id.
   - Acceptance: Users only see their own jobs.

2. Improve jobs list rendering
   - File: `kron/views/jobsList.go`
   - Action: Expand list items to include method, schedule summary, created time, and actions (edit/delete buttons with HTMX hooks).
   - Acceptance: Jobs list shows richer detail and action hooks.

### P2 — Scheduler / executor
1. Design choices
   - Options:
     - Built-in goroutine scheduler inside the same process (simple).
     - External worker (recommended if you need horizontal scaling).
   - For v1, a simple in-process scheduler that loads jobs on startup and on schedule changes is fine.

2. Implementation plan
   - Create a scheduler component:
     - File suggestions: `kron/scheduler.go` and `kron/worker.go`
     - Responsibilities:
       - Parse `schedule` JSON for each job into a typed `Schedule` (use `models.Schedule`).
       - Maintain a queue/timer to run jobs at the correct intervals.
       - Execute HTTP requests as specified by `job.Request`.
       - Compare response body (and possibly status) with `job.Expected_response`.
       - Log and/or persist results (maybe a new `job_runs` collection) and trigger notifications on failures.
   - Acceptance:
     - Jobs execute roughly on schedule and results are persisted or visible in logs.

3. Notifications
   - Add support for email/webhook notifications; require configuration for SMTP or webhook endpoint.
   - Consider adding a `notifications` collection or adding fields to `jobs` to configure notification targets.

### P3 — Edit/Delete, Tests, and Polishing
1. Edit/Delete endpoints & UI
   - Files: `kron/dashboard.go`, `kron/views/*`
   - Implement handlers and forms for editing job fields and deleting jobs safely (confirmations).

2. Tests
   - Add unit tests for `jobRecordToStruct`, schedule parsing, and scheduler logic.
   - Add integration tests for `pJob` (if feasible).

3. Logging & error handling
   - Improve logs in critical paths (job execution, job creation failures).
   - Return user-friendly error messages in the HTMX responses.

4. Security & validation
   - Ensure only authenticated users can access dashboard routes.
   - Rate-limit or restrict potentially dangerous targets (optional).
   - Sanitize and limit the size of `request.body` to prevent abuse.

---

## Implementation notes & hints
- PocketBase specifics:
  - Use the RequestEvent to access the current authenticated user. For example, get the auth record id and set it as `user_id` on new `jobs` records.
  - `core.Record` JSON fields typically accept raw JSON strings. Build the request/expected JSON with `encoding/json` and save as a string.
- HTMX:
  - Dashboard uses HTMX attributes in `views/dashHome.go`. Return small partials that HTMX can swap into `#newjobres` or `#jobslist`.
  - To trigger a jobs reload after successful creation, you can return an element with `hx-trigger="reload-jobs from:body"` or use `hx-swap` behavior.

---

## Acceptance criteria (overall)
- A logged-in user can create a job from the dashboard form; the job is persisted and associated with that user.
- The jobs list refreshes to show newly created jobs without a full page reload.
- Jobs execute according to their schedule (basic scheduling working).
- Basic error messages propagate to the UI and logs.

---

## Future enhancements (ideas)
- Job history and metrics (response time, success rate).
- Support for headers, content-types, authentication for outgoing requests.
- Retry/backoff logic and concurrency controls.
- Per-user quotas and admin dashboard.
- Dockerfile and deployment scripts.

---

## Quick checklist you can use while implementing
- [ ] Un-comment/register `POST /dash/jobs` route
- [ ] Implement `pJob` to create PocketBase job records with `user_id`
- [ ] Ensure `request` and `expected_response` saved as JSON
- [ ] Make jobs list show only the current user's jobs
- [ ] Add HTMX success fragment to trigger jobs reload
- [ ] Implement basic scheduler and job executor
- [ ] Add logging and minimal error handling
- [ ] Add edit/delete endpoints and UI
- [ ] Add unit tests for parser/utility functions

---

If you want, I can convert any of the items above into a concrete code patch or show the exact code snippets you'd add for `pJob` and route wiring — but per your direction I will not implement them now.
