--Roles
-- super_admin = Read/Write everything in all applications
-- integrand_admin = Read/Write in Integrand
-- integrand_user = Read only in Integrand

-- Integrand Securables
--- Topic
    -- Read_Topic
    -- Create_Topic
    -- Edit_Topic
    -- Delete_Topic
--- Endpoint
    -- Read_Endpoint
    -- Create_Endpoint
    -- Edit_Endpoint
    -- Delete_Endpoint
--- Workflow
    -- Read_Workflow
    -- Create_Workflow
    -- Edit_Workflow
    -- Delete_Workflow
--- API Key
    -- Read_API_Key
    -- Create_API_Key
    -- Delete
--- User
    -- Read_User
    -- Create_User
    -- Edit_User
    -- Delete_User

CREATE TABLE IF NOT EXISTS roles (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS securables (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_to_role (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(role_id) REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_to_securable (
    role_id INTEGER NOT NULL,
    securable_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(role_id) REFERENCES roles(id),
    FOREIGN KEY(securable_id) REFERENCES securables(id),
    PRIMARY KEY (role_id, securable_id)
);