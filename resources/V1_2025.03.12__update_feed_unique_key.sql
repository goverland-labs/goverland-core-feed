drop index if exists feed_items_dao_proposal_uindex;

create unique index if not exists feed_items_unique_index
    on feed_items (dao_id, proposal_id, type, action)
    where type != 'delegate';
