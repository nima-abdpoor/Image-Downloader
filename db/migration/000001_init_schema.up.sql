CREATE TABLE query
(
    id         bigserial PRIMARY KEY,
    query      varchar     NOT NULL,
    status     VARCHAR(10) NOT NULL,
    per_page   int         not null,
    page       int         not null,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE image_result
(
    "id"         bigserial PRIMARY KEY,
    "query_id"   BIGINT NOT NULL REFERENCES query (id) ON DELETE CASCADE,
    "image_url"  varchar,
    "image_data" varchar,
    "timestamp"  TIMESTAMPTZ DEFAULT NOW()
);

