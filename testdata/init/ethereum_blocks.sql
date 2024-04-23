CREATE SCHEMA IF NOT EXISTS ethereum;

-- Schema taken from prod database (v2_blocks).
CREATE TABLE ethereum.blocks (
  number integer NOT NULL,
  hash VARCHAR(100) NOT NULL,
  parent_hash VARCHAR(100) NOT NULL,
  nonce VARCHAR(255) NOT NULL,
  mix_hash VARCHAR(255) NOT NULL,
  sha3_uncles VARCHAR(255) NOT NULL,
  logs_bloom TEXT,
  transactions_root VARCHAR(100) NOT NULL,
  state_root VARCHAR(100) NOT NULL,
  receipts_root VARCHAR(100) NOT NULL,
  miner VARCHAR(42) NOT NULL,
  difficulty decimal(38, 0) NOT NULL,
  total_difficulty decimal(38, 0) NOT NULL,
  size bigint NOT NULL,
  gas_limit bigint NOT NULL,
  gas_used bigint NOT NULL,
  base_fee_per_gas bigint,
  timestamp timestamp NOT NULL,
  uncles VARCHAR(255) ARRAY,
  num_of_transactions integer NOT NULL,
  extra_data_raw TEXT,
  extra_data TEXT,
  process_time timestamp,
  block_date DATE,
  PRIMARY KEY (number)
);