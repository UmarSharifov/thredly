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
alter table tred add column parent_tred_id int ;
alter table tred add column category_id int ;
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

CREATE TABLE treds_relations (
    id INT AUTO_INCREMENT PRIMARY KEY,
	tred_parent INT,
    foreign key (tred_parent) REFERENCES tred(id),
    tred_child INT,
    foreign key (tred_child) REFERENCES tred(id)
);

CREATE TABLE category (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50)
);

INSERT INTO category (name) VALUES
('Бакалавры'),
('Магистры'),
('Абитуриенты'),
('Преподаватели'),
('Экзамены'),
('Пересдача'),
('Общежитие'),
('Иностранцам'),
('Олимпиады'),
('Мероприятия'),
('Волонтерство');


CREATE TABLE t_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    publication_date DATETIME,
    views_count INT,
    content TEXT,
    category_id INT,
    photo VARCHAR(255)
);

INSERT INTO t_events (user_id, publication_date, views_count, content, photo) VALUES
(1, NOW(), 100, 'Содержание треда 1', 1, 'photo1.jpg');

CREATE TABLE event_category (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50)
);

INSERT INTO event_category (name) VALUES
('Университет'),
('Общежитие'),
('Свободное от учебы время'),
('Волонтерство'),
('Спорт');


CREATE TABLE complaint (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    tred_id INT,
    complaint_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES userAccount(id),
    FOREIGN KEY (tred_id) REFERENCES tred(id)
);


CREATE TABLE userSubscription (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subscriber_id INT NOT NULL,
    subscribed_to_id INT NOT NULL,
    subscription_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscriber_id) REFERENCES userAccount(id) ON DELETE CASCADE,
    FOREIGN KEY (subscribed_to_id) REFERENCES userAccount(id) ON DELETE CASCADE
);