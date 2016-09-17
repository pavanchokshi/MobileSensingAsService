CREATE TABLE project_metadata (
project_id VARCHAR(32) NOT NULL,
project_name VARCHAR(32) NOT NULL,
project_desc VARCHAR(128),
project_owner VARCHAR(32) NOT NULL,
project_status VARCHAR(16) NOT NULL,
primary KEY (project_id));
