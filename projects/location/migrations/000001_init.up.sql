CREATE EXTENSION cube;
CREATE EXTENSION earthdistance;

CREATE TABLE drivers_locations
(
    id  UUID,
    lat FLOAT8,
    lng FLOAT8,
    PRIMARY KEY(id)
);
