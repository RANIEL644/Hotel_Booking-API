-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd



create table Room (
    id int primary key,
    type varchar(255) not null, 
    amenities varchar(255),
    price int not null, 
    availability varchar(50) not null );

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop table Room;
