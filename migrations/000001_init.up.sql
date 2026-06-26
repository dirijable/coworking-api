CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE rooms
(
    id       uuid DEFAULT uuidv7() PRIMARY KEY,
    name     varchar(32) NOT NULL UNIQUE CHECK ( char_length(name) >= 2 ),
    capacity int         NOT NULL CHECK ( capacity >= 1 AND capacity <= 1000 )
);

CREATE TABLE users
(
    id         uuid DEFAULT uuidv7() PRIMARY KEY,
    first_name varchar(64) NOT NULL CHECK ( char_length(first_name) >= 1 ),
    last_name  varchar(64) NOT NULL CHECK ( char_length(last_name) >= 1 ),
    email      varchar(64) NOT NULL UNIQUE
);

CREATE TABLE bookings
(
    id         uuid DEFAULT uuidv7() PRIMARY KEY,
    room_id    uuid REFERENCES rooms (id) ON DELETE CASCADE NOT NULL,
    user_id    uuid REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    start_time TIMESTAMPTZ                                  NOT NULL,
    end_time   TIMESTAMPTZ                                  NOT NULL,

    CONSTRAINT check_booking_dates CHECK (end_time > start_time),

    CONSTRAINT exclude_room_time_overlapping EXCLUDE USING gist (
        room_id WITH =,
        tstzrange(start_time, end_time) WITH &&
        )
);

CREATE INDEX IF NOT EXISTS idx_bookings_room_id ON bookings (room_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);
