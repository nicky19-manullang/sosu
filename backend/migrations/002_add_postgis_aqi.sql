ALTER TABLE aqi_readings ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

CREATE OR REPLACE FUNCTION sync_aqi_geom() RETURNS TRIGGER AS $$
BEGIN
  NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_aqi_geom ON aqi_readings;
CREATE TRIGGER trg_aqi_geom
BEFORE INSERT OR UPDATE ON aqi_readings
FOR EACH ROW EXECUTE FUNCTION sync_aqi_geom();

CREATE INDEX IF NOT EXISTS idx_aqi_geom ON aqi_readings USING GIST (geom);