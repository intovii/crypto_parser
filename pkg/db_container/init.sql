\c postgres
CREATE EXTENSION IF NOT EXISTS dblink;
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'currency') THEN
      PERFORM dblink_exec('dbname=postgres user=' || current_user, 'CREATE DATABASE currency');
   END IF;
END
$$;
\c currency
DO
$$
    BEGIN
        CREATE TABLE IF NOT EXISTS tickers
        (
            id                  serial PRIMARY KEY,
            ticker              text NOT NULL UNIQUE
        );

        CREATE TABLE IF NOT EXISTS records
        (
            id                  serial PRIMARY KEY,
            ticker_id           integer NOT NULL REFERENCES tickers (id) ON DELETE CASCADE,
            timestamp           BIGINT NOT NULL,
            price               TEXT NOT NULL
        
        );
        
        
        RAISE NOTICE 'Таблицы успешно созданы.';
    EXCEPTION
        WHEN OTHERS THEN
            RAISE EXCEPTION 'Ошибка при создании таблиц: %', SQLERRM;
    END
$$;
