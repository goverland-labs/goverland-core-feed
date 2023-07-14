create table feed_items
(
    id            uuid primary key,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    dao_id        uuid,
    proposal_id   text,
    discussion_id text,
    type          text,
    action        text,
    snapshot      jsonb,
    timeline      jsonb
);

create index idx_feed_items_dao_proposal_discussion_ids_action on feed_items (dao_id, proposal_id, discussion_id, action);

create table subscribers
(
    id          uuid not null primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    webhook_url text
);

create table subscriptions
(
    id            uuid primary key,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    subscriber_id text,
    dao_id        uuid
);

create index idx_subscriptions_deleted_at
    on subscriptions (deleted_at);

create index idx_subscriptions_subscriber_id_dao_id
    on subscriptions (subscriber_id, dao_id);
