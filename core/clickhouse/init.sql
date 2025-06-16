CREATE DATABASE IF NOT EXISTS stats;

CREATE TABLE IF NOT EXISTS stats.post_events (
    event_time DateTime,
    post_id String,
    user_id String,
    event_type String
) ENGINE = MergeTree()
ORDER BY (post_id, event_time);

CREATE MATERIALIZED VIEW IF NOT EXISTS stats.post_stats_daily
ENGINE = SummingMergeTree()
ORDER BY (post_id, date)
AS SELECT
    post_id,
    toDate(event_time) as date,
    countIf(event_type = 'view') as views,
    countIf(event_type = 'like') as likes,
    countIf(event_type = 'comment') as comments
FROM stats.post_events
GROUP BY post_id, date;
