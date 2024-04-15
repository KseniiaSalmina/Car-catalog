-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS persons (
                                       "id" SERIAL PRIMARY KEY NOT NULL,
                                       "name" TEXT NOT NULL,
                                       "surname" TEXT NOT NULL,
                                       "patronymic" TEXT
);

CREATE INDEX IF NOT EXISTS name_surname_patronymic_persons_idx ON persons (name, surname, patronymic);

CREATE TABLE IF NOT EXISTS cars (
                                    "reg_num" TEXT PRIMARY KEY NOT NULL,
                                    "mark" TEXT NOT NULL,
                                    "model" TEXT NOT NULL,
                                    "year" INT,
                                    "owner_id" INT NOT NULL,
                                    FOREIGN KEY (owner_id) REFERENCES persons(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS reg_num_cars_idx ON cars USING HASH (reg_num);
CREATE INDEX IF NOT EXISTS mark_cars_idx ON cars USING HASH (mark);
CREATE INDEX IF NOT EXISTS model_cars_idx ON cars USING HASH (model);
CREATE INDEX IF NOT EXISTS year_cars_idx ON cars (year);
CREATE INDEX IF NOT EXISTS owner_id_cars_idx ON cars USING HASH (owner_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS persons;
DROP INDEX IF EXISTS name_surname_patronymic_persons_idx;

DROP TABLE IF EXISTS cars;
DROP INDEX IF EXISTS reg_num_cars_idx;
DROP INDEX IF EXISTS mark_cars_idx;
DROP INDEX IF EXISTS model_cars_idx;
DROP INDEX IF EXISTS year_cars_idx;
DROP INDEX IF EXISTS owner_id_cars_idx;
-- +goose StatementEnd
