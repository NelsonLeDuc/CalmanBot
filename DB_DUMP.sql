--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.1
-- Dumped by pg_dump version 9.5.1

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
-- Name: actions; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE actions (
    id integer NOT NULL,
    type character varying NOT NULL,
    content character varying NOT NULL,
    data_path character varying,
    pattern character varying,
    main boolean NOT NULL,
    priority integer NOT NULL,
    fallback_action integer,
    post_text character varying
);


ALTER TABLE actions OWNER TO nelsonleduc;

--
-- Name: actions_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE actions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE actions_id_seq OWNER TO nelsonleduc;

--
-- Name: actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE actions_id_seq OWNED BY actions.id;


--
-- Name: bots; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE bots (
    group_id character varying NOT NULL,
    group_name character varying NOT NULL,
    bot_name character varying NOT NULL,
    key character varying NOT NULL
);


ALTER TABLE bots OWNER TO nelsonleduc;

--
-- Name: bots_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE bots_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE bots_id_seq OWNER TO nelsonleduc;

--
-- Name: bots_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE bots_id_seq OWNED BY bots.group_id;


--
-- Name: cached; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE cached (
    id integer NOT NULL,
    query text,
    result text
);


ALTER TABLE cached OWNER TO nelsonleduc;

--
-- Name: cached_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE cached_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cached_id_seq OWNER TO nelsonleduc;

--
-- Name: cached_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE cached_id_seq OWNED BY cached.id;


--
-- Name: groupme_posts; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE groupme_posts (
    id integer NOT NULL,
    cache_id integer,
    likes integer DEFAULT 0,
    message_id character varying,
    group_id character varying,
    posted_at timestamp without time zone DEFAULT now()
);


ALTER TABLE groupme_posts OWNER TO nelsonleduc;

--
-- Name: groupme_posts_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE groupme_posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE groupme_posts_id_seq OWNER TO nelsonleduc;

--
-- Name: groupme_posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE groupme_posts_id_seq OWNED BY groupme_posts.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY actions ALTER COLUMN id SET DEFAULT nextval('actions_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY cached ALTER COLUMN id SET DEFAULT nextval('cached_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY groupme_posts ALTER COLUMN id SET DEFAULT nextval('groupme_posts_id_seq'::regclass);


--
-- Data for Name: actions; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY actions (id, type, content, data_path, pattern, main, priority, fallback_action, post_text) FROM stdin;
11	TEXT	Lol 12	\N	\N	f	100	\N	\N
8	URL/IMAGE	http://www.reddit.com/r/calmangonewild/.json	data.children.{_randomInt_}.data.url	[@]{_botname_} (c(al|law)man)	t	1	12	\N
7	TEXT	Do I look like a map?	\N	[@]{_botname_} (where (is|am|are|be))	t	0	\N	\N
12	URL/IMAGE	http://api.giphy.com/v1/gifs/search?q={_text_}&api_key=dc6zaTOxFJmzC&limit=40	data.{_randomInt_}.images.original.url	[@]{_botname_} (.*)	t	20	10	\N
2	TEXT	http://crossfitsouthcobb.com/wp-content/uploads/2014/10/cookie_monster_original.jpg	\N	[&]{_botname_} (coo+kies*)	t	0	\N	\N
1	URL/TEXT	http://www.quandl.com/api/v1/datasets/CHRIS/CME_BZ1.json	data.0.6	[&]{_botname_} (oil)	t	2	\N	\N
9	URL/TEXT	http://utdeats.com/university-eats/json/dining.php?id=1&type=off	Off.{_randomInt_}.Name	[@]{_botname_} (where should we eat\\?|foods\\?)	t	2	12	\N
13	URL/TEXT	https://www.googleapis.com/youtube/v3/search?key={_key(yt)_}&part=snippet&type=video&q={_text_}	items.{_randomInt_}.id.videoId	[@]{_botname_} youtube (.*)	t	3	15	https://www.youtube.com/watch?v={_text_}
14	URL/TEXT	https://www.googleapis.com/youtube/v3/search?key={_key(yt)_}&part=snippet&maxResults=1&type=video&q={_text_}	items.{_randomInt_}.id.videoId	[@]{_botname_} \\$youtube (.*)	t	4	15	https://www.youtube.com/watch?v={_text_}
10	URL/IMAGE	http://calmanbot-production.herokuapp.com/animated?q={_text_}	{_randomInt_}	[@]{_botname_} (.*)	f	100	15	\N
15	URL/TEXT	https://www.reddit.com/r/copypasta/top.json?sort=top&t=all&limit=100	data.children.{_randomInt_}.data.selftext	\N	f	100	11	\N
3	URL/IMAGE	http://ajax.googleapis.com/ajax/services/search/images?v=1.0&as_filetype=png&rsz=8&q={_text_}	responseData.results.{_randomInt_}.url	[@]{_botname_} im(?:(?:age)|g) (.*)	f	18	15	\N
5	URL/TEXT	http://jsoncat.parseapp.com/numMessages?groupID={groupID}	messageCount	[&]{_botname_} (num$|number$)	f	1	\N	\N
\.


--
-- Name: actions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('actions_id_seq', 15, true);


--
-- Data for Name: bots; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY bots (group_id, group_name, bot_name, key) FROM stdin;
9214876	Test	CalmanBot	PLACEHOLDER
9288084	Games	CalmanBot	PLACEHOLDER
12515792	SpringBreak	CalmanBot	PLACEHOLDER
10866751	Archeage	CalmanBot	PLACEHOLDER
9197483	Food	CalmanBot	PLACEHOLDER
7903597	BR 	CalmanBot	PLACEHOLDER
4067479	JS	CalmanBot	PLACEHOLDER
5785582	Main	CalmanBot	PLACEHOLDER
\.


--
-- Name: bots_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('bots_id_seq', 3, true);


--
-- Data for Name: cached; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY cached (id, query, result) FROM stdin;
\.


--
-- Name: cached_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('cached_id_seq', 289, true);


--
-- Data for Name: groupme_posts; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY groupme_posts (id, cache_id, likes, message_id, group_id, posted_at) FROM stdin;
\.


--
-- Name: groupme_posts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('groupme_posts_id_seq', 292, true);


--
-- Name: actions_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY actions
    ADD CONSTRAINT actions_pkey PRIMARY KEY (id);


--
-- Name: bots_group_id_key; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY bots
    ADD CONSTRAINT bots_group_id_key UNIQUE (group_id);


--
-- Name: bots_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY bots
    ADD CONSTRAINT bots_pkey PRIMARY KEY (group_id);


--
-- Name: cached_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY cached
    ADD CONSTRAINT cached_pkey PRIMARY KEY (id);


--
-- Name: groupme_posts_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY groupme_posts
    ADD CONSTRAINT groupme_posts_pkey PRIMARY KEY (id);


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

