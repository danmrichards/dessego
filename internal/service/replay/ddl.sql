CREATE TABLE IF NOT EXISTS replay (
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
   data TEXT,
   legacy INTEGER DEFAULT 0
)