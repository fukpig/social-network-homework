CREATE TABLE users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(200) NOT NULL,
  name  VARCHAR(250) NOT NULL,
  surname  VARCHAR(250) NOT NULL,
  password  VARCHAR(250) NOT NULL,
  sex  VARCHAR(250) NOT NULL,
  city  VARCHAR(250) NOT NULL,
  interests TEXT NULL,
  created_at TIMESTAMP
);
