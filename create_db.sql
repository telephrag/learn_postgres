--
-- PostgreSQL database dump
--

-- Dumped from database version 12.7
-- Dumped by pg_dump version 12.7

-- Started on 2022-07-04 18:03:30 +04

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 219 (class 1255 OID 50016)
-- Name: stream_event_func(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.stream_event_func() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
    payload text;
begin
    --listen changestream;
 
    payload = new.diff::text;
    
    perform pg_notify(
      'changestream', 
      payload
    );
    
    return null;
end;
$$;


ALTER FUNCTION public.stream_event_func() OWNER TO postgres;

--
-- TOC entry 218 (class 1255 OID 49972)
-- Name: to_oplog_func(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.to_oplog_func() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
    payload jsonb;
    expire timestamp with time zone;
begin
    set timezone = 'UTC';

    expire = current_timestamp + interval '1 minute';
    payload = jsonb_build_object(
        'old',        row_to_json(old)::jsonb,
        'new',        row_to_json(new)::jsonb,
        'timestamp',  to_jsonb(current_timestamp),
        'table_name', to_jsonb(TG_TABLE_NAME),
        'optype',     to_jsonb(TG_OP)
    );
    
    insert into oplog(expire, diff) values(expire, payload);
    
    return null;
end;
$$;


ALTER FUNCTION public.to_oplog_func() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 203 (class 1259 OID 49992)
-- Name: interest; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.interest (
    id smallint NOT NULL,
    name character varying(64) NOT NULL
);


ALTER TABLE public.interest OWNER TO postgres;

--
-- TOC entry 202 (class 1259 OID 49990)
-- Name: interest_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.interest_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.interest_id_seq OWNER TO postgres;

--
-- TOC entry 3982 (class 0 OID 0)
-- Dependencies: 202
-- Name: interest_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.interest_id_seq OWNED BY public.interest.id;


--
-- TOC entry 205 (class 1259 OID 50001)
-- Name: oplog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.oplog (
    id smallint NOT NULL,
    diff jsonb NOT NULL,
    expire timestamp with time zone
);


ALTER TABLE public.oplog OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 49999)
-- Name: oplog_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.oplog_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.oplog_id_seq OWNER TO postgres;

--
-- TOC entry 3983 (class 0 OID 0)
-- Dependencies: 204
-- Name: oplog_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.oplog_id_seq OWNED BY public.oplog.id;


--
-- TOC entry 3839 (class 2604 OID 49995)
-- Name: interest id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.interest ALTER COLUMN id SET DEFAULT nextval('public.interest_id_seq'::regclass);


--
-- TOC entry 3840 (class 2604 OID 50004)
-- Name: oplog id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oplog ALTER COLUMN id SET DEFAULT nextval('public.oplog_id_seq'::regclass);


--
-- TOC entry 3842 (class 2606 OID 49997)
-- Name: interest interest_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.interest
    ADD CONSTRAINT interest_pkey PRIMARY KEY (id);


--
-- TOC entry 3844 (class 2606 OID 50009)
-- Name: oplog oplog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oplog
    ADD CONSTRAINT oplog_pkey PRIMARY KEY (id);


--
-- TOC entry 3846 (class 2620 OID 50018)
-- Name: oplog stream_event; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER stream_event AFTER INSERT ON public.oplog FOR EACH ROW EXECUTE FUNCTION public.stream_event_func();


--
-- TOC entry 3845 (class 2620 OID 49998)
-- Name: interest to_oplog; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER to_oplog AFTER INSERT OR DELETE OR UPDATE ON public.interest FOR EACH ROW EXECUTE FUNCTION public.to_oplog_func();


-- Completed on 2022-07-04 18:03:30 +04

--
-- PostgreSQL database dump complete
--

