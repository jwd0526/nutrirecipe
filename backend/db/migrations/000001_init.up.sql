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
