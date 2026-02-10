INSERT INTO app_admin_users (username, password, real_name, email, phone, status, first_login, created_time, updated_time)
VALUES (
  'admin',
  MD5(CONCAT('admin123', 'super_admin_salt_2026')),
  '超级管理员',
  'admin@company.com',
  '',
  1,
  0,
  UNIX_TIMESTAMP(),
  UNIX_TIMESTAMP()
);