-- +goose Up
INSERT INTO users (username, role, email, password_hash) VALUES
                                                             ('admin', 'admin', 'admin@email.com', '\x24326124313024385a313078336f514849444a7744435071575236564f5470766e6e4679716e46522f456a30535a516268614e79544f447a64736e61'),
                                                             ('test', 'user', 'test@email.com', '\x2432612431302443436e3231377541485a336d755a5a41426f726e514f7a423438744d7050667950702e2f4438554378396d34764f6a694c564c5a2e'),
                                                             ('john_doe_2007', 'user', 'johndoe2007@email.com', '\x24326124313024436d35507043473854766961666178312f324e2e774f7432582f6754596f484461316a6d456553355538787a3170302e6a5273354b');
-- admin admin123
-- test 12345678
-- john_doe_2007 megaSecret
-- +goose Down
DELETE FROM users;