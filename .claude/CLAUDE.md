# NutriRecipe — Claude Code Instructions

## Current State
- [x] Repo init + git hooks
- [x] DB + migrations
- [x] Go module + models
- [x] Config
- [x] DB connection
- [x] USDA service
- [x] Agent service (mock)
- [x] Handlers
- [x] Main + routes
- [x] Vite init
- [x] API wrappers
- [x] Components
- [x] Pages
- [x] Vite proxy config

---

## Workflow Rules

### General
- Work through the build order one step at a time
- Do not create files not listed in the project structure without asking first
- Do not modify the DB schema without creating a new numbered migration file
- Do not guess at ambiguous requirements — ask before implementing
- After completing each build order step, stop and wait for confirmation before proceeding
- If any discrepancy is found between this spec and actual implementation, update this file immediately before continuing. Discrepancies include: type changes, added/removed fields, renamed files, route changes, or any deviation from the documented structure

### Branching
- Always check current branch before starting: `git branch`
- Create a new branch per build step: `git checkout -b feat/<name>`
- Branch names must correspond to the current build order step
- Never commit directly to `main` or `dev`
- Do not merge — stop and wait for user to review and merge

### Commits
- One commit per logical unit (one handler, one component, one migration)
- Maximum ~200 lines changed per commit — split if needed
- Commit message must follow Conventional Commits exactly:
  `type(scope): description`
  Valid types: `feat`, `fix`, `chore`, `test`
  Examples: `feat(usda): add caching layer`, `test(agent): add mock parse tests`
- Description must be 15 words or fewer
- No "Co-authored-by" lines of any kind in commit messages
- Run all tests before every commit — do not commit if tests fail

### File size
- No file may exceed 300 lines
- If a file approaches 250 lines, split it before continuing
- One handler file per route group
- One service file per external dependency

### Testing
- Every handler must have a `_test.go` file created in the same commit
- Every service function must have tests covering: happy path, empty/null input, error case
- Every frontend API wrapper must have a vitest test
- Test files live next to the code they test
- Backend: `go test ./...` must pass before any commit
- Frontend: `npm run test -- --run` must pass before any commit

---

## Git Hooks

Create these files at repo init and run `scripts/setup-hooks.sh` once.

> **Discrepancy (documented):** The pre-commit hook guards `go test ./...` and `npm run test -- --run` with existence checks (`[ -f backend/go.mod ]`, `[ -f frontend/package.json ]`) so that early build steps (DB migrations, Go module init) can be committed before their test runtimes exist. This deviates from the hook shown below but is required for the build order to function correctly.

### `scripts/setup-hooks.sh`
```bash
#!/bin/bash
cp scripts/hooks/* .git/hooks/
chmod +x .git/hooks/*
echo "Git hooks installed."
```

### `scripts/hooks/pre-commit`
```bash
#!/bin/bash
cd backend && go test ./...
if [ $? -ne 0 ]; then echo "❌ Backend tests failed."; exit 1; fi

cd ../frontend && npm run test -- --run
if [ $? -ne 0 ]; then echo "❌ Frontend tests failed."; exit 1; fi

LARGE=$(git diff --cached --name-only | xargs wc -l 2>/dev/null | awk '$1 > 300 {print $2}')
if [ -n "$LARGE" ]; then echo "❌ Files exceed 300 lines: $LARGE"; exit 1; fi

echo "✅ All checks passed."
```

### `scripts/hooks/commit-msg`
```bash
#!/bin/bash
MSG=$(cat "$1")

# Block co-authored-by
if echo "$MSG" | grep -qi "co-authored-by"; then
  echo "❌ Co-authored-by not allowed in commit messages."
  exit 1
fi

# Enforce conventional commits
PATTERN="^(feat|fix|chore|test)\(.+\): .+"
if ! echo "$MSG" | grep -qE "$PATTERN"; then
  echo "❌ Commit message must match: type(scope): description"
  exit 1
fi

# Enforce 15-word limit on description
DESC=$(echo "$MSG" | head -1 | sed 's/^[^:]*: //')
WORDS=$(echo "$DESC" | wc -w)
if [ "$WORDS" -gt 15 ]; then
  echo "❌ Commit description exceeds 15 words ($WORDS)."
  exit 1
fi
```

---

## Project Structure

