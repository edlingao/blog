CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY,
  title VARCHAR(255) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  description TEXT,
  md_url TEXT,
  reactions TEXT,
  comments_available INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
  id INTEGER PRIMARY KEY,
  post_id INT NOT NULL,
  comment_id INT,
  author VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  reactions TEXT,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(50) DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  emoji VARCHAR(10) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS post_tags (
  post_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO tags (name, emoji) VALUES
('books', '📚'),
('tutorials', '📖'),
('music', '🎵'),
('projects', '🛠️'),
('art', '🎨'),
('photography', '📷'),
('travel', '✈️'),
('food', '🍜'),
('technology', '💻'),
('gaming', '🎮'),
('movies', '🎬');

