\c courses

-- Вставка тестовых данных в таблицу users
INSERT INTO users (role, name, email, phone) VALUES
('admin', 'John Doe', 'john.doe@example.com', '123-456-7890'),
('student', 'Jane Smith', 'jane.smith@example.com', '098-765-4321'),
('teacher', 'Alice Johnson', 'alice.johnson@example.com', '555-555-5555');

-- Вставка тестовых данных в таблицу courses
INSERT INTO courses (title, description, category, thumbnail_url, price) VALUES
('Introduction to Python', 'Learn the basics of Python programming', 'Programming', 'http://example.com/python.jpg', 50),
('Advanced SQL', 'Master advanced SQL techniques', 'Database', 'http://example.com/sql.jpg', 100),
('Web Development', 'Build modern web applications', 'Web Development', 'http://example.com/webdev.jpg', 150);

-- Вставка тестовых данных в таблицу courses_content
INSERT INTO courses_content (title, description, url, photo_url, course_id) VALUES
('Lesson 1: Basics', 'Introduction to Python basics', 'http://example.com/lesson1', 'http://example.com/lesson1.jpg', 1),
('Lesson 2: Functions', 'Learn about Python functions', 'http://example.com/lesson2', 'http://example.com/lesson2.jpg', 1),
('Lesson 1: Joins', 'Introduction to SQL joins', 'http://example.com/lesson3', 'http://example.com/lesson3.jpg', 2),
('Lesson 2: Indexes', 'Learn about SQL indexes', 'http://example.com/lesson4', 'http://example.com/lesson4.jpg', 2),
('Lesson 1: HTML', 'Introduction to HTML', 'http://example.com/lesson5', 'http://example.com/lesson5.jpg', 3),
('Lesson 2: CSS', 'Learn about CSS', 'http://example.com/lesson6', 'http://example.com/lesson6.jpg', 3);

-- Вставка тестовых данных в таблицу users_courses
INSERT INTO users_courses (user_id, course_id) VALUES
(1, 1), -- Admin enrolled in Introduction to Python
(2, 2), -- Student enrolled in Advanced SQL
(3, 3); -- Teacher enrolled in Web Development
