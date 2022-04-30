ALTER TABLE IF EXISTS ONLY public.users_orgs DROP CONSTRAINT IF EXISTS users_orgs_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.users_orgs DROP CONSTRAINT IF EXISTS users_orgs_org_id_fkey;
DROP INDEX IF EXISTS public.users_email_key;
DROP INDEX IF EXISTS public.users_code_key;
DROP INDEX IF EXISTS public.userid_orgid_uniq;
DROP INDEX IF EXISTS public.orgs_code_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users_orgs DROP CONSTRAINT IF EXISTS users_orgs_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.orgs DROP CONSTRAINT IF EXISTS orgs_pkey;
ALTER TABLE IF EXISTS public.users_orgs ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.users_orgs_id_seq;
DROP TABLE IF EXISTS public.users_orgs;
DROP TABLE IF EXISTS public.users;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.schema_migrations;
DROP SEQUENCE IF EXISTS public.roles_id_seq;
DROP SEQUENCE IF EXISTS public.orgs_users_id_seq;
DROP TABLE IF EXISTS public.orgs;
DROP SEQUENCE IF EXISTS public.orgs_id_seq;

CREATE SEQUENCE public.orgs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.orgs (
    id integer DEFAULT nextval('public.orgs_id_seq'::regclass) NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    domain text,
    logo text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    status character varying(30) DEFAULT 'active'::character varying
);

CREATE SEQUENCE public.orgs_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.users (
    id integer DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL,
    code character varying(50) NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    email character varying(300) NOT NULL,
    phone character varying(25),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    status character varying(30) DEFAULT 'active'::character varying NOT NULL
);

CREATE TABLE public.users_orgs (
    id integer NOT NULL,
    user_id integer NOT NULL,
    org_id integer NOT NULL,
    role character varying(50),
    status character varying(30) DEFAULT 'active'::character varying NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE SEQUENCE public.users_orgs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
-- Dependencies: 216

ALTER SEQUENCE public.users_orgs_id_seq OWNED BY public.users_orgs.id;

ALTER TABLE ONLY public.users_orgs ALTER COLUMN id SET DEFAULT nextval('public.users_orgs_id_seq'::regclass);

ALTER TABLE ONLY public.orgs
    ADD CONSTRAINT orgs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

ALTER TABLE ONLY public.users_orgs
    ADD CONSTRAINT users_orgs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX orgs_code_key ON public.orgs USING btree (code);

CREATE UNIQUE INDEX userid_orgid_uniq ON public.users_orgs USING btree (user_id, org_id);

CREATE UNIQUE INDEX users_code_key ON public.users USING btree (code);

CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email);

ALTER TABLE ONLY public.users_orgs
    ADD CONSTRAINT users_orgs_org_id_fkey FOREIGN KEY (org_id) REFERENCES public.orgs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.users_orgs
    ADD CONSTRAINT users_orgs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
--
-- PostgreSQL database dump complete
--
