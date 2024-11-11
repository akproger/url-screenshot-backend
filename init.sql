DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_tables WHERE tablename = 'urls') THEN
        CREATE TABLE urls (
            id SERIAL PRIMARY KEY,
            url VARCHAR(255) NOT NULL,
            gen VARCHAR(255) UNIQUE NOT NULL,
            user_ip VARCHAR(45) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_tables WHERE tablename = 'users') THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            role VARCHAR(20) DEFAULT 'user',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;
