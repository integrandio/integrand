-- Add our roles
INSERT INTO roles(name, description) VALUES('super_admin', 'the super admin'),
    ('integrand_admin', 'the regular admin'),
    ('integrand_user', 'the regular user');

-- Add our securables
INSERT INTO securables(name) VALUES('read_topic'), ('write_topic'),
    ('read_endpoint'), ('write_endpoint'),
    ('read_workflow'), ('write_workflow'),
    ('read_user'), ('write_user'),
    ('read_api_key'), ('write_api_key');

-- Create our roles_securables
INSERT INTO role_to_securable(role_id, securable_id) SELECT roles.id, securables.id FROM roles INNER JOIN securables ON roles.name = 'super_admin';

INSERT INTO role_to_securable(role_id, securable_id) 
    SELECT roles.id, securables.id FROM roles INNER JOIN securables ON
    roles.name = 'integrand_admin' AND 
      (securables.name = 'read_topic' OR securables.name = 'write_topic' 
      OR securables.name = 'read_endpoint' OR securables.name = 'write_endpoint'
      OR securables.name = 'read_workflow' OR securables.name = 'write_workflow' 
      OR securables.name = 'read_api_key' OR securables.name = 'read_user');

INSERT INTO role_to_securable(role_id, securable_id) 
    SELECT roles.id, securables.id FROM roles INNER JOIN securables ON
    roles.name = 'integrand_user' AND
      (securables.name = 'read_topic' OR securables.name = 'read_endpoint'
      OR securables.name = 'read_workflow' OR securables.name = 'read_user');