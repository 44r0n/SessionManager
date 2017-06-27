-- CREATING SCHEMA.
DROP SCHEMA IF EXISTS sessionmanager;
CREATE SCHEMA sessionmanager;
USE sessionmanager;

DROP TABLE IF EXISTS configuration;
CREATE TABLE configuration (
  name VARCHAR(165) NOT NULL,
  value VARCHAR(165) NOT NULL,
  PRIMARY KEY (name)
);

DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(165) UNIQUE NOT NULL,
  email VARCHAR(165) UNIQUE NOT NULL,
  password VARCHAR(128) NOT NULL,
  status TINYINT DEFAULT 1,
  date_created DATETIME NOT NULL,
  PRIMARY KEY (id),
  FULLTEXT (username,password)
);

DROP TABLE IF EXISTS user_tokens;
CREATE TABLE user_tokens (
  user INT UNSIGNED NOT NULL,
  token VARCHAR(32) NOT NULL,
  last_date_used DATETIME NOT NULL,
  PRIMARY KEY (user),
  FOREIGN KEY (user) REFERENCES users(id)
);
