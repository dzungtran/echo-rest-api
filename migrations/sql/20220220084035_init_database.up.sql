ALTER TABLE IF EXISTS ONLY users_orgs DROP CONSTRAINT IF EXISTS users_orgs_user_id_fkey;
ALTER TABLE IF EXISTS ONLY users_orgs DROP CONSTRAINT IF EXISTS users_orgs_org_id_fkey;
DROP INDEX IF EXISTS users_email_key;
DROP INDEX IF EXISTS users_code_key;
DROP INDEX IF EXISTS userid_orgid_uniq;
DROP INDEX IF EXISTS orgs_code_key;
ALTER TABLE IF EXISTS ONLY users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY users_orgs DROP CONSTRAINT IF EXISTS users_orgs_pkey;
ALTER TABLE IF EXISTS ONLY schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY orgs DROP CONSTRAINT IF EXISTS orgs_pkey;
ALTER TABLE IF EXISTS users_orgs ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS users_orgs_id_seq;
DROP TABLE IF EXISTS users_orgs;
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS users_id_seq;
DROP TABLE IF EXISTS schema_migrations;
DROP SEQUENCE IF EXISTS roles_id_seq;
DROP SEQUENCE IF EXISTS orgs_users_id_seq;
DROP TABLE IF EXISTS orgs;
DROP SEQUENCE IF EXISTS orgs_id_seq;

CREATE SEQUENCE orgs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE orgs (
    id integer DEFAULT nextval('orgs_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    domain text,
    logo text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    status character varying(30) DEFAULT 'active'::character varying
);

CREATE SEQUENCE orgs_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE users (
    id integer DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
    code character varying(50) NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    email character varying(300) NOT NULL,
    phone character varying(25),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    status character varying(30) DEFAULT 'active'::character varying NOT NULL
);

CREATE TABLE users_orgs (
    id integer NOT NULL,
    user_id integer NOT NULL,
    org_id integer NOT NULL,
    role character varying(50),
    status character varying(30) DEFAULT 'active'::character varying NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE SEQUENCE users_orgs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
-- Dependencies: 216

ALTER SEQUENCE users_orgs_id_seq OWNED BY users_orgs.id;

ALTER TABLE ONLY users_orgs ALTER COLUMN id SET DEFAULT nextval('users_orgs_id_seq'::regclass);

ALTER TABLE ONLY orgs
    ADD CONSTRAINT orgs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

ALTER TABLE ONLY users_orgs
    ADD CONSTRAINT users_orgs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX orgs_code_key ON orgs USING btree (code);

CREATE UNIQUE INDEX userid_orgid_uniq ON users_orgs USING btree (user_id, org_id);

CREATE UNIQUE INDEX users_code_key ON users USING btree (code);

CREATE UNIQUE INDEX users_email_key ON users USING btree (email);

ALTER TABLE ONLY users_orgs
    ADD CONSTRAINT users_orgs_org_id_fkey FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE;

ALTER TABLE ONLY users_orgs
    ADD CONSTRAINT users_orgs_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
--
-- PostgreSQL database dump complete
--
