-- GitHub-style sample data for SQL practice (PostgreSQL).
-- Loads 500+ rows: users, repositories, pull requests, workflow runs.
--
-- Usage:
--   psql -d mydb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
--   or in psql: \i path/to/github_practice_seed.sql
--
-- Example queries (see README.md in this folder for more):

SET client_min_messages TO NOTICE;

DROP TABLE IF EXISTS workflow_runs CASCADE;
DROP TABLE IF EXISTS pull_requests CASCADE;
DROP TABLE IF EXISTS repositories CASCADE;
DROP TABLE IF EXISTS github_users CASCADE;

-- ---------------------------------------------------------------------------
-- Schema
-- ---------------------------------------------------------------------------

CREATE TABLE github_users (
    id              SERIAL PRIMARY KEY,
    login           VARCHAR(100) NOT NULL UNIQUE,
    display_name    VARCHAR(200),
    company         VARCHAR(120),
    created_at      DATE NOT NULL
);

CREATE TABLE repositories (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    full_name       VARCHAR(260) NOT NULL UNIQUE,
    owner_login     VARCHAR(100) NOT NULL REFERENCES github_users (login),
    description     TEXT,
    stars           INT NOT NULL DEFAULT 0 CHECK (stars >= 0),
    is_private      BOOLEAN NOT NULL DEFAULT FALSE,
    default_branch  VARCHAR(100) NOT NULL DEFAULT 'main',
    created_at      DATE NOT NULL,
    UNIQUE (owner_login, name)
);

CREATE TABLE pull_requests (
    id              SERIAL PRIMARY KEY,
    repo_id         INT NOT NULL REFERENCES repositories (id) ON DELETE CASCADE,
    number          INT NOT NULL,
    title           TEXT NOT NULL,
    state           VARCHAR(20) NOT NULL CHECK (state IN ('open', 'closed')),
    merged          BOOLEAN NOT NULL DEFAULT FALSE,
    author_login    VARCHAR(100) NOT NULL REFERENCES github_users (login),
    head_branch     VARCHAR(200) NOT NULL,
    base_branch     VARCHAR(200) NOT NULL,
    additions       INT NOT NULL DEFAULT 0 CHECK (additions >= 0),
    deletions       INT NOT NULL DEFAULT 0 CHECK (deletions >= 0),
    changed_files   INT NOT NULL DEFAULT 0 CHECK (changed_files >= 0),
    created_at      TIMESTAMPTZ NOT NULL,
    merged_at       TIMESTAMPTZ,
    closed_at       TIMESTAMPTZ,
    UNIQUE (repo_id, number)
);

CREATE TABLE workflow_runs (
    id              SERIAL PRIMARY KEY,
    repo_id         INT NOT NULL REFERENCES repositories (id) ON DELETE CASCADE,
    pr_id           INT REFERENCES pull_requests (id) ON DELETE SET NULL,
    workflow_name   VARCHAR(200) NOT NULL,
    event           VARCHAR(50) NOT NULL,
    status          VARCHAR(30) NOT NULL,
    conclusion      VARCHAR(30),
    branch          VARCHAR(200) NOT NULL,
    run_attempt     SMALLINT NOT NULL DEFAULT 1 CHECK (run_attempt >= 1),
    duration_seconds INT CHECK (duration_seconds IS NULL OR duration_seconds >= 0),
    created_at      TIMESTAMPTZ NOT NULL,
    head_sha        CHAR(40) NOT NULL
);

CREATE INDEX idx_pr_repo ON pull_requests (repo_id);
CREATE INDEX idx_pr_author ON pull_requests (author_login);
CREATE INDEX idx_pr_state_merged ON pull_requests (state, merged);
CREATE INDEX idx_wr_repo ON workflow_runs (repo_id);
CREATE INDEX idx_wr_pr ON workflow_runs (pr_id);
CREATE INDEX idx_wr_conclusion ON workflow_runs (conclusion);

-- ---------------------------------------------------------------------------
-- Seed: 25 users
-- ---------------------------------------------------------------------------

