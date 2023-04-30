CREATE SCHEMA IF NOT EXISTS public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.example_table (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar(100) NOT NULL,
	CONSTRAINT example_table_pk PRIMARY KEY (id),
	CONSTRAINT example_table_un UNIQUE ("name")
);