ALTER TABLE venues
    ADD COLUMN location TEXT;
ALTER TABLE venues
    ADD COLUMN date_start DATE;
ALTER TABLE venues
    ADD COLUMN date_end DATE;
ALTER TABLE venues
    ADD COLUMN created DATETIME DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE venues
    ADD COLUMN placement TEXT DEFAULT 'agenda' NOT NULL CHECK (placement IN ('agenda', 'what-else'));

CREATE TABLE opening_times
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    venue_id   INTEGER NOT NULL,
    days       TEXT    NOT NULL,
    time_start TEXT    NOT NULL,
    time_end   TEXT    NOT NULL,
    created    DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (venue_id) REFERENCES venues (id)
);