/* business_code */
CREATE TABLE business_code (
    id integer NOT NULL,
    code character varying,
    description character varying
);

CREATE SEQUENCE business_code_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE business_code_id_seq OWNED BY business_code.id;

ALTER TABLE ONLY business_code ALTER COLUMN id SET DEFAULT nextval('business_code_id_seq'::regclass);

ALTER TABLE ONLY business_code
    ADD CONSTRAINT business_code_code_description UNIQUE (code, description);

ALTER TABLE ONLY business_code
    ADD CONSTRAINT business_code_pkey PRIMARY KEY (id);

/* sa3 */
CREATE TABLE sa3 (
    id integer NOT NULL,
    code integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE sa3_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sa3_id_seq OWNED BY sa3.id;

ALTER TABLE ONLY sa3 ALTER COLUMN id SET DEFAULT nextval('sa3_id_seq'::regclass);

ALTER TABLE ONLY sa3
    ADD CONSTRAINT sa3_code_name_key UNIQUE (code, name);

ALTER TABLE ONLY sa3
    ADD CONSTRAINT sa3_pkey PRIMARY KEY (id);

/* sa4 */
CREATE TABLE sa4 (
    id integer NOT NULL,
    code integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE sa4_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE sa4_id_seq OWNED BY sa4.id;

ALTER TABLE ONLY sa4 ALTER COLUMN id SET DEFAULT nextval('sa4_id_seq'::regclass);

ALTER TABLE ONLY sa4
    ADD CONSTRAINT sa4_code_name_key UNIQUE (code, name);

ALTER TABLE ONLY sa4
    ADD CONSTRAINT sa4_pkey PRIMARY KEY (id);

/* state */
CREATE TABLE state (
    id integer NOT NULL,
    name character varying NOT NULL
);

CREATE SEQUENCE state_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE state_id_seq OWNED BY state.id;

ALTER TABLE ONLY state ALTER COLUMN id SET DEFAULT nextval('state_id_seq'::regclass);

ALTER TABLE ONLY state
    ADD CONSTRAINT state_name_key UNIQUE (name);

ALTER TABLE ONLY state
    ADD CONSTRAINT state_pkey PRIMARY KEY (id);

/* postcode */
CREATE TABLE postcode (
    id integer NOT NULL,
    postcode character varying NOT NULL,
    locality character varying NOT NULL,
    state_id integer NOT NULL,
    sa3_id integer,
    sa4_id integer
);

CREATE SEQUENCE postcode_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE postcode_id_seq OWNED BY postcode.id;

ALTER TABLE ONLY postcode ALTER COLUMN id SET DEFAULT nextval('postcode_id_seq'::regclass);

ALTER TABLE ONLY postcode
    ADD CONSTRAINT postcode_locality_state UNIQUE (postcode, locality, state_id);

ALTER TABLE ONLY postcode
    ADD CONSTRAINT postcode_pkey PRIMARY KEY (id);

ALTER TABLE ONLY postcode
    ADD CONSTRAINT postcode_sa3_id_fkey FOREIGN KEY (sa3_id) REFERENCES sa3(id);

ALTER TABLE ONLY postcode
    ADD CONSTRAINT postcode_sa4_id_fkey FOREIGN KEY (sa4_id) REFERENCES sa4(id);

ALTER TABLE ONLY postcode
    ADD CONSTRAINT postcode_state_id_fkey FOREIGN KEY (state_id) REFERENCES state(id);

/* organisation */
CREATE TABLE organisation (
    id INTEGER NOT NULL,
    organisation_source_id INTEGER NOT NULL,
    anzsic_id INTEGER,
    business_code_id INTEGER REFERENCES business_code(id) NOT NULL,
    postcode_id INTEGER REFERENCES postcode(id) NOT NULL
);

CREATE SEQUENCE organisation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organisation_id_seq OWNED BY organisation.id;

ALTER TABLE ONLY organisation 
    ALTER COLUMN id SET DEFAULT nextval('organisation_id_seq'::regclass);

ALTER TABLE ONLY organisation 
    ADD CONSTRAINT organisation_pkey PRIMARY KEY (id);

ALTER TABLE ONLY organisation
    ADD CONSTRAINT organisation_organisation_source_id UNIQUE (organisation_source_id);

/* ansiz */
CREATE TABLE anzsic (
    id integer NOT NULL,
    code character varying,
    description character varying
);

CREATE SEQUENCE anzsic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE anzsic_id_seq OWNED BY anzsic.id;

ALTER TABLE ONLY anzsic ALTER COLUMN id SET DEFAULT nextval('anzsic_id_seq'::regclass);

ALTER TABLE ONLY anzsic
    ADD CONSTRAINT anzsic_pkey PRIMARY KEY (id);

ALTER TABLE ONLY anzsic
    ADD CONSTRAINT code_description UNIQUE (code, description);

/* email */
CREATE TABLE email (
    id integer NOT NULL,
    email CHARACTER VARYING NOT NULL
);

CREATE SEQUENCE email_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE email_id_seq OWNED BY email.id;

ALTER TABLE ONLY email ALTER COLUMN id SET DEFAULT nextval('email_id_seq'::regclass);

ALTER TABLE ONLY email
    ADD CONSTRAINT email_pkey PRIMARY KEY (id);

ALTER TABLE ONLY email
    ADD CONSTRAINT email_unique UNIQUE (email);

/* email_organisation */
CREATE TABLE email_organisation (
    id integer NOT NULL,
    email_id INTEGER REFERENCES email(id) NOT NULL,
    organisation_id INTEGER REFERENCES organisation(id) NOT NULL
);

CREATE SEQUENCE email_organisation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE email_organisation_id_seq OWNED BY email_organisation.id;

ALTER TABLE ONLY email_organisation ALTER COLUMN id SET DEFAULT nextval('email_organisation_id_seq'::regclass);

ALTER TABLE ONLY email_organisation
    ADD CONSTRAINT email_organisation_pkey PRIMARY KEY (id);

ALTER TABLE ONLY email_organisation
    ADD CONSTRAINT email_organisation_unique UNIQUE (email_id, organisation_id);

/* DEBUG:[]string{"EMAIL-2", "WEBSITE", "FACEBOOK", "TWITTER", "LINKEDIN", "EMPLOYEES", "REVENUE-$M", "YEAR-ESTABLISHED", "CONTACT-NAME", "CONTACT-FIRST-NAME", "CONTACT-JOB-TITLE", "ABN", "ABN-STATUS", "STATUS-DATE", "ENTITY-TYPE-CODE", "ANZSIC-CODE", "LATITUDE", "LONGITUDE", "MAPLINK", "ID-ORG"}              */

/* organisation_state_vic */
CREATE TABLE organisation_state_vic (
    id integer NOT NULL,
    organisation_id INTEGER REFERENCES organisation(id) NOT NULL,
    name character varying NOT NULL,
    abn character varying,
    address character varying,
    record_defunct_risk CHARACTER VARYING,
    region CHARACTER VARYING,
    phone CHARACTER VARYING,
    mobile CHARACTER VARYING,
    freecall CHARACTER VARYING,
    fax CHARACTER VARYING
);

CREATE SEQUENCE organisation_state_vic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organisation_state_vic_id_seq OWNED BY organisation_state_vic.id;

ALTER TABLE ONLY organisation_state_vic 
    ALTER COLUMN id SET DEFAULT nextval('organisation_state_vic_id_seq'::regclass);

ALTER TABLE ONLY organisation_state_vic 
    ADD CONSTRAINT organisation_state_vic_pkey PRIMARY KEY (id);

ALTER TABLE ONLY organisation_state_vic
    ADD CONSTRAINT organisation_state_vic_name_address_region_unique UNIQUE (name, address, region);

ALTER TABLE ONLY organisation_state_vic
    ADD CONSTRAINT organisation_state_vic_organisation_id_unique UNIQUE (organisation_id);

CREATE INDEX organistaion_state_vic_trgm_idx ON organisation_state_vic USING GIST (name gist_trgm_ops(siglen=32));

/* organisation_state_nsw */
CREATE TABLE organisation_state_nsw (
    id integer NOT NULL,
    organisation_id INTEGER REFERENCES organisation(id) NOT NULL,
    name character varying NOT NULL,
    abn character varying,
    address character varying,
    record_defunct_risk CHARACTER VARYING,
    region CHARACTER VARYING,
    phone CHARACTER VARYING,
    mobile CHARACTER VARYING,
    freecall CHARACTER VARYING,
    fax CHARACTER VARYING
);

CREATE SEQUENCE organisation_state_nsw_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organisation_state_nsw_id_seq OWNED BY organisation_state_nsw.id;

ALTER TABLE ONLY organisation_state_nsw 
    ALTER COLUMN id SET DEFAULT nextval('organisation_state_nsw_id_seq'::regclass);

ALTER TABLE ONLY organisation_state_nsw 
    ADD CONSTRAINT organisation_state_nsw_pkey PRIMARY KEY (id);

ALTER TABLE ONLY organisation_state_nsw
    ADD CONSTRAINT organisation_state_nsw_name_address_region_unique UNIQUE (name, address, region);

ALTER TABLE ONLY organisation_state_nsw
    ADD CONSTRAINT organisation_state_nsw_organisation_id_unique UNIQUE (organisation_id);

CREATE INDEX organistaion_state_nsw_trgm_idx ON organisation_state_nsw USING GIST (name gist_trgm_ops(siglen=32));
