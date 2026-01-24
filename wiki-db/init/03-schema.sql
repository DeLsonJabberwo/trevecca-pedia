INSERT INTO pages VALUES ('07918316-875e-4581-87ab-5b8d1d8bdd3a', 'dan-boone', 'Dan Boone', NULL, NULL, NULL);
INSERT INTO pages VALUES ('60b6b10c-db33-4b4c-9dcf-566f5b3c59a4', 'newsies', 'Newsies', NULL, '2025-11-12 00:00:00', NULL);
INSERT INTO pages VALUES ('80d10989-2dbb-4d64-98a0-eb8739ccb77f', 'dan-boone-deep-state', 'Dan Boone Deep State', NULL, NULL, '2025-10-20 00:00:00');
INSERT INTO pages VALUES ('c32b2206-ca5f-407c-b6d0-89b7da1e4c2b', 'spiritual-life', 'Spiritual Life', NULL, NULL, NULL);

INSERT INTO revisions VALUES ('4a76899b-0051-444d-92cb-8017f09f2fea', '07918316-875e-4581-87ab-5b8d1d8bdd3a', '2025-12-18 20:18:09.115549', '0000001');
INSERT INTO revisions VALUES ('e98eb510-8a71-401b-9797-2d08e45151b9', '07918316-875e-4581-87ab-5b8d1d8bdd3a', '2025-12-18 20:18:00.000000', '1197028');
INSERT INTO revisions VALUES ('a9519a06-c79d-4f41-a3e4-0cb919d0127c', 'c32b2206-ca5f-407c-b6d0-89b7da1e4c2b', '2026-01-24 22:13:23.402381', '1197028');

INSERT INTO snapshots VALUES ('46976289-d152-43d3-a2d7-f75b3102e471', '07918316-875e-4581-87ab-5b8d1d8bdd3a', 'e98eb510-8a71-401b-9797-2d08e45151b9');
INSERT INTO snapshots VALUES ('fc6037c1-296c-4ffa-be7f-b81b893680c8', 'c32b2206-ca5f-407c-b6d0-89b7da1e4c2b', 'a9519a06-c79d-4f41-a3e4-0cb919d0127c');

INSERT INTO categories (slug, name) VALUES ('people', 'People');
INSERT INTO categories (slug, name) VALUES ('buildings', 'Buildings');
INSERT INTO categories (slug, name) VALUES ('departments', 'Departments');
INSERT INTO categories (slug, name) VALUES ('student-life', 'Student Life');
INSERT INTO categories (slug, name) VALUES ('miscellaneous', 'Miscellaneous');

INSERT INTO page_categories VALUES ('07918316-875e-4581-87ab-5b8d1d8bdd3a', '1');
INSERT INTO page_categories VALUES ('60b6b10c-db33-4b4c-9dcf-566f5b3c59a4', '4');
INSERT INTO page_categories VALUES ('80d10989-2dbb-4d64-98a0-eb8739ccb77f', '5');
INSERT INTO page_categories VALUES ('c32b2206-ca5f-407c-b6d0-89b7da1e4c2b', '4');
