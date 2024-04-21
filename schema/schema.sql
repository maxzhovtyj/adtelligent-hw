CREATE TABLE IF NOT EXISTS campaigns
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sources
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS campaigns_sources
(
    id          SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns (id) ON DELETE CASCADE NOT NULL,
    source_id   INTEGER REFERENCES sources (id) ON DELETE CASCADE   NOT NULL
);

SELECT * FROM campaigns_sources;
SELECT * FROM campaigns;
SELECT * FROM sources;

TRUNCATE campaigns_sources;
DELETE FROM campaigns;
DELETE FROM sources;
