CREATE TABLE wallet (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id uuid NOT NULL,
    balance integer NOT NULL DEFAULT 0,
    currency integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE products(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    description integer NOT NULL,
    min_transaction_amount integer NOT NULL,
    max_transaction_amount integer NOT NULL,
    authorised_max_transaction_amount integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone
);

CREATE TABLE transactions(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id uuid NOT NULL,
    product_id uuid NOT NULL,
    correlation_id uuid NOT NULL,
    amount integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone
);

-- unique constraints to prevent duplicate transactions
CREATE UNIQUE INDEX transactions_correlation_id_idx ON transactions(correlation_id);

-- index for transaction by ID
CREATE INDEX transactions_id_idx ON transactions(id);

-- index for wallet and product by ID
CREATE INDEX wallet_id_idx ON wallet(id);

