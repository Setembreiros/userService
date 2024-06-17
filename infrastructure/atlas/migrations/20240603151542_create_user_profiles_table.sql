CREATE TABLE userservice.user_profiles (
    user_id INT NOT NULL,
    full_name VARCHAR(255),
    bio TEXT,
    link VARCHAR(255),
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
      REFERENCES userservice.users(id)
);
