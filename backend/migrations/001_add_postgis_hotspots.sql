ALTER TABLE hotspots ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

CREATE OR REPLACE FUNCTION sync_hotspot_geom() RETURNS TRIGGER AS $$
BEGIN
  NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_hotspot_geom ON hotspots;
CREATE TRIGGER trg_hotspot_geom
BEFORE INSERT OR UPDATE ON hotspots
FOR EACH ROW EXECUTE FUNCTION sync_hotspot_geom();

UPDATE hotspots
SET geom = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
WHERE geom IS NULL;

CREATE INDEX IF NOT EXISTS idx_hotspots_geom ON hotspots USING GIST (geom);
CREATE INDEX IF NOT EXISTS idx_hotspots_geom_geography ON hotspots USING GIST ((geom::geography));
CREATE INDEX IF NOT EXISTS idx_hotspots_acquired_at ON hotspots (acquired_at DESC);