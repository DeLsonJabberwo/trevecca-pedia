#!/bin/sh

curl -X POST 127.0.0.1:9454/pages/new \
-F "slug=spiritual-life" \
-F "name=Spiritual Life" \
-F "author=1197028" \
-F "archive_date=" \
-F "new_page=@new_page.md"

