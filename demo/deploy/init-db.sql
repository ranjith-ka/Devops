CREATE TABLE IF NOT EXISTS cases (
  id SERIAL PRIMARY KEY,
  title TEXT,
  details TEXT,
  sender TEXT,
  eta TIMESTAMP,
  sla_days INT,
  hypercare BOOLEAN,
  label TEXT,
  hs_code TEXT,
  preference TEXT,
  supp_units TEXT,
  assigned_to TEXT,
  priority_score INT,
  created_at TIMESTAMP DEFAULT now()
);
