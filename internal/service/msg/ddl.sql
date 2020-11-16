CREATE TABLE IF NOT EXISTS message (
    id INTEGER PRIMARY KEY autoincrement,
    character_id TEXT,
    block_id INTEGER DEFAULT 0,
    posx REAL,
    posy REAL,
    posz REAL,
    angx REAL,
    angy REAL,
    angz REAL,
    msg_id INTEGER DEFAULT 0,
    main_msg_id INTEGER DEFAULT 0,
    add_msg_cate_id INTEGER DEFAULT 0,
    rating INTEGER DEFAULT 0,
    legacy INTEGER DEFAULT 0
)