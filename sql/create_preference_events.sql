CREATE TABLE preference_events(
  user_id int not null, 
  image_id int not null, 
  liked boolean not null, 
  timestamp timestamp default current_timestamp
);
