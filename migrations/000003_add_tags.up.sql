CREATE TABLE tags
(
    id              serial          not null unique,
    title           varchar(255)    not null   
);

CREATE TABLE items_tags
(
    id              serial                                              not null unique,
    item_id         int references todo_items (id) on delete cascade    not null,
    tag_id          int references tags (id) on delete cascade          not null  
);