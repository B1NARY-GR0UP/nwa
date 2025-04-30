CREATE TABLE hello_world (
                             id INT PRIMARY KEY,
                             message VARCHAR(100)
);

-- Insert a hello world message
INSERT INTO hello_world (id, message)
VALUES (1, 'Hello, World!');

-- Query the message
SELECT message FROM hello_world WHERE id = 1;