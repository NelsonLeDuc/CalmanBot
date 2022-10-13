--
-- PostgreSQL database dump
--

-- Dumped from database version 11.17 (Ubuntu 11.17-1.pgdg20.04+1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-1.pgdg22.04+1)

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
-- Name: heroku_ext; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA heroku_ext;


ALTER SCHEMA heroku_ext OWNER TO postgres;

SET default_tablespace = '';

--
-- Name: actions; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.actions (
    id integer NOT NULL,
    type character varying NOT NULL,
    content character varying NOT NULL,
    data_path character varying,
    pattern character varying,
    main boolean NOT NULL,
    priority integer NOT NULL,
    fallback_action integer,
    post_text character varying,
    description text,
    note_process_mode integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.actions OWNER TO nelsonleduc;

--
-- Name: actions_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.actions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.actions_id_seq OWNER TO nelsonleduc;

--
-- Name: actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.actions_id_seq OWNED BY public.actions.id;


--
-- Name: bots; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.bots (
    group_id character varying NOT NULL,
    group_name character varying NOT NULL,
    bot_name character varying NOT NULL,
    key character varying NOT NULL
);


ALTER TABLE public.bots OWNER TO nelsonleduc;

--
-- Name: bots_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.bots_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.bots_id_seq OWNER TO nelsonleduc;

--
-- Name: bots_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.bots_id_seq OWNED BY public.bots.group_id;


--
-- Name: cached; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.cached (
    id integer NOT NULL,
    query text,
    result text
);


ALTER TABLE public.cached OWNER TO nelsonleduc;

--
-- Name: cached_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.cached_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cached_id_seq OWNER TO nelsonleduc;

--
-- Name: cached_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.cached_id_seq OWNED BY public.cached.id;


--
-- Name: discord_status; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.discord_status (
    id integer NOT NULL,
    type integer NOT NULL,
    text character varying NOT NULL
);


ALTER TABLE public.discord_status OWNER TO nelsonleduc;

--
-- Name: discord_status_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.discord_status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.discord_status_id_seq OWNER TO nelsonleduc;

--
-- Name: discord_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.discord_status_id_seq OWNED BY public.discord_status.id;


--
-- Name: discord_triggers; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.discord_triggers (
    id integer NOT NULL,
    channel_id character varying NOT NULL,
    trigger_id character varying NOT NULL,
    server_id character varying
);


ALTER TABLE public.discord_triggers OWNER TO nelsonleduc;

--
-- Name: discord_triggers_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.discord_triggers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.discord_triggers_id_seq OWNER TO nelsonleduc;

--
-- Name: discord_triggers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.discord_triggers_id_seq OWNED BY public.discord_triggers.id;


--
-- Name: groupme_posts; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.groupme_posts (
    id integer NOT NULL,
    cache_id integer,
    likes integer DEFAULT 0,
    message_id character varying,
    group_id character varying,
    posted_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.groupme_posts OWNER TO nelsonleduc;

--
-- Name: groupme_posts_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.groupme_posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.groupme_posts_id_seq OWNER TO nelsonleduc;

--
-- Name: groupme_posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.groupme_posts_id_seq OWNED BY public.groupme_posts.id;


--
-- Name: minecraft_servers; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.minecraft_servers (
    id integer NOT NULL,
    address character varying NOT NULL,
    name character varying
);


ALTER TABLE public.minecraft_servers OWNER TO nelsonleduc;

--
-- Name: minecraft_servers_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.minecraft_servers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.minecraft_servers_id_seq OWNER TO nelsonleduc;

--
-- Name: minecraft_servers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.minecraft_servers_id_seq OWNED BY public.minecraft_servers.id;


--
-- Name: spotify_playlists; Type: TABLE; Schema: public; Owner: nelsonleduc
--

CREATE TABLE public.spotify_playlists (
    id integer NOT NULL,
    group_id character varying NOT NULL,
    playlist_id character varying NOT NULL
);


ALTER TABLE public.spotify_playlists OWNER TO nelsonleduc;

--
-- Name: spotify_playlists_id_seq; Type: SEQUENCE; Schema: public; Owner: nelsonleduc
--

CREATE SEQUENCE public.spotify_playlists_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.spotify_playlists_id_seq OWNER TO nelsonleduc;

--
-- Name: spotify_playlists_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: nelsonleduc
--

ALTER SEQUENCE public.spotify_playlists_id_seq OWNED BY public.spotify_playlists.id;


--
-- Name: actions id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.actions ALTER COLUMN id SET DEFAULT nextval('public.actions_id_seq'::regclass);


--
-- Name: cached id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.cached ALTER COLUMN id SET DEFAULT nextval('public.cached_id_seq'::regclass);


--
-- Name: discord_status id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.discord_status ALTER COLUMN id SET DEFAULT nextval('public.discord_status_id_seq'::regclass);


--
-- Name: discord_triggers id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.discord_triggers ALTER COLUMN id SET DEFAULT nextval('public.discord_triggers_id_seq'::regclass);


--
-- Name: groupme_posts id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.groupme_posts ALTER COLUMN id SET DEFAULT nextval('public.groupme_posts_id_seq'::regclass);


--
-- Name: minecraft_servers id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.minecraft_servers ALTER COLUMN id SET DEFAULT nextval('public.minecraft_servers_id_seq'::regclass);


--
-- Name: spotify_playlists id; Type: DEFAULT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.spotify_playlists ALTER COLUMN id SET DEFAULT nextval('public.spotify_playlists_id_seq'::regclass);


--
-- Data for Name: actions; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.actions (id, type, content, data_path, pattern, main, priority, fallback_action, post_text, description, note_process_mode) FROM stdin;
11	TEXT	Lol 12	\N	\N	f	100	\N	\N	\N	0
7	TEXT	Do I look like a map?	\N	[@]{_botname_} (where (is|am|are|be))	t	0	\N	\N	\N	0
3	URL/IMAGE	http://ajax.googleapis.com/ajax/services/search/images?v=1.0&as_filetype=png&rsz=8&q={_text_}	responseData.results.{_randomInt_}.url	[@]{_botname_} im(?:(?:age)|g) (.*)	f	18	15	\N	\N	0
5	URL/TEXT	http://jsoncat.parseapp.com/numMessages?groupID={groupID}	messageCount	[&]{_botname_} (num$|number$)	f	1	\N	\N	\N	0
8	URL/IMAGE	http://www.reddit.com/r/calmangonewild/.json	data.children.{_randomInt_}.data.url	[@]{_botname_} (c(al|law)man)	t	1	12	\N	Calman	0
9	URL/TEXT	http://utdeats.com/university-eats/json/dining.php?id=1&type=off	Off.{_randomInt_}.Name	[@]{_botname_} (where should we eat\\?|foods\\?)	t	2	12	\N	Food near UTD	0
2	TEXT	https://media.giphy.com/media/whNK1SAMSQjwQ/giphy.gif	\N	[@]{_botname_} !(coo+kies*|🍪)	t	0	\N	\N	Cookie Monster	0
1	URL/TEXT	http://www.quandl.com/api/v1/datasets/CHRIS/CME_BZ1.json	data.0.6	[@]{_botname_} !(oil)	t	2	\N	\N	Current oil price	0
16	URL/TEXT	http://{_me_}/jeff?n=3	output	[@]{_botname_} !(jeff.*)	t	0	11	\N	Some Jeff	0
17	URL/TEXT	http://{_me_}/minecraftStatus?addr=m.mwe.st%3A25565	description	[@]{_botname_} (is the server.*)	t	1	11	\N	\N	0
10	URL/IMAGE	http://{_me_}/animated?q={_text_}	{_randomInt_}	[@]{_botname_} (.*)	f	100	15	\N	\N	0
41	URL/TEXT	http://{_me_}/playlistGet?groupid={_groupid_}&groupName={_groupname_}&serverid={_serverid_}	output	[@]{_botname_} !(spotify playlist)	t	0	\N	\N	\N	0
42	TEXT/TRIGGER/DISABLE	{_trigger(spotifyPlaylist)_}	\N	[@]{_botname_} !(disable spotify tracking)	t	0	11	No longer storing spotify matches in a playlist		0
40	TEXT/TRIGGER/ENABLE	{_trigger(spotifyPlaylist)_}	\N	[@]{_botname_} !(enable spotify tracking)	t	0	11	Now storing spotify matches in a playlist	\N	0
0	TEXT	https://media.giphy.com/media/3oEdv22bKDUluFKkxi/giphy.gif	\N	[@]{_botname_} (show me the money!*|[💰💵💸]*)	t	0	\N	\N	Money baby!	0
14	URL/TEXT	https://www.googleapis.com/youtube/v3/search?key={_key(yt)_}&part=snippet&maxResults=1&type=video&q={_text_}	items.{_randomInt_}.id.videoId	[@]{_botname_} !youtube (.*)	t	4	15	https://www.youtube.com/watch?v={_text_}	First result for YouTube	0
13	URL/TEXT	https://www.googleapis.com/youtube/v3/search?key={_key(yt)_}&part=snippet&type=video&q={_text_}	items.{_randomInt_}.id.videoId	[@]{_botname_} youtube (.*)	t	3	15	https://www.youtube.com/watch?v={_text_}	Gives a random YouTube video	0
22	URL/URL	http://{_me_}/youtubeSong?link={_text_}&groupid={_groupid_}&groupName={_groupname_}&serverid={_serverid_}	pageUrl	(https://www.youtube.com/watch\\?v=.+|https://open.spotify.com/track/.+|https://music.apple.com/us/album/.+)	t	0	\N	Found on Spotify, Apple Music, and Google Play Music	\N	1
12	URL/IMAGE	http://api.giphy.com/v1/gifs/search?q={_text_}&api_key={_key(giphy)_}&limit=40	data.{_randomInt_}.images.original.url	[@]{_botname_} (.*)	t	20	10	\N	Gif search	0
15	URL/TEXT	https://www.reddit.com/r/copypasta/top.json?sort=top&t=all&limit=100	data.children.{_randomInt_}.data.selftext	[@]{_botname_} !(pasta)	t	19	11	\N	Get your helping of pasta	0
\.


--
-- Data for Name: bots; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.bots (group_id, group_name, bot_name, key) FROM stdin;
9214876	Test	CalmanBot	17df3e511d405eb60452a767c8
discord	all_discord_groups	CalmanBot	PLACEHOLDER
\.


--
-- Data for Name: cached; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.cached (id, query, result) FROM stdin;
\.


--
-- Data for Name: discord_status; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.discord_status (id, type, text) FROM stdin;
1	1	to some sick beats
2	0	shit that doesn't suck yo
3	2	something I guess?
4	2	whatever Zach thinks is good
5	3	to Zach ask "Could this game get any worse?"
6	0	nothing, I have no games :(
7	2	these l33t skillz
8	1	to Alexa play Despacito
9	2	Jeff Goldblum movies for quotes
10	2	jet fuel not melt steel beams
11	2	jet beams not melt steel fuel
12	2	jet steel not melt beams fuel
13	2	steel beams melt jet fuel
14	2	Gazorpazorpfield
15	0	Star Citizen
16	2	The Tempest 2: Here we blow again
\.


--
-- Data for Name: discord_triggers; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.discord_triggers (id, channel_id, trigger_id, server_id) FROM stdin;
\.


--
-- Data for Name: groupme_posts; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.groupme_posts (id, cache_id, likes, message_id, group_id, posted_at) FROM stdin;
\.


--
-- Data for Name: minecraft_servers; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.minecraft_servers (id, address, name) FROM stdin;
\.


--
-- Data for Name: spotify_playlists; Type: TABLE DATA; Schema: public; Owner: nelsonleduc
--

COPY public.spotify_playlists (id, group_id, playlist_id) FROM stdin;
\.


--
-- Name: actions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.actions_id_seq', 22, true);


--
-- Name: bots_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.bots_id_seq', 3, true);


--
-- Name: cached_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.cached_id_seq', 6167, true);


--
-- Name: discord_status_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.discord_status_id_seq', 16, true);


--
-- Name: discord_triggers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.discord_triggers_id_seq', 135, true);


--
-- Name: groupme_posts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.groupme_posts_id_seq', 4341, true);


--
-- Name: minecraft_servers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.minecraft_servers_id_seq', 58, true);


--
-- Name: spotify_playlists_id_seq; Type: SEQUENCE SET; Schema: public; Owner: nelsonleduc
--

SELECT pg_catalog.setval('public.spotify_playlists_id_seq', 11, true);


--
-- Name: actions actions_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.actions
    ADD CONSTRAINT actions_pkey PRIMARY KEY (id);


--
-- Name: bots bots_group_id_key; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.bots
    ADD CONSTRAINT bots_group_id_key UNIQUE (group_id);


--
-- Name: bots bots_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.bots
    ADD CONSTRAINT bots_pkey PRIMARY KEY (group_id);


--
-- Name: cached cached_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.cached
    ADD CONSTRAINT cached_pkey PRIMARY KEY (id);


--
-- Name: discord_triggers discord_triggers_pk; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.discord_triggers
    ADD CONSTRAINT discord_triggers_pk PRIMARY KEY (id);


--
-- Name: discord_triggers discord_triggers_un; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.discord_triggers
    ADD CONSTRAINT discord_triggers_un UNIQUE (channel_id, trigger_id, server_id);


--
-- Name: groupme_posts groupme_posts_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.groupme_posts
    ADD CONSTRAINT groupme_posts_pkey PRIMARY KEY (id);


--
-- Name: minecraft_servers minecraft_servers_pk; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.minecraft_servers
    ADD CONSTRAINT minecraft_servers_pk PRIMARY KEY (id);


--
-- Name: minecraft_servers minecraft_servers_un; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.minecraft_servers
    ADD CONSTRAINT minecraft_servers_un UNIQUE (address);


--
-- Name: spotify_playlists spotify_playlists_pkey; Type: CONSTRAINT; Schema: public; Owner: nelsonleduc
--

ALTER TABLE ONLY public.spotify_playlists
    ADD CONSTRAINT spotify_playlists_pkey PRIMARY KEY (id);


--
-- Name: discord_triggers_trigger_id_idx; Type: INDEX; Schema: public; Owner: nelsonleduc
--

CREATE INDEX discord_triggers_trigger_id_idx ON public.discord_triggers USING btree (trigger_id);


--
-- Name: SCHEMA heroku_ext; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA heroku_ext TO nelsonleduc WITH GRANT OPTION;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: nelsonleduc
--

REVOKE ALL ON SCHEMA public FROM postgres;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO nelsonleduc;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- Name: LANGUAGE plpgsql; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON LANGUAGE plpgsql TO nelsonleduc;


--
-- PostgreSQL database dump complete
--