```
nutrirecipe/
├── .claude/
│   └── CLAUDE.md
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── IngredientInput.tsx
│   │   │   ├── ClarificationDialog.tsx
│   │   │   ├── ReviewTable.tsx
│   │   │   └── RecipeCard.tsx
│   │   ├── pages/
│   │   │   ├── NewRecipe.tsx
│   │   │   └── RecipeList.tsx
│   │   ├── api/
│   │   └── main.tsx
│   ├── index.html
│   └── vite.config.ts
├── backend/
│   ├── main.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   │       ├── 000001_init.up.sql
│   │       └── 000001_init.down.sql
│   ├── handlers/
│   │   ├── agent.go
│   │   ├── recipes.go
│   │   └── usda.go
│   ├── models/
│   │   └── models.go
│   ├── services/
│   │   ├── agent.go
│   │   └── usda.go
│   └── go.mod
├── scripts/
│   ├── setup-hooks.sh
│   └── hooks/
│       ├── pre-commit
│       └── commit-msg
└── docker-compose.yml
```

---

## Stack

- **Frontend:** Vite + React + TypeScript
- **Backend:** Go + Gin
- **Database:** PostgreSQL 16
- **Local dev:** Docker Compose (Postgres only), `go run ./...`, `vite dev`

---

## Environment

`backend/.env`:
```
DATABASE_URL=postgres://nutrirecipe:nutrirecipe@localhost:5432/nutrirecipe?sslmode=disable
USDA_API_KEY=DEMO_KEY
PORT=8080
```

Vite proxy in `vite.config.ts` forwards `/api/*` → `localhost:8080`.

---

## Database Schema

