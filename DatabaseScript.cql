CREATE KEYSPACE tenant with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE tenant.tenant(tenant_id UUID, secret_key text, PRIMARY KEY(tenant_id));
