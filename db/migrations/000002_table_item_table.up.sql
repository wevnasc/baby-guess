CREATE TABLE IF NOT EXISTS tables (
    id UUID DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    account_id uuid,
    PRIMARY KEY (id),
    CONSTRAINT fk_accounts
        FOREIGN KEY(account_id)
            REFERENCES accounts(id)
                ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    id uuid DEFAULT uuid_generate_v4(),
    description VARCHAR NOT NULL,
    account_id UUID,
    table_id UUID,
    status INT,
    PRIMARY KEY (id),
    CONSTRAINT fk_accounts
        FOREIGN KEY(account_id)
            REFERENCES accounts(id),
    CONSTRAINT fk_tables
        FOREIGN KEY(table_id)
            REFERENCES tables(id)
                ON DELETE CASCADE
);
