create table feed_items
(
    id            bigserial
        primary key,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    dao_id        text,
    proposal_id   text,
    discussion_id text,
    -- dao, proposal, discussion, etc
    type          text,
    -- core.dao.created, core.proposal.voting.started, etc
    action        text,
    snapshot      jsonb
);

create index idx_feed_items_dao_proposal_discussion_ids_action
    on feed_items (dao_id, proposal_id, discussion_id, action);

create table subscribers
(
    id             text not null
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    webhook_url text
);
