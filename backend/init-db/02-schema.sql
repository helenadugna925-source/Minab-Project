-- Core tables are already present; add functions, triggers, and computed helpers.

-- Computed field: event_total_bookmarks
CREATE OR REPLACE FUNCTION public.event_total_bookmarks(e events)
RETURNS integer
LANGUAGE sql
STABLE
AS $$
  SELECT COUNT(*) FROM event_bookmarks eb WHERE eb.event_id = e.id;
$$;

-- Computed field: event_total_followers
CREATE OR REPLACE FUNCTION public.event_total_followers(e events)
RETURNS integer
LANGUAGE sql
STABLE
AS $$
  SELECT COUNT(*) FROM event_followers ef WHERE ef.event_id = e.id;
$$;

-- Computed field: event_lowest_ticket_price
CREATE OR REPLACE FUNCTION public.event_lowest_ticket_price(e events)
RETURNS numeric
LANGUAGE sql
STABLE
AS $$
  SELECT MIN(t.price) FROM tickets t WHERE t.event_id = e.id;
$$;

-- Trigger: keep events.updated_at in sync
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_events_updated_at ON events;
CREATE TRIGGER trg_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Trigger: soft constraint to ensure at least one ticket per event (optional)
CREATE OR REPLACE FUNCTION public.ensure_ticket_on_event()
RETURNS trigger AS $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM tickets t WHERE t.event_id = NEW.id) THEN
    RAISE EXCEPTION 'Event must have at least one ticket';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- You can enable this by attaching to events AFTER INSERT/UPDATE if desired.

-- Ensure core users and events-related tables exist
CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  email text UNIQUE NOT NULL,
  password text NOT NULL,
  name text,
  phone text,
  created_at timestamptz DEFAULT now()
);

-- Backwards-compat: ensure 'phone' column exists for older installations
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone text;

CREATE TABLE IF NOT EXISTS events (
  id serial PRIMARY KEY,
  title text NOT NULL,
  description text,
  event_date date,
  price numeric DEFAULT 0,
  location_lat numeric,
  location_lng numeric,
  venue_name text,
  address text,
  category_id integer,
  user_id integer,
  featured_image text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS event_images (
  id serial PRIMARY KEY,
  event_id integer NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  image_url text NOT NULL
);

CREATE TABLE IF NOT EXISTS event_tags (
  id serial PRIMARY KEY,
  event_id integer NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name text NOT NULL
);

-- Materialized view for fast searches (title/description/tags)
DROP MATERIALIZED VIEW IF EXISTS public.events_search_mv;
CREATE MATERIALIZED VIEW public.events_search_mv AS
SELECT
  e.id,
  e.title,
  e.description,
  e.category_id,
  e.event_date,
  coalesce(string_agg(t.name, ' '), '') AS tag_text
FROM events e
LEFT JOIN event_tags t ON t.event_id = e.id
GROUP BY e.id;

CREATE INDEX IF NOT EXISTS idx_events_search_mv_title ON events_search_mv USING gin (to_tsvector('english', title));
CREATE INDEX IF NOT EXISTS idx_events_search_mv_desc ON events_search_mv USING gin (to_tsvector('english', description));
CREATE INDEX IF NOT EXISTS idx_events_search_mv_tags ON events_search_mv USING gin (to_tsvector('english', tag_text));

-- Add payment reference columns if not present
ALTER TABLE ticket_sales ADD COLUMN IF NOT EXISTS tx_ref text;
ALTER TABLE ticket_sales ADD COLUMN IF NOT EXISTS checkout_url text;

