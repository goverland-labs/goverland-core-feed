create table feed_items
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type       text,
    type_id    text,
    action     text,
    snapshot   jsonb
);

create index idx_feed_items_type_type_id_action
    on feed_items (type, type_id, action);
