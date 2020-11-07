CREATE TABLE IF NOT EXISTS player (
    id TEXT,
    idx INTEGER DEFAULT 0,
    grade_s INTEGER DEFAULT 0,
    grade_a INTEGER DEFAULT 0,
    grade_b INTEGER DEFAULT 0,
    grade_c INTEGER DEFAULT 0,
    grade_d INTEGER DEFAULT 0,
    sessions INTEGER DEFAULT 0,
    msg_rating INTEGER DEFAULT 0,
    desired_tendency INTEGER DEFAULT 0,
    UNIQUE (id, idx)
)