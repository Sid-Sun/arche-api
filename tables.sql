-- Table structure for table `Users`
create TABLE dbo.Users
(
    user_id          int identity not null
        constraint Users_pk
            primary key,
    email            varchar(255) not null
        constraint Users_email_index
            unique,
    encryption_key   varchar(255) not null,
    key_hash         varchar(255) not null,
    verification_key varchar(255) not null,
    verified         bit          not null default 0
)

-- Table structure for table `Folders`
create table dbo.Folders
(
    folder_id int identity not null
        constraint Folders_pk
            primary key,
    user_id   int          not null,
    name      varchar(200) not null
)

-- Table structure for table `Notes`
create table Notes
(
    note_id   int identity not null
        constraint Note_pk
            primary key,
    folder_id int          not null,
    data      varchar(max) not null,
    name      varchar(200) not null
)