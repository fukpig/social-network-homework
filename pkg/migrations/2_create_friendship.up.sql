CREATE TABLE friendship(
  user VARCHAR(250),
  friend VARCHAR(250),
  primary key(user, friend),
  key(friend, user)
);