Migrations managed via [golang-migrate](https://github.com/golang-migrate/migrate).

**`000001_init.up.sql`**

```sql
CREATE TABLE recipes (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name           TEXT NOT NULL,
  total_weight_g NUMERIC NOT NULL,
  serving_size_g NUMERIC NOT NULL DEFAULT 100,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE recipe_ingredients (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id         UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
  name              TEXT NOT NULL,
  quantity_g        NUMERIC NOT NULL,
  weight_ratio      NUMERIC NOT NULL,
  portion           NUMERIC,
  unit              TEXT,
  fdc_id            TEXT,
  source            TEXT NOT NULL DEFAULT 'usda',
  notes             TEXT,
  calories_per_100g NUMERIC,
  protein_per_100g  NUMERIC,
  carbs_per_100g    NUMERIC,
  fat_per_100g      NUMERIC
);

CREATE TABLE usda_cache (
  fdc_id     TEXT PRIMARY KEY,
  name       TEXT NOT NULL,
  calories   NUMERIC,
  protein    NUMERIC,
  carbs      NUMERIC,
  fat        NUMERIC,
  fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE local_ingredients (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name              TEXT NOT NULL UNIQUE,
  calories_per_100g NUMERIC,
  protein_per_100g  NUMERIC,
  carbs_per_100g    NUMERIC,
  fat_per_100g      NUMERIC,
  source            TEXT NOT NULL DEFAULT 'user_defined',
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**`000001_init.down.sql`**

```sql
DROP TABLE IF EXISTS local_ingredients;
DROP TABLE IF EXISTS usda_cache;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
```

---

## Models (`models/models.go`)

```go
type AgentParseRequest struct {
    Input     string   `json:"input"`
    QAHistory []QAPair `json:"qa_history,omitempty"`
}

type QAPair struct {
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

type AgentParseResponse struct {
    Status      string            `json:"status"` // "resolved" | "needs_clarification"
    Questions   []ClarificationQ  `json:"questions,omitempty"`
    Ingredients []ParsedIngredient `json:"ingredients,omitempty"`
}

type ClarificationQ struct {
    ID            string   `json:"id"`
    IngredientRaw string   `json:"ingredient_raw"`
    Question      string   `json:"question"`
    Options       []string `json:"options"`
}

type ParsedIngredient struct {
    Name          string        `json:"name"`
    Portion       float64       `json:"portion"`
    Unit          string        `json:"unit"`
    QuantityG     float64       `json:"quantity_g"`
    SearchQueries SearchQueries `json:"search_queries"`
    Confidence    string        `json:"confidence"` // "high" | "low"
    Notes         string        `json:"notes,omitempty"`
}

type SearchQueries struct {
    Primary      string   `json:"primary"`
    Alternatives []string `json:"alternatives"`
}

type ValidatedIngredient struct {
    ParsedIngredient
    FdcID           string      `json:"fdc_id,omitempty"`
    FdcName         string      `json:"fdc_name,omitempty"`
    MatchStatus     string      `json:"match_status"` // "matched" | "low_confidence" | "unresolved"
    MatchWarning    string      `json:"match_warning,omitempty"`
    Options         []USDAOption `json:"options,omitempty"`
    CaloriesPer100g float64     `json:"calories_per_100g,omitempty"`
    ProteinPer100g  float64     `json:"protein_per_100g,omitempty"`
    CarbsPer100g    float64     `json:"carbs_per_100g,omitempty"`
    FatPer100g      float64     `json:"fat_per_100g,omitempty"`
}

type USDAOption struct {
    FdcID    string  `json:"fdc_id"`
    Name     string  `json:"name"`
    Category string  `json:"category"`
    Calories float64 `json:"calories"`
    Protein  float64 `json:"protein"`
    Carbs    float64 `json:"carbs"`
    Fat      float64 `json:"fat"`
}

type SaveRecipeRequest struct {
    Name        string               `json:"name"`
    Ingredients []ValidatedIngredient `json:"ingredients"`
}
```

---

## Routes

```
POST /api/agent/parse     → handlers.AgentParse
POST /api/usda/validate   → handlers.USDAValidate
POST /api/recipes         → handlers.SaveRecipe
GET  /api/recipes         → handlers.ListRecipes
GET  /api/recipes/:id     → handlers.GetRecipe
```

---

## Service Behavior

### `services/agent.go` (mock)

**Agent 1:** Returns `needs_clarification` if input contains `"syrup"` without a qualifier. Otherwise returns `resolved` with placeholder gram conversions.

**Agent 2:** Returns `"low_confidence"` if the USDA result name shares fewer than 2 words with the search query. Otherwise returns `"matched"`.

### `services/usda.go`

Lookup order per ingredient:
1. `local_ingredients` table
2. `usda_cache` table
3. USDA API — primary query
4. USDA API — alternatives in order, max 3 attempts
5. Mark unresolved if all fail

USDA endpoints:
```
GET https://api.nal.usda.gov/fdc/v1/foods/search?query=<q>&api_key=<key>&pageSize=5
GET https://api.nal.usda.gov/fdc/v1/food/<fdcId>?api_key=<key>
```

Nutrient extraction: Use `nutrient.number` field (NOT `nutrient.id`):
- Energy (kcal) — nutrient number **208**
- Protein — nutrient number **203**
- Total lipid (fat) — nutrient number **204**
- Carbohydrate, by difference — nutrient number **205**

Cache all results in `usda_cache` after first fetch.

---

## Frontend Behavior

### `NewRecipe.tsx` — three sequential steps:
1. **Input** — recipe name + ingredient textarea → `POST /api/agent/parse`
2. **Clarification** — render questions as radio buttons → re-POST with `qa_history`
3. **Review** — `ReviewTable` → user approves/overrides/manual entry → `POST /api/recipes`

### `ReviewTable.tsx` row states:
- `matched` → green badge, collapsed, editable
- `low_confidence` → yellow badge, expanded, options dropdown
- `unresolved` → red badge, expanded, manual macro entry (4 fields)

---

## Build Order

Execute one step at a time. Stop after each for confirmation.

0. Repo init — `git init`, branches `main` + `dev`, `.gitignore`, run `scripts/setup-hooks.sh`
1. DB + migrations — `docker-compose.yml`, migration files (`000001_init.up.sql`, `000001_init.down.sql`)
2. Go module init — `go.mod` with `gin`, `pgx`, `godotenv`, `golang-migrate`
3. Models — `models/models.go`
4. Config — `config/config.go`
5. DB connection — `db/db.go`
6. USDA service — `services/usda.go`
7. Agent service (mock) — `services/agent.go`
8. Handlers — `handlers/agent.go`, `handlers/usda.go`, `handlers/recipes.go`
9. Main + routes — `main.go`
10. Vite init — `npm create vite` (React + TypeScript)
11. API wrappers — `src/api/`
12. Components — `IngredientInput`, `ClarificationDialog`, `ReviewTable`, `RecipeCard`
13. Pages — `NewRecipe`, `RecipeList`
14. Vite proxy — `vite.config.ts`

---

## Out of Scope (v1)

- Authentication
- Recipe editing after save
- Cooking loss / cooked weight
- Serving size customization
- Real LLM agent integration