INSERT INTO github_users (login, display_name, company, created_at) VALUES
    ('ranjith-dev', 'Ranjith A', 'Acme Cloud', '2018-03-12'),
    ('hema-s', 'Hema S', 'TechWithHema', '2019-01-04'),
    ('alex-k', 'Alex Kumar', 'Acme Cloud', '2017-11-20'),
    ('priya-m', 'Priya Menon', NULL, '2020-06-01'),
    ('devops-bot', 'DevOps Bot', 'Acme Cloud', '2021-02-14'),
    ('ravi-p', 'Ravi P', 'StartupXYZ', '2019-09-09'),
    ('sneha-r', 'Sneha R', 'StartupXYZ', '2020-04-22'),
    ('amit-v', 'Amit V', 'FinTech Co', '2016-05-30'),
    ('noor-f', 'Noor F', 'FinTech Co', '2018-08-15'),
    ('li-wei', 'Li Wei', 'OpenSource', '2015-12-01'),
    ('maria-g', 'Maria G', 'OpenSource', '2019-03-28'),
    ('omar-h', 'Omar H', NULL, '2021-07-07'),
    ('yuki-t', 'Yuki T', 'Acme Cloud', '2022-01-10'),
    ('ben-c', 'Ben C', 'StartupXYZ', '2020-11-11'),
    ('fatima-z', 'Fatima Z', 'FinTech Co', '2017-04-18'),
    ('vikram-s', 'Vikram S', 'Acme Cloud', '2018-10-25'),
    ('ananya-k', 'Ananya K', 'TechWithHema', '2023-02-01'),
    ('carlos-m', 'Carlos M', 'OpenSource', '2014-06-06'),
    ('diana-l', 'Diana L', 'StartupXYZ', '2021-09-19'),
    ('erik-n', 'Erik N', NULL, '2022-05-05'),
    ('fiona-q', 'Fiona Q', 'FinTech Co', '2019-12-12'),
    ('george-w', 'George W', 'Acme Cloud', '2016-01-23'),
    ('hana-y', 'Hana Y', 'TechWithHema', '2023-08-08'),
    ('ivan-p', 'Ivan P', 'OpenSource', '2015-09-17'),
    ('julia-x', 'Julia X', 'StartupXYZ', '2020-02-29');

-- ---------------------------------------------------------------------------
-- Seed: 20 repositories
-- ---------------------------------------------------------------------------

INSERT INTO repositories (name, full_name, owner_login, description, stars, is_private, default_branch, created_at) VALUES
    ('devops', 'ranjith-dev/devops', 'ranjith-dev', 'Infra & tutorials', 420, FALSE, 'main', '2019-04-01'),
    ('backend-api', 'ranjith-dev/backend-api', 'ranjith-dev', 'REST API service', 128, FALSE, 'main', '2020-02-10'),
    ('frontend-app', 'hema-s/frontend-app', 'hema-s', 'React dashboard', 89, FALSE, 'main', '2021-03-15'),
    ('data-pipeline', 'alex-k/data-pipeline', 'alex-k', 'ETL jobs', 56, TRUE, 'develop', '2020-08-20'),
    ('mobile-sdk', 'priya-m/mobile-sdk', 'priya-m', 'Mobile SDK', 210, FALSE, 'main', '2019-11-05'),
    ('infra-terraform', 'devops-bot/infra-terraform', 'devops-bot', 'Terraform modules', 340, FALSE, 'main', '2021-01-01'),
    ('auth-service', 'ravi-p/auth-service', 'ravi-p', 'OAuth2 / JWT', 95, FALSE, 'main', '2022-04-12'),
    ('payments', 'sneha-r/payments', 'sneha-r', 'Payment gateway', 67, TRUE, 'main', '2021-06-30'),
    ('risk-engine', 'amit-v/risk-engine', 'amit-v', 'Risk scoring', 44, TRUE, 'develop', '2018-09-01'),
    ('ledger', 'noor-f/ledger', 'noor-f', 'Double-entry ledger', 31, FALSE, 'main', '2019-12-20'),
    ('docs-site', 'li-wei/docs-site', 'li-wei', 'Documentation', 512, FALSE, 'main', '2016-03-03'),
    ('cli-tool', 'maria-g/cli-tool', 'maria-g', 'CLI utilities', 178, FALSE, 'main', '2020-10-10'),
    ('ml-models', 'omar-h/ml-models', 'omar-h', 'Model training', 92, FALSE, 'main', '2023-01-15'),
    ('edge-worker', 'yuki-t/edge-worker', 'yuki-t', 'Edge compute', 55, FALSE, 'main', '2022-11-22'),
    ('notifications', 'ben-c/notifications', 'ben-c', 'Email & push', 38, FALSE, 'develop', '2021-04-04'),
    ('compliance', 'fatima-z/compliance', 'fatima-z', 'Audit logs', 22, TRUE, 'main', '2018-07-07'),
    ('cache-layer', 'vikram-s/cache-layer', 'vikram-s', 'Redis wrappers', 61, FALSE, 'main', '2019-05-18'),
    ('search-index', 'ananya-k/search-index', 'ananya-k', 'OpenSearch', 47, FALSE, 'main', '2023-06-01'),
    ('legacy-monolith', 'carlos-m/legacy-monolith', 'carlos-m', 'Old stack', 15, TRUE, 'master', '2014-01-01'),
    ('sandbox', 'diana-l/sandbox', 'diana-l', 'Experiments', 8, FALSE, 'main', '2022-02-14');

