CREATE TABLE IF NOT EXISTS world_tendency (
    id INTEGER PRIMARY KEY autoincrement,
    character_id TEXT,
    area_1 INTEGER DEFAULT 0,
    wb_1 INTEGER DEFAULT 0,
    lr_1 INTEGER DEFAULT 0,
    area_2 INTEGER DEFAULT 0,
    wb_2 INTEGER DEFAULT 0,
    lr_2 INTEGER DEFAULT 0,
    area_3 INTEGER DEFAULT 0,
    wb_3 INTEGER DEFAULT 0,
    lr_3 INTEGER DEFAULT 0,
    area_4 INTEGER DEFAULT 0,
    wb_4 INTEGER DEFAULT 0,
    lr_4 INTEGER DEFAULT 0,
    area_5 INTEGER DEFAULT 0,
    wb_5 INTEGER DEFAULT 0,
    lr_5 INTEGER DEFAULT 0,
    area_6 INTEGER DEFAULT 0,
    wb_6 INTEGER DEFAULT 0,
    lr_6 INTEGER DEFAULT 0,
    area_7 INTEGER DEFAULT 0,
    wb_7 INTEGER DEFAULT 0,
    lr_7 INTEGER DEFAULT 0,
    FOREIGN KEY(character_id) REFERENCES character(id)
);