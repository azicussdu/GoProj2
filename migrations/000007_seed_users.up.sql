INSERT INTO users (
    full_name,
    email,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
) VALUES
      (
          'John Doe',
          'john@example.com',
          '$2a$10$examplehash1',
          'teacher',
          true,
          NOW(),
          NOW()
      ),
      (
          'Jane Smith',
          'jane@example.com',
          '$2a$10$examplehash2',
          'admin',
          true,
          NOW(),
          NOW()
      ),
      (
          'Admin User',
          'admin@example.com',
          '$2a$10$examplehash3',
          'admin',
          true,
          NOW(),
          NOW()
      );