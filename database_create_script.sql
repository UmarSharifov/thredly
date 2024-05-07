use snippetbox;

CREATE TABLE userAccount (
    id INT AUTO_INCREMENT PRIMARY KEY,
    last_name VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    photo VARCHAR(255),
    date_of_birth DATE,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    userLogin VARCHAR(100),
    userPassword VARCHAR(100)
);


INSERT INTO userAccount (last_name, first_name, photo, date_of_birth, email, phone_number, userLogin, userPassword) VALUES
('Иванов', 'Иван', 'ivan.jpg', '1990-05-15', 'ivanov@example.com', '+123456789', 'user1', 'user1'),
('Петров', 'Петр', 'petr.jpg', '1985-09-20', 'petrov@example.com', '+987654321', 'user2', 'user2'),
('Сидорова', 'Анна', 'anna.jpg', '1995-02-10', 'sidorova@example.com', '+111222333','user3', 'user3'),
('Козлов', 'Алексей', 'alex.jpg', '1988-12-03', 'kozlov@example.com', '+444555666','user4', 'user4'),
('Смирнова', 'Екатерина', 'katya.jpg', '1992-07-25', 'smirnova@example.com', '+777888999','user5', 'user5');


CREATE TABLE tred (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES user(id),
    publication_date DATETIME,
    views_count INT,
    content TEXT,
    photo VARCHAR(255)
);

INSERT INTO tred (user_id, publication_date, views_count, content, photo) VALUES
(1, NOW(), 100, 'Содержание треда 1', 'photo1.jpg'),
(2, NOW(), 150, 'Содержание треда 2', 'photo2.jpg'),
(3, NOW(), 200, 'Содержание треда 3', 'photo3.jpg'),
(4, NOW(), 250, 'Содержание треда 4', 'photo4.jpg'),
(5, NOW(), 300, 'Содержание треда 5', 'photo5.jpg');


CREATE TABLE thread_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    thread_id INT,
    FOREIGN KEY (thread_id) REFERENCES tred(id),
    tag VARCHAR(50)
);


INSERT INTO thread_tags (thread_id, tag) VALUES
(1, 'MySQL'),
(1, 'Database'),
(2, 'Programming'),
(3, 'Technology'),
(4, 'Science'),
(5, 'Art');
