CREATE TABLE cities(
    city varchar(40) unique,
    country varchar(40),
    lattitude varchar(40),
    longitude varchar(40)
);


CREATE TABLE predictions(
    city varchar(40),
    temp numeric,
    date timestamp
);