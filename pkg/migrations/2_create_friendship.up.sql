CREATE TABLE friendship(
  user BIGINT,
  friend BIGINT,
  primary key(user, friend),
  key(friend, user)
);
