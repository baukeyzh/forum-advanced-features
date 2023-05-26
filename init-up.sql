CREATE TABLE IF NOT EXISTS users (
    id              integer         primary key autoincrement,
    email           varchar(50)     not null unique,
    username        varchar(50)     not null unique,
    password_hash   varchar(255)    not null,
    token           varchar(50)     unique,
    expire_at       date            
);

CREATE TABLE IF NOT EXISTS posts (
    id			integer         primary key autoincrement,
    user_id	    integer         not null,
    date       date,            
    title		varchar(255),
    content		varchar(255),
    image_name		varchar(255),
    foreign key (user_id)       references users(id)

);

CREATE TABLE IF NOT EXISTS categories (
    id			integer        primary key autoincrement,
    name		varchar(255)
);

CREATE TABLE IF NOT EXISTS posts_categories (
    id			integer         primary key autoincrement,
    post_id     integer         not null,
    category_id integer         not null,
    foreign key (post_id)       references posts(id),
    foreign key (category_id)   references categories(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id			integer         primary key autoincrement,
    user_id	integer,
    date       date,            
    post_id		integer,
    content		varchar(255),
    foreign key (user_id)       references users(id),
    foreign key (post_id)       references posts(id)
);

CREATE TABLE IF NOT EXISTS posts_likes (
    id			integer        primary key autoincrement,
    user_id	    integer,
    post_id		integer,
    type        boolean not null,
    unique      (post_id, user_id),
    foreign key (user_id)    references users(id)
    foreign key (post_id)    references posts(id)
);

CREATE TABLE IF NOT EXISTS comments_likes (
    id			integer           primary key autoincrement,
    user_id	    integer,
    comment_id	integer,
    type        boolean not null,
    unique      (comment_id, user_id),
    foreign key (user_id)        references users(id),
    foreign key (comment_id)     references comments(id)
);

CREATE TABLE IF NOT EXISTS activity (
    id			integer primary key autoincrement,
    user_id	    integer not null,
    comment_id	integer,
    post_id     integer,
    author_id	integer not null,
    is_shown bool default false, 
    action      string not null,
    created_at       date,
    foreign key (comment_id)        references comment(id),
    foreign key (user_id)        references users(id),
    foreign key (comment_id)     references comments(id),
    foreign key (author_id)        references users(id)
);
INSERT INTO categories (name) values ("GO");
INSERT INTO categories (name) values ("JS");
INSERT INTO categories (name) values ("PHP");
INSERT INTO categories (name) values ("HTML");
DELETE FROM categories WHERE categories.id > 4;
