INSERT INTO users (username, role, email, password_hash) VALUES
                                                             ('admin', 'admin', 'admin@email.com', 'admin123'),
                                                             ('test', 'user', 'test@email.com', '12345678'),
                                                             ('john_doe_2007', 'user', 'johndoe2007@email.com', 'megaSecret');

INSERT INTO tasks (title, text, priority, user_id) VALUES
                                                       ('Do some admin stuff', '', 2, 1),
                                                       ('Do some other admin stuff', 'idk', 2, 1),

                                                       ('Run end to end tests', 'the ones with full access', 3, 2),
                                                       ('Run integration tests', 'those which access db and web', 2, 2),
                                                       ('Run unit tests', 'the func tests', 1, 2),


                                                       ('Workout', 'go to the gym and workout', 2, 3),
                                                       ('Have a Meal', '', 0, 3);