--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: builds; Type: TABLE; Schema: public; Owner: quodeck
--

CREATE TABLE builds (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE builds OWNER TO quodeck;

--
-- Name: phases; Type: TABLE; Schema: public; Owner: quodeck
--

CREATE TABLE phases (
    id uuid NOT NULL,
    data jsonb NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE phases OWNER TO quodeck;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: quodeck
--

CREATE TABLE schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE schema_migration OWNER TO quodeck;

--
-- Name: users; Type: TABLE; Schema: public; Owner: quodeck
--

CREATE TABLE users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    tokens jsonb
);


ALTER TABLE users OWNER TO quodeck;

--
-- Name: builds_pkey; Type: CONSTRAINT; Schema: public; Owner: quodeck
--

ALTER TABLE ONLY builds
    ADD CONSTRAINT builds_pkey PRIMARY KEY (id);


--
-- Name: phases_pkey; Type: CONSTRAINT; Schema: public; Owner: quodeck
--

ALTER TABLE ONLY phases
    ADD CONSTRAINT phases_pkey PRIMARY KEY (id);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: quodeck
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: quodeck
--

CREATE UNIQUE INDEX schema_migration_version_idx ON schema_migration USING btree (version);


--
-- Name: builds_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: quodeck
--

ALTER TABLE ONLY builds
    ADD CONSTRAINT builds_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- Name: phases_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: quodeck
--

ALTER TABLE ONLY phases
    ADD CONSTRAINT phases_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

