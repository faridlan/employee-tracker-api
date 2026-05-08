-- Enable UUID extension if not yet enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Create Employees Table
CREATE TABLE "employees" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "name" VARCHAR(255) NOT NULL,
    "position" VARCHAR(100) NOT NULL,
    "office_location" VARCHAR(100) NOT NULL,
    "entry_date" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "deleted_at" TIMESTAMPTZ
);

-- 2. Create Categories Table
CREATE TABLE "categories" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "name" VARCHAR(255) NOT NULL,
    "name_norm" VARCHAR(255) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "deleted_at" TIMESTAMPTZ
);

-- 3. Create Products Table
CREATE TABLE "products" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "name" VARCHAR(255) NOT NULL,
    "name_norm" VARCHAR(255) UNIQUE NOT NULL,
    "category_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_product_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE RESTRICT
);

-- 4. Create Targets Table
CREATE TABLE "targets" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "employee_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "nominal" BIGINT NOT NULL,
    "month" INT NOT NULL,
    "year" INT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_target_employee" FOREIGN KEY ("employee_id") REFERENCES "employees" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_target_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE RESTRICT
);

-- 5. Create Achievements Table (Ledger Style)
CREATE TABLE "achievements" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    "target_id" UUID NOT NULL,
    "nominal" BIGINT NOT NULL,
    "description" TEXT,
    "closing_date" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_achievement_target" FOREIGN KEY ("target_id") REFERENCES "targets" ("id") ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_employees_deleted_at ON employees (deleted_at);

CREATE INDEX idx_categories_name_norm ON categories (name_norm);

CREATE INDEX idx_products_category_id ON products (category_id);

CREATE INDEX idx_targets_employee_period ON targets (employee_id, month, year);

CREATE INDEX idx_achievements_target_id ON achievements (target_id);