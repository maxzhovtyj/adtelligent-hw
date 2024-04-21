INSERT INTO sources (name)
SELECT concat('Source #', md5(random()::text))
FROM generate_series(1, 99)
RETURNING *;

select *
from sources;

select *
from campaigns;

select *
from campaigns_sources;

INSERT INTO campaigns (name)
SELECT concat('Campaign #', md5(random()::text))
FROM generate_series(1, 100)
RETURNING *;

DELETE
FROM campaigns_sources;
DELETE
FROM sources;
DELETE
FROM campaigns;

SELECT name FROM sources
UNION
SELECT name FROM campaigns;

SELECT source_id, count(source_id) AS count
FROM campaigns_sources
GROUP BY source_id
ORDER BY count DESC
LIMIT 5;

SELECT c.id, c.name, cs.id
FROM campaigns c
         LEFT JOIN campaigns_sources cs ON c.id = cs.campaign_id
ORDER BY c.id;

SELECT c.id, c.name
FROM campaigns c
         LEFT JOIN campaigns_sources cs ON c.id = cs.campaign_id
WHERE cs.id IS NULL;

SELECT s.id as source_id, c.id AS campaign_id, c.name AS campaign_name
FROM sources s
         INNER JOIN campaigns_sources cs ON cs.source_id = s.id
         INNER JOIN campaigns c ON cs.campaign_id = c.id
WHERE s.id = 500