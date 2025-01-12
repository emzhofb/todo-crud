CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" text NOT NULL,
  "status" varchar CHECK (status IN ('pending', 'completed')) NOT NULL DEFAULT 'pending',
  "due_date" timestamptz NOT NULL DEFAULT (now())
);
