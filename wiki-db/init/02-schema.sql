-- JOINs here
/*
CREATE VIEW page_snapshot AS 
SELECT snapshots.uuid AS snapshot, revision, page_id AS page
FROM snapshots
JOIN revisions ON snapshots.revision = revisions.uuid;
*/

CREATE VIEW page_category_slugs AS
SELECT pages.slug AS page, categories.slug AS category
FROM page_categories
JOIN pages ON pages.uuid = page_categories.page_id
JOIN categories ON categories.id = page_categories.category;
