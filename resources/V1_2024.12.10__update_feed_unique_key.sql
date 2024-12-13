drop index if exists feed_items_dao_proposal_uindex;

update feed_items set action = 'dao.updated' where type = 'dao';
update feed_items set action = 'proposal.updated' where type = 'proposal';

create unique index feed_items_unique_index on feed_items (dao_id, proposal_id, type, action);
