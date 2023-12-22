/* ansiz */
CREATE TABLE public.anzsic (
    id integer NOT NULL,
    code character varying,
    description character varying
);

CREATE SEQUENCE public.anzsic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.anzsic_id_seq OWNED BY public.anzsic.id;

ALTER TABLE ONLY public.anzsic ALTER COLUMN id SET DEFAULT nextval('public.anzsic_id_seq'::regclass);

ALTER TABLE ONLY public.anzsic
    ADD CONSTRAINT anzsic_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.anzsic
    ADD CONSTRAINT code_description UNIQUE (code, description);

/* business_code */
CREATE TABLE public.business_code (
    id integer NOT NULL,
    code character varying,
    description character varying
);

CREATE SEQUENCE public.business_code_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.business_code_id_seq OWNED BY public.business_code.id;

ALTER TABLE ONLY public.business_code ALTER COLUMN id SET DEFAULT nextval('public.business_code_id_seq'::regclass);

ALTER TABLE ONLY public.business_code
    ADD CONSTRAINT business_code_code_description UNIQUE (code, description);

ALTER TABLE ONLY public.business_code
    ADD CONSTRAINT business_code_pkey PRIMARY KEY (id);

/* sa3 */
CREATE TABLE public.sa3 (
    id integer NOT NULL,
    code integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE public.sa3_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.sa3_id_seq OWNED BY public.sa3.id;

ALTER TABLE ONLY public.sa3 ALTER COLUMN id SET DEFAULT nextval('public.sa3_id_seq'::regclass);

ALTER TABLE ONLY public.sa3
    ADD CONSTRAINT sa3_code_name_key UNIQUE (code, name);

ALTER TABLE ONLY public.sa3
    ADD CONSTRAINT sa3_pkey PRIMARY KEY (id);

/* sa4 */
CREATE TABLE public.sa4 (
    id integer NOT NULL,
    code integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE public.sa4_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.sa4_id_seq OWNED BY public.sa4.id;

ALTER TABLE ONLY public.sa4 ALTER COLUMN id SET DEFAULT nextval('public.sa4_id_seq'::regclass);

ALTER TABLE ONLY public.sa4
    ADD CONSTRAINT sa4_code_name_key UNIQUE (code, name);

ALTER TABLE ONLY public.sa4
    ADD CONSTRAINT sa4_pkey PRIMARY KEY (id);

/* state */
CREATE TABLE public.state (
    id integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE public.state_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.state_id_seq OWNED BY public.state.id;

ALTER TABLE ONLY public.state ALTER COLUMN id SET DEFAULT nextval('public.state_id_seq'::regclass);

ALTER TABLE ONLY public.state
    ADD CONSTRAINT state_name_key UNIQUE (name);

ALTER TABLE ONLY public.state
    ADD CONSTRAINT state_pkey PRIMARY KEY (id);

/* postcode */
CREATE TABLE public.postcode (
    id integer NOT NULL,
    postcode character varying NOT NULL,
    locality character varying NOT NULL,
    state_id integer NOT NULL,
    sa3_id integer,
    sa4_id integer
);

CREATE SEQUENCE public.postcode_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.postcode_id_seq OWNED BY public.postcode.id;

ALTER TABLE ONLY public.postcode ALTER COLUMN id SET DEFAULT nextval('public.postcode_id_seq'::regclass);

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_locality UNIQUE (postcode, locality);

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_sa3_id_fkey FOREIGN KEY (sa3_id) REFERENCES public.sa3(id);

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_sa4_id_fkey FOREIGN KEY (sa4_id) REFERENCES public.sa4(id);

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_state_id_fkey FOREIGN KEY (state_id) REFERENCES public.state(id);

/* organisation_source */
CREATE TABLE public.organisation_source (
    id integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE public.organisation_source_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.organisation_source_id_seq OWNED BY public.organisation_source.id;

ALTER TABLE ONLY public.organisation_source ALTER COLUMN id SET DEFAULT nextval('public.organisation_source_id_seq'::regclass);

ALTER TABLE ONLY public.organisation_source
    ADD CONSTRAINT organisation_source_name_key UNIQUE (name);

ALTER TABLE ONLY public.organisation_source
    ADD CONSTRAINT organisation_source_pkey PRIMARY KEY (id);

/* organisation */
CREATE TABLE public.organisation (
    id integer NOT NULL,
    name character varying NOT NULL,
    abn character varying,
    address character varying,
    anzsic_id integer,
    business_code_id integer NOT NULL,
    postcode_id integer,
    source_id integer,
    organisation_source_id integer NOT NULL
);

CREATE SEQUENCE public.organisation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.organisation_id_seq OWNED BY public.organisation.id;

ALTER TABLE ONLY public.organisation ALTER COLUMN id SET DEFAULT nextval('public.organisation_id_seq'::regclass);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_source_id_organisation_source_id UNIQUE (source_id, organisation_source_id);

CREATE INDEX organisation_name ON public.organisation USING btree (name);

CREATE INDEX organisation_name_postcode_id ON public.organisation USING btree (name, postcode_id);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_anzsic_id_fkey FOREIGN KEY (anzsic_id) REFERENCES public.anzsic(id);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_business_code_id_fkey FOREIGN KEY (business_code_id) REFERENCES public.business_code(id);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_organisation_source_id_fkey FOREIGN KEY (organisation_source_id) REFERENCES public.organisation_source(id);

ALTER TABLE ONLY public.organisation
    ADD CONSTRAINT organisation_postcode_id_fkey FOREIGN KEY (postcode_id) REFERENCES public.postcode(id);

/* transaction */

CREATE TABLE public.transaction (
    id integer NOT NULL,
    memo character varying NOT NULL
);

CREATE SEQUENCE public.transaction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.transaction_id_seq OWNED BY public.transaction.id;

ALTER TABLE ONLY public.transaction ALTER COLUMN id SET DEFAULT nextval('public.transaction_id_seq'::regclass);

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_pkey PRIMARY KEY (id);

/* alembic */
CREATE TABLE public.alembic_version (
    version_num character varying(32) NOT NULL
);

ALTER TABLE ONLY public.alembic_version
    ADD CONSTRAINT alembic_version_pkc PRIMARY KEY (version_num);
