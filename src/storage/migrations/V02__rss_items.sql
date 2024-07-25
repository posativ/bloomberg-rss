CREATE TABLE rss_items
(
    url      TEXT      NOT NULL,
    category TEXT      NOT NULL,
    pub_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (url, category)
)