-- ---------------------------------------------------------------------------
-- Seed: 300 pull requests (15 per repo × 20 repos)
-- ---------------------------------------------------------------------------

INSERT INTO pull_requests (
    repo_id, number, title, state, merged, author_login,
    head_branch, base_branch, additions, deletions, changed_files,
    created_at, merged_at, closed_at
)
SELECT
    ((g - 1) / 15) + 1 AS repo_id,
    ((g - 1) % 15) + 1 AS number,
    CASE (g % 5)
        WHEN 0 THEN 'feat: add workflow for integration tests #' || g::text
        WHEN 1 THEN 'fix: flaky PR checks on ' || ((g % 20) + 1)::text
        WHEN 2 THEN 'chore: bump deps and lockfile #' || g::text
        WHEN 3 THEN 'docs: update README for API v' || ((g % 3) + 1)::text
        ELSE 'refactor: extract module ' || ((g % 50) + 1)::text
    END AS title,
    CASE WHEN (g % 7) = 0 THEN 'open' ELSE 'closed' END AS state,
    CASE WHEN (g % 7) = 0 THEN FALSE WHEN (g % 5) = 0 THEN FALSE ELSE TRUE END AS merged,
    (ARRAY[
        'ranjith-dev', 'hema-s', 'alex-k', 'priya-m', 'devops-bot',
        'ravi-p', 'sneha-r', 'amit-v', 'noor-f', 'li-wei',
        'maria-g', 'omar-h', 'yuki-t', 'ben-c', 'fatima-z',
        'vikram-s', 'ananya-k', 'carlos-m', 'diana-l', 'erik-n',
        'fiona-q', 'george-w', 'hana-y', 'ivan-p', 'julia-x'
    ])[( (g - 1) % 25 ) + 1] AS author_login,
    'feature/pr-' || g::text AS head_branch,
    'main' AS base_branch,
    (50 + (g % 800))::int AS additions,
    (10 + (g % 200))::int AS deletions,
    (1 + (g % 25))::int AS changed_files,
    TIMESTAMPTZ '2024-01-01 09:00:00+00' + (g * interval '6 hours') AS created_at,
    CASE
        WHEN (g % 7) = 0 THEN NULL
        WHEN (g % 5) = 0 THEN NULL
        ELSE TIMESTAMPTZ '2024-01-01 09:00:00+00' + (g * interval '6 hours') + interval '45 minutes'
    END AS merged_at,
    CASE
        WHEN (g % 7) = 0 THEN NULL
        ELSE TIMESTAMPTZ '2024-01-01 09:00:00+00' + (g * interval '6 hours') + interval '2 hours'
    END AS closed_at
FROM generate_series(1, 300) AS g;

-- ---------------------------------------------------------------------------
-- Seed: 220 workflow runs (linked to PRs; conclusions vary)
-- ---------------------------------------------------------------------------

INSERT INTO workflow_runs (
    repo_id, pr_id, workflow_name, event, status, conclusion,
    branch, run_attempt, duration_seconds, created_at, head_sha
)
SELECT
    pr.repo_id,
    pr.id AS pr_id,
    (ARRAY['CI', 'Build', 'Deploy', 'Lint', 'Test', 'E2E', 'Security scan'])[( (g - 1) % 7 ) + 1] AS workflow_name,
    (ARRAY['pull_request', 'push', 'schedule', 'workflow_dispatch', 'merge_group'])[( (g - 1) % 5 ) + 1] AS event,
    'completed' AS status,
    CASE (g % 11)
        WHEN 0 THEN 'failure'
        WHEN 1 THEN 'cancelled'
        WHEN 2 THEN 'skipped'
        ELSE 'success'
    END AS conclusion,
    pr.head_branch,
    (1 + (g % 3))::smallint AS run_attempt,
    (45 + (g % 1800))::int AS duration_seconds,
    pr.created_at + ((g % 200) || ' minutes')::interval AS created_at,
    substr(md5(g::text || pr.id::text), 1, 40) AS head_sha
FROM generate_series(1, 220) AS g
JOIN pull_requests pr ON pr.id = ((g * 17) % 300) + 1;

-- ---------------------------------------------------------------------------
-- Row counts (informational)
-- ---------------------------------------------------------------------------

DO $$
DECLARE
    u INT; r INT; p INT; w INT;
BEGIN
    SELECT count(*) INTO u FROM github_users;
    SELECT count(*) INTO r FROM repositories;
    SELECT count(*) INTO p FROM pull_requests;
    SELECT count(*) INTO w FROM workflow_runs;
    RAISE NOTICE 'Loaded github_users=%, repositories=%, pull_requests=%, workflow_runs=% (total %)',
        u, r, p, w, u + r + p + w;
END $$;
