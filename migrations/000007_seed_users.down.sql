DELETE FROM users
WHERE email IN (
                'john@example.com',
                'jane@example.com',
                'admin@example.com'
    );