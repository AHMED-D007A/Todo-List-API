CREATE TABLE lists (
    list_id SERIAL PRIMARY KEY,            -- Auto-incrementing list ID
    user_id INTEGER NOT NULL,              -- Foreign key referencing users table
    list_title VARCHAR(100) NOT NULL,      -- Title of the list
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the list is created
    
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE -- If a user is deleted, their lists are deleted as well
);
