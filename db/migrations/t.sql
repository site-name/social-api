--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.addresses (
    id character varying(36) NOT NULL,
    firstname character varying(64),
    lastname character varying(64),
    companyname character varying(256),
    streetaddress1 character varying(256),
    streetaddress2 character varying(256),
    city character varying(256),
    cityarea character varying(128),
    postalcode character varying(20),
    country character varying(3),
    countryarea character varying(128),
    phone character varying(20),
    createat bigint,
    updateat bigint
);


ALTER TABLE public.addresses OWNER TO minh;

--
-- Name: allocations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.allocations (
    id character varying(36) NOT NULL,
    createat bigint,
    orderlineid character varying(36),
    stockid character varying(36),
    quantityallocated integer
);


ALTER TABLE public.allocations OWNER TO minh;

--
-- Name: apps; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.apps (
    id character varying(36) NOT NULL,
    name character varying(60),
    createat bigint,
    isactive boolean,
    type text,
    identifier character varying(256),
    aboutapp text,
    dataprivacy text,
    dataprivacyurl text,
    homepageurl text,
    supporturl text,
    configurationurl text,
    appurl text,
    version character varying(60)
);


ALTER TABLE public.apps OWNER TO minh;

--
-- Name: apptokens; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.apptokens (
    id character varying(36) NOT NULL,
    appid character varying(36),
    name character varying(128),
    authtoken character varying(30)
);


ALTER TABLE public.apptokens OWNER TO minh;

--
-- Name: assignedpageattributes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedpageattributes (
    id character varying(36) NOT NULL,
    pageid character varying(36),
    assignmentid character varying(36)
);


ALTER TABLE public.assignedpageattributes OWNER TO minh;

--
-- Name: assignedpageattributevalues; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedpageattributevalues (
    id character varying(36) NOT NULL,
    valueid character varying(36),
    assignmentid character varying(36),
    sortorder integer
);


ALTER TABLE public.assignedpageattributevalues OWNER TO minh;

--
-- Name: assignedproductattributes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedproductattributes (
    id character varying(36) NOT NULL,
    productid character varying(36),
    assignmentid character varying(36)
);


ALTER TABLE public.assignedproductattributes OWNER TO minh;

--
-- Name: assignedproductattributevalues; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedproductattributevalues (
    id character varying(36) NOT NULL,
    valueid character varying(36),
    assignmentid character varying(36),
    sortorder integer
);


ALTER TABLE public.assignedproductattributevalues OWNER TO minh;

--
-- Name: assignedvariantattributes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedvariantattributes (
    id character varying(36) NOT NULL,
    variantid character varying(36),
    assignmentid character varying(36)
);


ALTER TABLE public.assignedvariantattributes OWNER TO minh;

--
-- Name: assignedvariantattributevalues; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.assignedvariantattributevalues (
    id character varying(36) NOT NULL,
    valueid character varying(36),
    assignmentid character varying(36),
    sortorder integer
);


ALTER TABLE public.assignedvariantattributevalues OWNER TO minh;

--
-- Name: attributepages; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributepages (
    id character varying(36) NOT NULL,
    attributeid character varying(36),
    pagetypeid character varying(36),
    sortorder integer
);


ALTER TABLE public.attributepages OWNER TO minh;

--
-- Name: attributeproducts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributeproducts (
    id character varying(36) NOT NULL,
    attributeid character varying(36),
    producttypeid character varying(36),
    sortorder integer
);


ALTER TABLE public.attributeproducts OWNER TO minh;

--
-- Name: attributes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributes (
    id character varying(36) NOT NULL,
    slug character varying(255),
    name character varying(250),
    type character varying(50),
    inputtype character varying(50),
    entitytype character varying(50),
    unit character varying(100),
    valuerequired boolean,
    isvariantonly boolean,
    visibleinstorefront boolean,
    filterableinstorefront boolean,
    filterableindashboard boolean,
    storefrontsearchposition integer,
    availableingrid boolean,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.attributes OWNER TO minh;

--
-- Name: attributetranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributetranslations (
    id character varying(36) NOT NULL,
    attributeid character varying(36),
    languagecode character varying(5),
    name character varying(100)
);


ALTER TABLE public.attributetranslations OWNER TO minh;

--
-- Name: attributevalues; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributevalues (
    id character varying(36) NOT NULL,
    name character varying(250),
    value character varying(9),
    slug character varying(255),
    fileurl character varying(200),
    contenttype character varying(50),
    attributeid character varying(36),
    richtext text,
    "boolean" boolean,
    datetime timestamp with time zone,
    sortorder integer
);


ALTER TABLE public.attributevalues OWNER TO minh;

--
-- Name: attributevaluetranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributevaluetranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    attributevalueid character varying(36),
    name character varying(100),
    richtext text
);


ALTER TABLE public.attributevaluetranslations OWNER TO minh;

--
-- Name: attributevariants; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.attributevariants (
    id character varying(36) NOT NULL,
    attributeid character varying(36),
    producttypeid character varying(36),
    variantselection boolean,
    sortorder integer
);


ALTER TABLE public.attributevariants OWNER TO minh;

--
-- Name: audits; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.audits (
    id character varying(36) NOT NULL,
    createat bigint,
    userid character varying(36),
    action character varying(512),
    extrainfo character varying(1024),
    ipaddress character varying(64),
    sessionid character varying(36)
);


ALTER TABLE public.audits OWNER TO minh;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.categories (
    id character varying(36) NOT NULL,
    name character varying(250),
    slug character varying(255),
    description text,
    parentid character varying(36),
    backgroundimage character varying(200),
    backgroundimagealt character varying(128),
    seotitle character varying(70),
    seodescription character varying(300),
    metadata text,
    privatemetadata text
);


ALTER TABLE public.categories OWNER TO minh;

--
-- Name: categorytranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.categorytranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    categoryid character varying(36),
    name character varying(250),
    description text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.categorytranslations OWNER TO minh;

--
-- Name: channels; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.channels (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    name character varying(250),
    isactive boolean,
    slug character varying(255),
    currency text,
    defaultcountry character varying(5)
);


ALTER TABLE public.channels OWNER TO minh;

--
-- Name: checkoutlines; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.checkoutlines (
    id character varying(36) NOT NULL,
    createat bigint,
    checkoutid character varying(36),
    variantid character varying(36),
    quantity integer
);


ALTER TABLE public.checkoutlines OWNER TO minh;

--
-- Name: checkouts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.checkouts (
    token character varying(36) NOT NULL,
    createat bigint,
    updateat bigint,
    userid character varying(36),
    shopid character varying(36),
    email text,
    quantity integer,
    channelid character varying(36),
    billingaddressid character varying(36),
    shippingaddressid character varying(36),
    shippingmethodid character varying(36),
    collectionpointid character varying(36),
    note text,
    currency text,
    country character varying(5),
    discountamount double precision,
    discountname character varying(255),
    translateddiscountname character varying(255),
    vouchercode character varying(12),
    redirecturl text,
    trackingcode character varying(255),
    languagecode text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.checkouts OWNER TO minh;

--
-- Name: clusterdiscovery; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.clusterdiscovery (
    id character varying(36) NOT NULL,
    type character varying(64),
    clustername character varying(64),
    hostname character varying(512),
    gossipport integer,
    port integer,
    createat bigint,
    lastpingat bigint
);


ALTER TABLE public.clusterdiscovery OWNER TO minh;

--
-- Name: collectionchannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.collectionchannellistings (
    id character varying(36) NOT NULL,
    createat bigint,
    collectionid character varying(36),
    channelid character varying(36),
    publicationdate timestamp with time zone,
    ispublished boolean
);


ALTER TABLE public.collectionchannellistings OWNER TO minh;

--
-- Name: collections; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.collections (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    name character varying(250),
    slug character varying(255),
    backgroundimage character varying(200),
    backgroundimagealt character varying(128),
    description text,
    metadata text,
    privatemetadata text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.collections OWNER TO minh;

--
-- Name: collectiontranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.collectiontranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    collectionid character varying(36),
    name character varying(250),
    description text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.collectiontranslations OWNER TO minh;

--
-- Name: compliances; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.compliances (
    id character varying(36) NOT NULL,
    createat bigint,
    userid character varying(36),
    status character varying(64),
    count integer,
    "desc" character varying(512),
    type character varying(64),
    startat bigint,
    endat bigint,
    keywords character varying(512),
    emails character varying(1024)
);


ALTER TABLE public.compliances OWNER TO minh;

--
-- Name: customerevents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.customerevents (
    id character varying(36) NOT NULL,
    date bigint,
    type character varying(255),
    orderid character varying(36),
    userid character varying(36),
    parameters text
);


ALTER TABLE public.customerevents OWNER TO minh;

--
-- Name: customernotes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.customernotes (
    id character varying(36),
    userid character varying(36),
    date bigint,
    content text,
    ispublic boolean,
    customerid character varying(36)
);


ALTER TABLE public.customernotes OWNER TO minh;

--
-- Name: db_lock; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.db_lock (
    id character varying(64) NOT NULL,
    expireat bigint
);


ALTER TABLE public.db_lock OWNER TO minh;

--
-- Name: db_migrations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.db_migrations (
    version bigint NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.db_migrations OWNER TO minh;

--
-- Name: digitalcontents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.digitalcontents (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    usedefaultsettings boolean,
    automaticfulfillment boolean,
    contenttype character varying(128),
    productvariantid character varying(36),
    contentfile character varying(200),
    maxdownloads integer,
    urlvaliddays integer,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.digitalcontents OWNER TO minh;

--
-- Name: digitalcontenturls; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.digitalcontenturls (
    id character varying(36) NOT NULL,
    token character varying(36),
    contentid character varying(36),
    createat bigint,
    downloadnum integer,
    lineid character varying(36)
);


ALTER TABLE public.digitalcontenturls OWNER TO minh;

--
-- Name: exportevents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.exportevents (
    id character varying(36) NOT NULL,
    date bigint,
    type character varying(255),
    parameters text,
    exportfileid character varying(36),
    userid character varying(36)
);


ALTER TABLE public.exportevents OWNER TO minh;

--
-- Name: exportfiles; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.exportfiles (
    id character varying(36) NOT NULL,
    userid character varying(36),
    contentfile text,
    createat bigint,
    updateat bigint
);


ALTER TABLE public.exportfiles OWNER TO minh;

--
-- Name: fileinfos; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.fileinfos (
    id character varying(36) NOT NULL,
    creatorid character varying(36),
    parentid character varying(36),
    createat bigint,
    updateat bigint,
    deleteat bigint,
    path character varying(512),
    thumbnailpath character varying(512),
    previewpath character varying(512),
    name character varying(256),
    extension character varying(64),
    size bigint,
    mimetype character varying(256),
    width integer,
    height integer,
    haspreviewimage boolean,
    minipreview bytea,
    content text,
    remoteid character varying(26)
);


ALTER TABLE public.fileinfos OWNER TO minh;

--
-- Name: fulfillmentlines; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.fulfillmentlines (
    id character varying(36) NOT NULL,
    orderlineid character varying(36),
    fulfillmentid character varying(36),
    quantity integer,
    stockid character varying(36)
);


ALTER TABLE public.fulfillmentlines OWNER TO minh;

--
-- Name: fulfillments; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.fulfillments (
    id character varying(36) NOT NULL,
    fulfillmentorder integer,
    orderid character varying(36),
    status character varying(32),
    trackingnumber character varying(255),
    createat bigint,
    shippingrefundamount double precision,
    totalrefundamount double precision,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.fulfillments OWNER TO minh;

--
-- Name: giftcardcheckouts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.giftcardcheckouts (
    id character varying(36) NOT NULL,
    giftcardid character varying(36),
    checkoutid character varying(36)
);


ALTER TABLE public.giftcardcheckouts OWNER TO minh;

--
-- Name: giftcardevents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.giftcardevents (
    id character varying(36) NOT NULL,
    date bigint,
    type character varying(255),
    parameters jsonb,
    userid character varying(36),
    giftcardid character varying(36)
);


ALTER TABLE public.giftcardevents OWNER TO minh;

--
-- Name: giftcards; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.giftcards (
    id character varying(36) NOT NULL,
    code character varying(40),
    createdbyid character varying(36),
    usedbyid character varying(36),
    createdbyemail character varying(128),
    usedbyemail character varying(128),
    createat bigint,
    startdate timestamp with time zone,
    expirydate timestamp with time zone,
    tag character varying(255),
    productid character varying(36),
    lastusedon bigint,
    isactive boolean,
    currency character varying(3),
    initialbalanceamount double precision,
    currentbalanceamount double precision,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.giftcards OWNER TO minh;

--
-- Name: invoiceevents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.invoiceevents (
    id character varying(36) NOT NULL,
    createat bigint,
    type character varying(255),
    invoiceid character varying(36),
    orderid character varying(36),
    userid character varying(36),
    parameters text
);


ALTER TABLE public.invoiceevents OWNER TO minh;

--
-- Name: invoices; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.invoices (
    id character varying(36) NOT NULL,
    orderid character varying(36),
    number character varying(255),
    createat bigint,
    externalurl character varying(2048),
    invoicefile text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.invoices OWNER TO minh;

--
-- Name: jobs; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.jobs (
    id character varying(36) NOT NULL,
    type character varying(32),
    priority bigint,
    createat bigint,
    startat bigint,
    lastactivityat bigint,
    status character varying(32),
    progress bigint,
    data jsonb
);


ALTER TABLE public.jobs OWNER TO minh;

--
-- Name: menuitems; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.menuitems (
    id character varying(36) NOT NULL,
    menuid character varying(36),
    name character varying(128),
    parentid character varying(36),
    url character varying(256),
    categoryid character varying(36),
    collectionid character varying(36),
    pageid character varying(36),
    metadata text,
    privatemetadata text,
    sortorder integer
);


ALTER TABLE public.menuitems OWNER TO minh;

--
-- Name: menuitemtranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.menuitemtranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(10),
    menuitemid character varying(36),
    name character varying(128)
);


ALTER TABLE public.menuitemtranslations OWNER TO minh;

--
-- Name: menus; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.menus (
    id character varying(36) NOT NULL,
    name character varying(250),
    slug character varying(255),
    createat bigint,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.menus OWNER TO minh;

--
-- Name: openexchangerates; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.openexchangerates (
    id character varying(36) NOT NULL,
    tocurrency character varying(3),
    rate double precision
);


ALTER TABLE public.openexchangerates OWNER TO minh;

--
-- Name: orderdiscounts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.orderdiscounts (
    id text NOT NULL,
    orderid character varying(36),
    type character varying(10),
    valuetype character varying(10),
    value double precision,
    amountvalue double precision,
    currency text,
    name character varying(255),
    translatedname character varying(255),
    reason text
);


ALTER TABLE public.orderdiscounts OWNER TO minh;

--
-- Name: orderevents; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.orderevents (
    id character varying(36) NOT NULL,
    createat bigint,
    type character varying(255),
    orderid character varying(36),
    parameters text,
    userid character varying(36)
);


ALTER TABLE public.orderevents OWNER TO minh;

--
-- Name: ordergiftcards; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.ordergiftcards (
    id character varying(36) NOT NULL,
    giftcardid character varying(36),
    orderid character varying(36)
);


ALTER TABLE public.ordergiftcards OWNER TO minh;

--
-- Name: orderlines; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.orderlines (
    id character varying(36) NOT NULL,
    createat bigint,
    orderid character varying(36),
    variantid character varying(36),
    productname character varying(386),
    variantname character varying(255),
    translatedproductname character varying(386),
    translatedvariantname character varying(255),
    productsku character varying(255),
    productvariantid character varying(255),
    isshippingrequired boolean,
    isgiftcard boolean,
    quantity integer,
    quantityfulfilled integer,
    currency character varying(3),
    unitdiscountamount double precision,
    unitdiscounttype character varying(10),
    unitdiscountreason text,
    unitpricenetamount double precision,
    unitdiscountvalue double precision,
    unitpricegrossamount double precision,
    totalpricenetamount double precision,
    totalpricegrossamount double precision,
    undiscountedunitpricegrossamount double precision,
    undiscountedunitpricenetamount double precision,
    undiscountedtotalpricegrossamount double precision,
    undiscountedtotalpricenetamount double precision,
    taxrate double precision,
    allocations text
);


ALTER TABLE public.orderlines OWNER TO minh;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.orders (
    id character varying(36) NOT NULL,
    createat bigint,
    status character varying(32),
    userid character varying(36),
    shopid character varying(36),
    languagecode character varying(5),
    trackingclientid character varying(36),
    billingaddressid character varying(36),
    shippingaddressid character varying(36),
    useremail character varying(128),
    originalid character varying(36),
    origin character varying(32),
    currency character varying(200),
    shippingmethodid character varying(36),
    collectionpointid character varying(36),
    shippingmethodname character varying(255),
    collectionpointname character varying(255),
    channelid character varying(36),
    shippingpricenetamount double precision,
    shippingpricegrossamount double precision,
    shippingtaxrate double precision,
    token character varying(36),
    checkouttoken character varying(36),
    totalnetamount double precision,
    undiscountedtotalnetamount double precision,
    totalgrossamount double precision,
    undiscountedtotalgrossamount double precision,
    totalpaidamount double precision,
    voucherid character varying(36),
    displaygrossprices boolean,
    customernote text,
    weightamount real,
    weightunit text,
    redirecturl text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.orders OWNER TO minh;

--
-- Name: pages; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.pages (
    id character varying(36) NOT NULL,
    title character varying(250),
    slug character varying(255),
    pagetypeid character varying(36),
    content text,
    createat bigint,
    metadata text,
    privatemetadata text,
    publicationdate timestamp with time zone,
    ispublished boolean,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.pages OWNER TO minh;

--
-- Name: pagetranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.pagetranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    pageid character varying(36),
    title character varying(250),
    content text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.pagetranslations OWNER TO minh;

--
-- Name: pagetypes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.pagetypes (
    id character varying(36) NOT NULL,
    name character varying(250),
    slug character varying(255),
    metadata text,
    privatemetadata text
);


ALTER TABLE public.pagetypes OWNER TO minh;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.payments (
    id character varying(36) NOT NULL,
    gateway character varying(255),
    isactive boolean,
    toconfirm boolean,
    createat bigint,
    updateat bigint,
    chargestatus character varying(20),
    token character varying(512),
    total double precision,
    capturedamount double precision,
    currency character varying(3),
    checkoutid character varying(36),
    orderid character varying(36),
    billingemail character varying(128),
    billingfirstname character varying(256),
    billinglastname character varying(256),
    billingcompanyname character varying(256),
    billingaddress1 character varying(256),
    billingaddress2 character varying(256),
    billingcity character varying(256),
    billingcityarea character varying(128),
    billingpostalcode character varying(20),
    billingcountrycode character varying(5),
    billingcountryarea character varying(256),
    ccfirstdigits character varying(6),
    cclastdigits character varying(4),
    ccbrand character varying(40),
    ccexpmonth integer,
    ccexpyear integer,
    paymentmethodtype character varying(256),
    customeripaddress character varying(39),
    extradata text,
    returnurl character varying(200),
    pspreference character varying(512),
    storepaymentmethod character varying(11),
    metadata text,
    privatemetadata text
);


ALTER TABLE public.payments OWNER TO minh;

--
-- Name: pluginconfigurations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.pluginconfigurations (
    id character varying(36) NOT NULL,
    identifier character varying(128),
    name character varying(128),
    channelid character varying(36),
    description character varying(1000),
    active boolean,
    configuration text
);


ALTER TABLE public.pluginconfigurations OWNER TO minh;

--
-- Name: pluginkeyvaluestore; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.pluginkeyvaluestore (
    pluginid character varying(190) NOT NULL,
    pkey character varying(50) NOT NULL,
    pvalue bytea,
    expireat bigint
);


ALTER TABLE public.pluginkeyvaluestore OWNER TO minh;

--
-- Name: preferences; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.preferences (
    userid character varying(36) NOT NULL,
    category character varying(32) NOT NULL,
    name character varying(32) NOT NULL,
    value character varying(2000)
);


ALTER TABLE public.preferences OWNER TO minh;

--
-- Name: preorderallocations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.preorderallocations (
    id character varying(36) NOT NULL,
    orderlineid character varying(36),
    quantity integer,
    productvariantchannellistingid character varying(36)
);


ALTER TABLE public.preorderallocations OWNER TO minh;

--
-- Name: productchannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productchannellistings (
    id character varying(36) NOT NULL,
    productid character varying(36),
    channelid character varying(36),
    visibleinlistings boolean,
    availableforpurchase timestamp with time zone,
    currency character varying(3),
    discountedpriceamount double precision,
    createat bigint,
    publicationdate timestamp with time zone,
    ispublished boolean
);


ALTER TABLE public.productchannellistings OWNER TO minh;

--
-- Name: productcollections; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productcollections (
    id character varying(36) NOT NULL,
    collectionid character varying(36),
    productid character varying(36)
);


ALTER TABLE public.productcollections OWNER TO minh;

--
-- Name: productmedias; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productmedias (
    id character varying(36) NOT NULL,
    createat bigint,
    productid character varying(36),
    ppoi character varying(20),
    image character varying(200),
    alt character varying(128),
    type character varying(32),
    externalurl character varying(256),
    oembeddata text,
    sortorder integer
);


ALTER TABLE public.productmedias OWNER TO minh;

--
-- Name: products; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.products (
    id character varying(36) NOT NULL,
    producttypeid character varying(36),
    name character varying(250),
    slug character varying(255),
    description text,
    descriptionplaintext text,
    categoryid character varying(36),
    createat bigint,
    updateat bigint,
    chargetaxes boolean,
    weight real,
    weightunit text,
    defaultvariantid character varying(36),
    rating real,
    metadata text,
    privatemetadata text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.products OWNER TO minh;

--
-- Name: producttranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.producttranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    productid character varying(36),
    name character varying(250),
    description text,
    seotitle character varying(70),
    seodescription character varying(300)
);


ALTER TABLE public.producttranslations OWNER TO minh;

--
-- Name: producttypes; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.producttypes (
    id character varying(36) NOT NULL,
    name character varying(250),
    slug character varying(255),
    kind character varying(32),
    hasvariants boolean,
    isshippingrequired boolean,
    isdigital boolean,
    weight real,
    weightunit text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.producttypes OWNER TO minh;

--
-- Name: productvariantchannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productvariantchannellistings (
    id character varying(36) NOT NULL,
    variantid character varying(36) NOT NULL,
    channelid character varying(36) NOT NULL,
    currency character varying(3),
    priceamount double precision,
    costpriceamount double precision,
    preorderquantitythreshold integer,
    createat bigint
);


ALTER TABLE public.productvariantchannellistings OWNER TO minh;

--
-- Name: productvariants; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productvariants (
    id character varying(36) NOT NULL,
    name character varying(255),
    productid character varying(36),
    sku character varying(255),
    weight real,
    weightunit text,
    trackinventory boolean,
    ispreorder boolean,
    preorderenddate bigint,
    preorderglobalthreshold integer,
    sortorder integer,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.productvariants OWNER TO minh;

--
-- Name: productvarianttranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.productvarianttranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(5),
    productvariantid character varying(36),
    name character varying(255)
);


ALTER TABLE public.productvarianttranslations OWNER TO minh;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.roles (
    id character varying(36) NOT NULL,
    name character varying(64),
    displayname character varying(128),
    description character varying(1024),
    createat bigint,
    updateat bigint,
    deleteat bigint,
    permissionsstr text,
    schememanaged boolean,
    builtin boolean
);


ALTER TABLE public.roles OWNER TO minh;

--
-- Name: salecategories; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.salecategories (
    id character varying(36) NOT NULL,
    saleid character varying(36),
    categoryid character varying(36),
    createat bigint
);


ALTER TABLE public.salecategories OWNER TO minh;

--
-- Name: salechannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.salechannellistings (
    id character varying(36) NOT NULL,
    saleid character varying(36),
    channelid character varying(36) NOT NULL,
    discountvalue double precision,
    currency text,
    createat bigint
);


ALTER TABLE public.salechannellistings OWNER TO minh;

--
-- Name: salecollections; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.salecollections (
    id character varying(36) NOT NULL,
    saleid character varying(36),
    collectionid character varying(36),
    createat bigint
);


ALTER TABLE public.salecollections OWNER TO minh;

--
-- Name: saleproducts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.saleproducts (
    id character varying(36) NOT NULL,
    saleid character varying(36),
    productid character varying(36),
    createat bigint
);


ALTER TABLE public.saleproducts OWNER TO minh;

--
-- Name: saleproductvariants; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.saleproductvariants (
    id character varying(36) NOT NULL,
    saleid character varying(36),
    productvariantid character varying(36),
    createat bigint
);


ALTER TABLE public.saleproductvariants OWNER TO minh;

--
-- Name: sales; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.sales (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    name character varying(255),
    type character varying(10),
    startdate bigint,
    enddate bigint,
    createat bigint,
    updateat bigint,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.sales OWNER TO minh;

--
-- Name: saletranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.saletranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(10),
    name character varying(255),
    saleid character varying(36)
);


ALTER TABLE public.saletranslations OWNER TO minh;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.sessions (
    id character varying(36) NOT NULL,
    token character varying(36),
    createat bigint,
    expiresat bigint,
    lastactivityat bigint,
    userid character varying(36),
    deviceid character varying(512),
    roles character varying(64),
    isoauth boolean,
    expirednotify boolean,
    props character varying(1000)
);


ALTER TABLE public.sessions OWNER TO minh;

--
-- Name: shippingmethodchannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingmethodchannellistings (
    id character varying(36) NOT NULL,
    shippingmethodid character varying(36),
    channelid character varying(36),
    minimumorderpriceamount double precision,
    currency character varying(3),
    maximumorderpriceamount double precision,
    priceamount double precision,
    createat bigint
);


ALTER TABLE public.shippingmethodchannellistings OWNER TO minh;

--
-- Name: shippingmethodexcludedproducts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingmethodexcludedproducts (
    id character varying(36) NOT NULL,
    shippingmethodid character varying(36),
    productid character varying(36)
);


ALTER TABLE public.shippingmethodexcludedproducts OWNER TO minh;

--
-- Name: shippingmethodpostalcoderules; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingmethodpostalcoderules (
    id character varying(36) NOT NULL,
    shippingmethodid character varying(36),
    start character varying(32),
    "end" character varying(32),
    inclusiontype character varying(32)
);


ALTER TABLE public.shippingmethodpostalcoderules OWNER TO minh;

--
-- Name: shippingmethods; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingmethods (
    id character varying(36) NOT NULL,
    name character varying(100),
    type character varying(30),
    shippingzoneid character varying(36),
    minimumorderweight real,
    maximumorderweight real,
    weightunit character varying(5),
    maximumdeliverydays integer,
    minimumdeliverydays integer,
    description text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.shippingmethods OWNER TO minh;

--
-- Name: shippingmethodtranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingmethodtranslations (
    id character varying(36) NOT NULL,
    shippingmethodid character varying(36),
    languagecode character varying(5),
    name character varying(100),
    description text
);


ALTER TABLE public.shippingmethodtranslations OWNER TO minh;

--
-- Name: shippingzonechannels; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingzonechannels (
    id character varying(36) NOT NULL,
    shippingzoneid character varying(36),
    channelid character varying(36)
);


ALTER TABLE public.shippingzonechannels OWNER TO minh;

--
-- Name: shippingzones; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shippingzones (
    id character varying(36) NOT NULL,
    name character varying(100),
    countries character varying(749),
    "default" boolean,
    description text,
    createat bigint,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.shippingzones OWNER TO minh;

--
-- Name: shops; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shops (
    id character varying(36) NOT NULL,
    ownerid character varying(36),
    createat bigint,
    updateat bigint,
    name character varying(100),
    description character varying(200),
    topmenuid character varying(36),
    includetaxesinprice boolean,
    displaygrossprices boolean,
    chargetaxesonshipping boolean,
    trackinventorybydefault boolean,
    defaultweightunit character varying(10),
    automaticfulfillmentdigitalproducts boolean,
    defaultdigitalmaxdownloads integer,
    defaultdigitalurlvaliddays integer,
    addressid character varying(36),
    defaultmailsendername character varying(78),
    defaultmailsenderaddress text,
    customersetpasswordurl text,
    automaticallyconfirmallneworders boolean,
    fulfillmentautoapprove boolean,
    fulfillmentallowunpaid boolean,
    giftcardexpirytype character varying(32),
    giftcardexpiryperiodtype character varying(32),
    giftcardexpiryperiod integer,
    automaticallyfulfillnonshippablegiftcard boolean
);


ALTER TABLE public.shops OWNER TO minh;

--
-- Name: shopstaffs; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shopstaffs (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    staffid character varying(36),
    createat bigint,
    endat bigint
);


ALTER TABLE public.shopstaffs OWNER TO minh;

--
-- Name: shoptranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.shoptranslations (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    languagecode character varying(5),
    name character varying(110),
    description character varying(110),
    createat bigint,
    updateat bigint
);


ALTER TABLE public.shoptranslations OWNER TO minh;

--
-- Name: staffnotificationrecipients; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.staffnotificationrecipients (
    id character varying(36) NOT NULL,
    userid character varying(36),
    staffemail character varying(128),
    active boolean
);


ALTER TABLE public.staffnotificationrecipients OWNER TO minh;

--
-- Name: status; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.status (
    userid character varying(36) NOT NULL,
    status character varying(32),
    manual boolean,
    lastactivityat bigint
);


ALTER TABLE public.status OWNER TO minh;

--
-- Name: stocks; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.stocks (
    id character varying(36) NOT NULL,
    createat bigint,
    warehouseid character varying(36),
    productvariantid character varying(36),
    quantity integer
);


ALTER TABLE public.stocks OWNER TO minh;

--
-- Name: systems; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.systems (
    name character varying(64) NOT NULL,
    value character varying(1024)
);


ALTER TABLE public.systems OWNER TO minh;

--
-- Name: termsofservices; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.termsofservices (
    id character varying(36) NOT NULL,
    createat bigint,
    userid character varying(36),
    text character varying(65535)
);


ALTER TABLE public.termsofservices OWNER TO minh;

--
-- Name: tokens; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.tokens (
    token character varying(64) NOT NULL,
    createat bigint,
    type character varying(64),
    extra character varying(2048)
);


ALTER TABLE public.tokens OWNER TO minh;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.transactions (
    id character varying(36) NOT NULL,
    createat bigint,
    paymentid character varying(36),
    token character varying(512),
    kind character varying(25),
    issuccess boolean,
    actionrequired boolean,
    actionrequireddata text,
    currency character varying(3),
    amount double precision,
    error character varying(256),
    customerid character varying(256),
    gatewayresponse text,
    alreadyprocessed boolean
);


ALTER TABLE public.transactions OWNER TO minh;

--
-- Name: uploadsessions; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.uploadsessions (
    id character varying(36) NOT NULL,
    type character varying(32),
    createat bigint,
    userid character varying(36),
    filename character varying(256),
    path character varying(512),
    filesize bigint,
    fileoffset bigint
);


ALTER TABLE public.uploadsessions OWNER TO minh;

--
-- Name: useraccesstokens; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.useraccesstokens (
    id character varying(36) NOT NULL,
    token character varying(36),
    userid character varying(36),
    description character varying(255),
    isactive boolean
);


ALTER TABLE public.useraccesstokens OWNER TO minh;

--
-- Name: useraddresses; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.useraddresses (
    id character varying(36) NOT NULL,
    userid character varying(36),
    addressid character varying(36)
);


ALTER TABLE public.useraddresses OWNER TO minh;

--
-- Name: users; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.users (
    id character varying(36) NOT NULL,
    email character varying(128),
    username character varying(64),
    firstname character varying(64),
    lastname character varying(64),
    defaultshippingaddressid character varying(36),
    defaultbillingaddressid character varying(36),
    password character varying(128),
    authdata character varying(128),
    authservice character varying(32),
    emailverified boolean,
    nickname character varying(64),
    roles character varying(256),
    props character varying(4000),
    notifyprops character varying(2000),
    lastpasswordupdate bigint,
    lastpictureupdate bigint,
    failedattempts integer,
    locale character varying(5),
    timezone character varying(256),
    mfaactive boolean,
    mfasecret character varying(128),
    createat bigint,
    updateat bigint,
    deleteat bigint,
    isactive boolean,
    note text,
    jwttokenkey text,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.users OWNER TO minh;

--
-- Name: usertermofservices; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.usertermofservices (
    userid character varying(36) NOT NULL,
    termsofserviceid character varying(36),
    createat bigint
);


ALTER TABLE public.usertermofservices OWNER TO minh;

--
-- Name: variantmedias; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.variantmedias (
    id character varying(36) NOT NULL,
    variantid character varying(36),
    mediaid character varying(36)
);


ALTER TABLE public.variantmedias OWNER TO minh;

--
-- Name: vouchercategories; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.vouchercategories (
    id character varying(36) NOT NULL,
    voucherid character varying(36),
    categoryid character varying(36),
    createat bigint
);


ALTER TABLE public.vouchercategories OWNER TO minh;

--
-- Name: voucherchannellistings; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.voucherchannellistings (
    id character varying(36) NOT NULL,
    createat bigint,
    voucherid character varying(36) NOT NULL,
    channelid character varying(36) NOT NULL,
    discountvalue double precision,
    currency character varying(3),
    minspenamount double precision
);


ALTER TABLE public.voucherchannellistings OWNER TO minh;

--
-- Name: vouchercollections; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.vouchercollections (
    id character varying(36) NOT NULL,
    voucherid character varying(36),
    collectionid character varying(36)
);


ALTER TABLE public.vouchercollections OWNER TO minh;

--
-- Name: vouchercustomers; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.vouchercustomers (
    id character varying(36) NOT NULL,
    voucherid character varying(36),
    customeremail character varying(128)
);


ALTER TABLE public.vouchercustomers OWNER TO minh;

--
-- Name: voucherproducts; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.voucherproducts (
    id character varying(36) NOT NULL,
    voucherid character varying(36),
    productid character varying(36)
);


ALTER TABLE public.voucherproducts OWNER TO minh;

--
-- Name: voucherproductvariants; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.voucherproductvariants (
    id character varying(36) NOT NULL,
    voucherid character varying(36),
    productvariantid character varying(36),
    createat bigint
);


ALTER TABLE public.voucherproductvariants OWNER TO minh;

--
-- Name: vouchers; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.vouchers (
    id character varying(36) NOT NULL,
    shopid character varying(36),
    type character varying(20),
    name character varying(255),
    code character varying(16),
    usagelimit integer,
    used integer,
    startdate bigint,
    enddate bigint,
    applyonceperorder boolean,
    applyoncepercustomer boolean,
    onlyforstaff boolean,
    discountvaluetype character varying(10),
    countries character varying(749),
    mincheckoutitemsquantity integer,
    createat bigint,
    updateat bigint,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.vouchers OWNER TO minh;

--
-- Name: vouchertranslations; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.vouchertranslations (
    id character varying(36) NOT NULL,
    languagecode character varying(10),
    name character varying(255),
    voucherid character varying(36),
    createat bigint
);


ALTER TABLE public.vouchertranslations OWNER TO minh;

--
-- Name: warehouses; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.warehouses (
    id character varying(36) NOT NULL,
    name character varying(250),
    slug character varying(255),
    addressid character varying(36),
    email character varying(128),
    clickandcollectoption character varying(30),
    isprivate boolean,
    metadata text,
    privatemetadata text
);


ALTER TABLE public.warehouses OWNER TO minh;

--
-- Name: warehouseshippingzones; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.warehouseshippingzones (
    id character varying(36) NOT NULL,
    warehouseid character varying(36),
    shippingzoneid character varying(36)
);


ALTER TABLE public.warehouseshippingzones OWNER TO minh;

--
-- Name: wishlistitemproductvariants; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.wishlistitemproductvariants (
    id character varying(36) NOT NULL,
    wishlistitemid character varying(36),
    productvariantid character varying(36)
);


ALTER TABLE public.wishlistitemproductvariants OWNER TO minh;

--
-- Name: wishlistitems; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.wishlistitems (
    id character varying(36) NOT NULL,
    wishlistid character varying(36),
    productid character varying(36),
    createat bigint
);


ALTER TABLE public.wishlistitems OWNER TO minh;

--
-- Name: wishlists; Type: TABLE; Schema: public; Owner: minh
--

CREATE TABLE public.wishlists (
    id character varying(36) NOT NULL,
    token character varying(36),
    userid character varying(36),
    createat bigint
);


ALTER TABLE public.wishlists OWNER TO minh;

--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.addresses (id, firstname, lastname, companyname, streetaddress1, streetaddress2, city, cityarea, postalcode, country, countryarea, phone, createat, updateat) FROM stdin;
\.


--
-- Data for Name: allocations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.allocations (id, createat, orderlineid, stockid, quantityallocated) FROM stdin;
\.


--
-- Data for Name: apps; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.apps (id, name, createat, isactive, type, identifier, aboutapp, dataprivacy, dataprivacyurl, homepageurl, supporturl, configurationurl, appurl, version) FROM stdin;
\.


--
-- Data for Name: apptokens; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.apptokens (id, appid, name, authtoken) FROM stdin;
\.


--
-- Data for Name: assignedpageattributes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedpageattributes (id, pageid, assignmentid) FROM stdin;
\.


--
-- Data for Name: assignedpageattributevalues; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedpageattributevalues (id, valueid, assignmentid, sortorder) FROM stdin;
\.


--
-- Data for Name: assignedproductattributes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedproductattributes (id, productid, assignmentid) FROM stdin;
\.


--
-- Data for Name: assignedproductattributevalues; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedproductattributevalues (id, valueid, assignmentid, sortorder) FROM stdin;
\.


--
-- Data for Name: assignedvariantattributes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedvariantattributes (id, variantid, assignmentid) FROM stdin;
\.


--
-- Data for Name: assignedvariantattributevalues; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.assignedvariantattributevalues (id, valueid, assignmentid, sortorder) FROM stdin;
\.


--
-- Data for Name: attributepages; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributepages (id, attributeid, pagetypeid, sortorder) FROM stdin;
\.


--
-- Data for Name: attributeproducts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributeproducts (id, attributeid, producttypeid, sortorder) FROM stdin;
\.


--
-- Data for Name: attributes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributes (id, slug, name, type, inputtype, entitytype, unit, valuerequired, isvariantonly, visibleinstorefront, filterableinstorefront, filterableindashboard, storefrontsearchposition, availableingrid, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: attributetranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributetranslations (id, attributeid, languagecode, name) FROM stdin;
\.


--
-- Data for Name: attributevalues; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributevalues (id, name, value, slug, fileurl, contenttype, attributeid, richtext, "boolean", datetime, sortorder) FROM stdin;
\.


--
-- Data for Name: attributevaluetranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributevaluetranslations (id, languagecode, attributevalueid, name, richtext) FROM stdin;
\.


--
-- Data for Name: attributevariants; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.attributevariants (id, attributeid, producttypeid, variantselection, sortorder) FROM stdin;
\.


--
-- Data for Name: audits; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.audits (id, createat, userid, action, extrainfo, ipaddress, sessionid) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.categories (id, name, slug, description, parentid, backgroundimage, backgroundimagealt, seotitle, seodescription, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: categorytranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.categorytranslations (id, languagecode, categoryid, name, description, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: channels; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.channels (id, shopid, name, isactive, slug, currency, defaultcountry) FROM stdin;
\.


--
-- Data for Name: checkoutlines; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.checkoutlines (id, createat, checkoutid, variantid, quantity) FROM stdin;
\.


--
-- Data for Name: checkouts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.checkouts (token, createat, updateat, userid, shopid, email, quantity, channelid, billingaddressid, shippingaddressid, shippingmethodid, collectionpointid, note, currency, country, discountamount, discountname, translateddiscountname, vouchercode, redirecturl, trackingcode, languagecode, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: clusterdiscovery; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.clusterdiscovery (id, type, clustername, hostname, gossipport, port, createat, lastpingat) FROM stdin;
\.


--
-- Data for Name: collectionchannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.collectionchannellistings (id, createat, collectionid, channelid, publicationdate, ispublished) FROM stdin;
\.


--
-- Data for Name: collections; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.collections (id, shopid, name, slug, backgroundimage, backgroundimagealt, description, metadata, privatemetadata, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: collectiontranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.collectiontranslations (id, languagecode, collectionid, name, description, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: compliances; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.compliances (id, createat, userid, status, count, "desc", type, startat, endat, keywords, emails) FROM stdin;
\.


--
-- Data for Name: customerevents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.customerevents (id, date, type, orderid, userid, parameters) FROM stdin;
\.


--
-- Data for Name: customernotes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.customernotes (id, userid, date, content, ispublic, customerid) FROM stdin;
\.


--
-- Data for Name: db_lock; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.db_lock (id, expireat) FROM stdin;
\.


--
-- Data for Name: db_migrations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.db_migrations (version, name) FROM stdin;
\.


--
-- Data for Name: digitalcontents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.digitalcontents (id, shopid, usedefaultsettings, automaticfulfillment, contenttype, productvariantid, contentfile, maxdownloads, urlvaliddays, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: digitalcontenturls; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.digitalcontenturls (id, token, contentid, createat, downloadnum, lineid) FROM stdin;
\.


--
-- Data for Name: exportevents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.exportevents (id, date, type, parameters, exportfileid, userid) FROM stdin;
\.


--
-- Data for Name: exportfiles; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.exportfiles (id, userid, contentfile, createat, updateat) FROM stdin;
\.


--
-- Data for Name: fileinfos; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.fileinfos (id, creatorid, parentid, createat, updateat, deleteat, path, thumbnailpath, previewpath, name, extension, size, mimetype, width, height, haspreviewimage, minipreview, content, remoteid) FROM stdin;
\.


--
-- Data for Name: fulfillmentlines; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.fulfillmentlines (id, orderlineid, fulfillmentid, quantity, stockid) FROM stdin;
\.


--
-- Data for Name: fulfillments; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.fulfillments (id, fulfillmentorder, orderid, status, trackingnumber, createat, shippingrefundamount, totalrefundamount, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: giftcardcheckouts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.giftcardcheckouts (id, giftcardid, checkoutid) FROM stdin;
\.


--
-- Data for Name: giftcardevents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.giftcardevents (id, date, type, parameters, userid, giftcardid) FROM stdin;
\.


--
-- Data for Name: giftcards; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.giftcards (id, code, createdbyid, usedbyid, createdbyemail, usedbyemail, createat, startdate, expirydate, tag, productid, lastusedon, isactive, currency, initialbalanceamount, currentbalanceamount, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: invoiceevents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.invoiceevents (id, createat, type, invoiceid, orderid, userid, parameters) FROM stdin;
\.


--
-- Data for Name: invoices; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.invoices (id, orderid, number, createat, externalurl, invoicefile, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: jobs; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.jobs (id, type, priority, createat, startat, lastactivityat, status, progress, data) FROM stdin;
65e30556-b36e-455a-b8bc-5895d259f95e	migrations	0	1660901377034	1660901388228	1660901388228	in_progress	0	{"last_done": "", "migration_key": "migration_advanced_permissions_phase_2"}
\.


--
-- Data for Name: menuitems; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.menuitems (id, menuid, name, parentid, url, categoryid, collectionid, pageid, metadata, privatemetadata, sortorder) FROM stdin;
\.


--
-- Data for Name: menuitemtranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.menuitemtranslations (id, languagecode, menuitemid, name) FROM stdin;
\.


--
-- Data for Name: menus; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.menus (id, name, slug, createat, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: openexchangerates; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.openexchangerates (id, tocurrency, rate) FROM stdin;
\.


--
-- Data for Name: orderdiscounts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.orderdiscounts (id, orderid, type, valuetype, value, amountvalue, currency, name, translatedname, reason) FROM stdin;
\.


--
-- Data for Name: orderevents; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.orderevents (id, createat, type, orderid, parameters, userid) FROM stdin;
\.


--
-- Data for Name: ordergiftcards; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.ordergiftcards (id, giftcardid, orderid) FROM stdin;
\.


--
-- Data for Name: orderlines; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.orderlines (id, createat, orderid, variantid, productname, variantname, translatedproductname, translatedvariantname, productsku, productvariantid, isshippingrequired, isgiftcard, quantity, quantityfulfilled, currency, unitdiscountamount, unitdiscounttype, unitdiscountreason, unitpricenetamount, unitdiscountvalue, unitpricegrossamount, totalpricenetamount, totalpricegrossamount, undiscountedunitpricegrossamount, undiscountedunitpricenetamount, undiscountedtotalpricegrossamount, undiscountedtotalpricenetamount, taxrate, allocations) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.orders (id, createat, status, userid, shopid, languagecode, trackingclientid, billingaddressid, shippingaddressid, useremail, originalid, origin, currency, shippingmethodid, collectionpointid, shippingmethodname, collectionpointname, channelid, shippingpricenetamount, shippingpricegrossamount, shippingtaxrate, token, checkouttoken, totalnetamount, undiscountedtotalnetamount, totalgrossamount, undiscountedtotalgrossamount, totalpaidamount, voucherid, displaygrossprices, customernote, weightamount, weightunit, redirecturl, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: pages; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.pages (id, title, slug, pagetypeid, content, createat, metadata, privatemetadata, publicationdate, ispublished, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: pagetranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.pagetranslations (id, languagecode, pageid, title, content, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: pagetypes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.pagetypes (id, name, slug, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.payments (id, gateway, isactive, toconfirm, createat, updateat, chargestatus, token, total, capturedamount, currency, checkoutid, orderid, billingemail, billingfirstname, billinglastname, billingcompanyname, billingaddress1, billingaddress2, billingcity, billingcityarea, billingpostalcode, billingcountrycode, billingcountryarea, ccfirstdigits, cclastdigits, ccbrand, ccexpmonth, ccexpyear, paymentmethodtype, customeripaddress, extradata, returnurl, pspreference, storepaymentmethod, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: pluginconfigurations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.pluginconfigurations (id, identifier, name, channelid, description, active, configuration) FROM stdin;
\.


--
-- Data for Name: pluginkeyvaluestore; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.pluginkeyvaluestore (pluginid, pkey, pvalue, expireat) FROM stdin;
\.


--
-- Data for Name: preferences; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.preferences (userid, category, name, value) FROM stdin;
\.


--
-- Data for Name: preorderallocations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.preorderallocations (id, orderlineid, quantity, productvariantchannellistingid) FROM stdin;
\.


--
-- Data for Name: productchannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productchannellistings (id, productid, channelid, visibleinlistings, availableforpurchase, currency, discountedpriceamount, createat, publicationdate, ispublished) FROM stdin;
\.


--
-- Data for Name: productcollections; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productcollections (id, collectionid, productid) FROM stdin;
\.


--
-- Data for Name: productmedias; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productmedias (id, createat, productid, ppoi, image, alt, type, externalurl, oembeddata, sortorder) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.products (id, producttypeid, name, slug, description, descriptionplaintext, categoryid, createat, updateat, chargetaxes, weight, weightunit, defaultvariantid, rating, metadata, privatemetadata, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: producttranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.producttranslations (id, languagecode, productid, name, description, seotitle, seodescription) FROM stdin;
\.


--
-- Data for Name: producttypes; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.producttypes (id, name, slug, kind, hasvariants, isshippingrequired, isdigital, weight, weightunit, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: productvariantchannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productvariantchannellistings (id, variantid, channelid, currency, priceamount, costpriceamount, preorderquantitythreshold, createat) FROM stdin;
\.


--
-- Data for Name: productvariants; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productvariants (id, name, productid, sku, weight, weightunit, trackinventory, ispreorder, preorderenddate, preorderglobalthreshold, sortorder, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: productvarianttranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.productvarianttranslations (id, languagecode, productvariantid, name) FROM stdin;
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.roles (id, name, displayname, description, createat, updateat, deleteat, permissionsstr, schememanaged, builtin) FROM stdin;
c4371d3c-3392-47c8-9479-7ddc43e391f6	system_post_all	authentication.roles.system_post_all.name	authentication.roles.system_post_all.description	1660901317607	1660901318050	0	create_post	f	t
530f394e-134c-458a-b7d1-f5bbef5d8366	system_read_only_admin	authentication.roles.system_read_only_admin.name	authentication.roles.system_read_only_admin.description	1660901317610	1660901318052	0	sysconsole_read_integrations_gif sysconsole_read_experimental_bleve sysconsole_read_authentication_email sysconsole_read_environment_file_storage sysconsole_read_authentication_mfa sysconsole_read_authentication_ldap sysconsole_read_integrations_integration_management sysconsole_read_site_public_links sysconsole_read_compliance_compliance_monitoring sysconsole_read_environment_high_availability download_compliance_export_result sysconsole_read_authentication_guest_access sysconsole_read_reporting_server_logs read_ldap_sync_job sysconsole_read_user_management_users sysconsole_read_integrations_bot_accounts sysconsole_read_environment_database sysconsole_read_site_announcement_banner sysconsole_read_site_file_sharing_and_downloads sysconsole_read_environment_image_proxy sysconsole_read_authentication_password sysconsole_read_site_notices sysconsole_read_compliance_custom_terms_of_service sysconsole_read_user_management_permissions sysconsole_read_environment_elasticsearch sysconsole_read_environment_web_server sysconsole_read_user_management_groups sysconsole_read_authentication_saml sysconsole_read_site_localization sysconsole_read_environment_performance_monitoring sysconsole_read_site_notifications get_analytics sysconsole_read_environment_developer sysconsole_read_compliance_compliance_export read_compliance_export_job read_audits sysconsole_read_authentication_signup sysconsole_read_reporting_site_statistics sysconsole_read_environment_push_notification_server sysconsole_read_experimental_feature_flags sysconsole_read_plugins read_elasticsearch_post_aggregation_job sysconsole_read_site_posts test_ldap sysconsole_read_environment_rate_limiting sysconsole_read_integrations_cors read_data_retention_job sysconsole_read_compliance_data_retention_policy sysconsole_read_experimental_features sysconsole_read_site_customization sysconsole_read_environment_smtp sysconsole_read_environment_logging read_elasticsearch_post_indexing_job sysconsole_read_authentication_openid get_logs sysconsole_read_environment_session_lengths	f	t
154cc966-89bf-45c3-9da6-6697d94b2bcd	system_user	authentication.roles.global_user.name	authentication.roles.global_user.description	1660901317594	1660901318037	0	view_members	t	t
51ff73a6-db0a-4d01-abc1-3c54f7db1024	system_post_all_public	authentication.roles.system_post_all_public.name	authentication.roles.system_post_all_public.description	1660901317597	1660901318040	0	create_post_public	f	t
f38592b8-b07b-45c6-bbeb-6536c6598954	system_user_access_token	authentication.roles.system_user_access_token.name	authentication.roles.system_user_access_token.description	1660901317599	1660901318043	0	revoke_user_access_token create_user_access_token read_user_access_token	f	t
d526f343-720d-4108-b4cb-0ed6a4f617b2	system_user_manager	authentication.roles.system_user_manager.name	authentication.roles.system_user_manager.description	1660901317602	1660901318045	0	sysconsole_read_authentication_saml sysconsole_read_authentication_mfa sysconsole_read_user_management_groups sysconsole_read_authentication_ldap read_ldap_sync_job sysconsole_read_user_management_permissions test_ldap sysconsole_read_authentication_email sysconsole_read_authentication_signup sysconsole_read_authentication_guest_access sysconsole_read_authentication_password sysconsole_read_authentication_openid sysconsole_write_user_management_groups	f	t
48d0f0ee-25e9-4696-ad6c-711fab68e9da	system_guest	authentication.roles.global_guest.name	authentication.roles.global_guest.description	1660901317605	1660901318048	0		t	t
4d2aa2d1-2d56-47d5-a0ca-bd04fd373af9	system_admin	authentication.roles.global_admin.name	authentication.roles.global_admin.description	1660901317616	1660901318058	0	sysconsole_read_environment_high_availability sysconsole_read_site_localization sysconsole_read_compliance_custom_terms_of_service manage_others_outgoing_webhooks edit_brand sysconsole_read_site_posts sysconsole_read_authentication_openid sysconsole_write_integrations_gif remove_saml_private_cert manage_oauth sysconsole_read_compliance_data_retention_policy add_ldap_private_cert edit_others_posts sysconsole_read_authentication_guest_access manage_jobs invalidate_email_invite create_data_retention_job sysconsole_read_environment_developer test_site_url read_audits sysconsole_write_environment_elasticsearch sysconsole_read_authentication_ldap manage_gift_card sysconsole_read_environment_session_lengths manage_orders sysconsole_write_site_posts read_elasticsearch_post_indexing_job sysconsole_read_experimental_feature_flags sysconsole_write_authentication_ldap manage_staff manage_page_types_and_attributes sysconsole_write_environment_high_availability sysconsole_write_site_announcement_banner sysconsole_read_authentication_saml manage_discounts sysconsole_write_environment_image_proxy sysconsole_read_user_management_users remove_saml_idp_cert sysconsole_read_environment_file_storage sysconsole_write_user_management_system_roles sysconsole_read_billing manage_pages delete_others_posts sysconsole_read_site_notifications sysconsole_write_compliance_data_retention_policy manage_channels test_s3 sysconsole_write_reporting_server_logs sysconsole_read_environment_web_server create_compliance_export_job revoke_user_access_token sysconsole_read_site_customization upload_file test_ldap get_public_link sysconsole_read_integrations_cors sysconsole_write_experimental_bleve sysconsole_read_environment_smtp sysconsole_read_user_management_permissions remove_ldap_public_cert sysconsole_read_experimental_features sysconsole_read_authentication_password sysconsole_write_authentication_email manage_menus create_elasticsearch_post_aggregation_job sysconsole_write_site_localization sysconsole_read_compliance_compliance_monitoring sysconsole_write_integrations_integration_management sysconsole_write_site_customization create_ldap_sync_job sysconsole_write_environment_performance_monitoring sysconsole_read_environment_elasticsearch sysconsole_write_experimental_features manage_apps sysconsole_write_user_management_users sysconsole_read_reporting_server_logs sysconsole_write_environment_push_notification_server sysconsole_write_site_notices manage_system sysconsole_write_environment_file_storage add_reaction read_jobs remove_ldap_private_cert remove_saml_public_cert sysconsole_write_environment_logging recycle_database_connections add_saml_public_cert sysconsole_read_integrations_integration_management sysconsole_write_integrations_cors sysconsole_write_billing sysconsole_read_integrations_gif manage_checkouts sysconsole_read_environment_image_proxy manage_remote_clusters sysconsole_write_environment_smtp sysconsole_write_environment_session_lengths read_data_retention_job manage_plugins sysconsole_read_environment_logging sysconsole_read_environment_rate_limiting sysconsole_write_compliance_compliance_export manage_users sysconsole_write_plugins sysconsole_read_environment_database create_user_access_token sysconsole_write_user_management_groups sysconsole_write_compliance_custom_terms_of_service sysconsole_write_site_file_sharing_and_downloads sysconsole_read_site_file_sharing_and_downloads manage_incoming_webhooks sysconsole_read_user_management_groups delete_post sysconsole_read_compliance_compliance_export edit_other_users create_post_bleve_indexes_job sysconsole_read_integrations_bot_accounts sysconsole_write_environment_database get_analytics sysconsole_read_authentication_email sysconsole_write_experimental_feature_flags sysconsole_write_authentication_guest_access test_elasticsearch sysconsole_read_site_notices test_email manage_shipping manage_product_types_and_attributes manage_outgoing_webhooks get_saml_cert_status sysconsole_read_site_announcement_banner add_saml_private_cert manage_translations sysconsole_write_environment_developer sysconsole_write_authentication_password read_elasticsearch_post_aggregation_job manage_roles sysconsole_write_environment_web_server sysconsole_read_user_management_system_roles sysconsole_read_site_public_links download_compliance_export_result sysconsole_read_environment_push_notification_server purge_bleve_indexes sysconsole_write_reporting_site_statistics sysconsole_write_user_management_permissions manage_settings reload_config create_post view_members sysconsole_write_site_notifications invalidate_caches read_user_access_token sysconsole_read_environment_performance_monitoring sysconsole_read_authentication_signup sysconsole_read_experimental_bleve read_ldap_sync_job add_ldap_public_cert read_compliance_export_job add_saml_idp_cert remove_reaction sysconsole_read_reporting_site_statistics invite_user get_saml_metadata_from_idp sysconsole_write_environment_rate_limiting manage_products create_elasticsearch_post_indexing_job sysconsole_write_authentication_openid purge_elasticsearch_indexes sysconsole_write_compliance_compliance_monitoring sysconsole_write_site_public_links sysconsole_write_authentication_saml edit_post sysconsole_read_authentication_mfa impersonate_user sysconsole_read_plugins sysconsole_write_authentication_mfa create_post_public manage_others_incoming_webhooks sysconsole_write_authentication_signup manage_system_wide_oauth get_logs sysconsole_write_integrations_bot_accounts create_post_ephemeral	t	t
4b0943be-e029-4cb2-9a68-cf2250bdfdbb	system_manager	authentication.roles.system_manager.name	authentication.roles.system_manager.description	1660901317613	1660901318055	0	sysconsole_read_integrations_gif sysconsole_read_environment_rate_limiting sysconsole_write_environment_session_lengths sysconsole_read_authentication_mfa sysconsole_read_site_public_links sysconsole_write_environment_rate_limiting read_elasticsearch_post_indexing_job create_elasticsearch_post_aggregation_job sysconsole_write_user_management_groups sysconsole_read_environment_elasticsearch sysconsole_read_authentication_openid sysconsole_write_site_announcement_banner sysconsole_write_environment_image_proxy sysconsole_read_integrations_integration_management sysconsole_read_environment_image_proxy invalidate_caches sysconsole_write_environment_developer test_email sysconsole_write_integrations_integration_management sysconsole_read_reporting_site_statistics recycle_database_connections sysconsole_read_environment_developer sysconsole_write_environment_elasticsearch reload_config sysconsole_write_site_notices sysconsole_read_authentication_saml sysconsole_write_site_customization sysconsole_read_user_management_groups test_site_url sysconsole_read_environment_file_storage sysconsole_read_environment_logging sysconsole_write_environment_database sysconsole_read_environment_web_server sysconsole_read_environment_smtp read_ldap_sync_job sysconsole_read_environment_database sysconsole_read_authentication_password sysconsole_write_environment_push_notification_server sysconsole_write_environment_web_server sysconsole_read_site_announcement_banner sysconsole_write_integrations_cors test_s3 sysconsole_read_environment_push_notification_server sysconsole_write_integrations_gif sysconsole_read_authentication_signup edit_brand sysconsole_read_environment_session_lengths sysconsole_write_user_management_permissions sysconsole_read_authentication_email sysconsole_read_site_customization sysconsole_read_integrations_cors sysconsole_read_site_file_sharing_and_downloads purge_elasticsearch_indexes sysconsole_read_plugins sysconsole_read_environment_high_availability sysconsole_write_site_public_links sysconsole_write_site_posts test_elasticsearch sysconsole_read_authentication_ldap sysconsole_write_environment_performance_monitoring sysconsole_write_environment_smtp get_logs sysconsole_read_site_notifications create_elasticsearch_post_indexing_job sysconsole_read_reporting_server_logs sysconsole_write_environment_logging sysconsole_read_environment_performance_monitoring sysconsole_write_site_notifications sysconsole_write_environment_high_availability sysconsole_read_site_notices read_elasticsearch_post_aggregation_job sysconsole_write_integrations_bot_accounts get_analytics sysconsole_read_integrations_bot_accounts sysconsole_read_authentication_guest_access sysconsole_read_site_localization sysconsole_write_site_localization sysconsole_write_environment_file_storage sysconsole_read_site_posts sysconsole_write_site_file_sharing_and_downloads test_ldap sysconsole_read_user_management_permissions	f	t
\.


--
-- Data for Name: salecategories; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.salecategories (id, saleid, categoryid, createat) FROM stdin;
\.


--
-- Data for Name: salechannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.salechannellistings (id, saleid, channelid, discountvalue, currency, createat) FROM stdin;
\.


--
-- Data for Name: salecollections; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.salecollections (id, saleid, collectionid, createat) FROM stdin;
\.


--
-- Data for Name: saleproducts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.saleproducts (id, saleid, productid, createat) FROM stdin;
\.


--
-- Data for Name: saleproductvariants; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.saleproductvariants (id, saleid, productvariantid, createat) FROM stdin;
\.


--
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.sales (id, shopid, name, type, startdate, enddate, createat, updateat, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: saletranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.saletranslations (id, languagecode, name, saleid) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.sessions (id, token, createat, expiresat, lastactivityat, userid, deviceid, roles, isoauth, expirednotify, props) FROM stdin;
\.


--
-- Data for Name: shippingmethodchannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingmethodchannellistings (id, shippingmethodid, channelid, minimumorderpriceamount, currency, maximumorderpriceamount, priceamount, createat) FROM stdin;
\.


--
-- Data for Name: shippingmethodexcludedproducts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingmethodexcludedproducts (id, shippingmethodid, productid) FROM stdin;
\.


--
-- Data for Name: shippingmethodpostalcoderules; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingmethodpostalcoderules (id, shippingmethodid, start, "end", inclusiontype) FROM stdin;
\.


--
-- Data for Name: shippingmethods; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingmethods (id, name, type, shippingzoneid, minimumorderweight, maximumorderweight, weightunit, maximumdeliverydays, minimumdeliverydays, description, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: shippingmethodtranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingmethodtranslations (id, shippingmethodid, languagecode, name, description) FROM stdin;
\.


--
-- Data for Name: shippingzonechannels; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingzonechannels (id, shippingzoneid, channelid) FROM stdin;
\.


--
-- Data for Name: shippingzones; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shippingzones (id, name, countries, "default", description, createat, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: shops; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shops (id, ownerid, createat, updateat, name, description, topmenuid, includetaxesinprice, displaygrossprices, chargetaxesonshipping, trackinventorybydefault, defaultweightunit, automaticfulfillmentdigitalproducts, defaultdigitalmaxdownloads, defaultdigitalurlvaliddays, addressid, defaultmailsendername, defaultmailsenderaddress, customersetpasswordurl, automaticallyconfirmallneworders, fulfillmentautoapprove, fulfillmentallowunpaid, giftcardexpirytype, giftcardexpiryperiodtype, giftcardexpiryperiod, automaticallyfulfillnonshippablegiftcard) FROM stdin;
\.


--
-- Data for Name: shopstaffs; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shopstaffs (id, shopid, staffid, createat, endat) FROM stdin;
\.


--
-- Data for Name: shoptranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.shoptranslations (id, shopid, languagecode, name, description, createat, updateat) FROM stdin;
\.


--
-- Data for Name: staffnotificationrecipients; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.staffnotificationrecipients (id, userid, staffemail, active) FROM stdin;
\.


--
-- Data for Name: status; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.status (userid, status, manual, lastactivityat) FROM stdin;
\.


--
-- Data for Name: stocks; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.stocks (id, createat, warehouseid, productvariantid, quantity) FROM stdin;
\.


--
-- Data for Name: systems; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.systems (name, value) FROM stdin;
AsymmetricSigningKey	{"ecdsa_key":{"curve":"P-256","x":60994663154940964224851498999282045830503865022427608004163590976762732432287,"y":79992441428556457832805863823092740744010174656252171566363834367025350740685,"d":106032714709113018694787692407169974053106921122688902321141744677825111735358}}
PostActionCookieSecret	{"key":"qO+CSLg5jLZdV/KXZoMIY3Lt7kibmlMWs4SeaD7e7lY="}
InstallationDate	1660901316994
FirstServerRunTimestamp	1660901316996
LastSecurityTime	1660901317004
AdvancedPermissionsMigrationComplete	true
SystemConsoleRolesCreationMigrationComplete	true
webhook_permissions_split	true
remove_permanent_delete_user	true
view_members_new_permission	true
add_system_console_permissions	true
manage_secure_connections_permissions	true
add_system_roles_permissions	true
add_billing_permissions	true
download_compliance_export_results	true
experimental_subsection_permissions	true
authentication_subsection_permissions	true
integrations_subsection_permissions	true
site_subsection_permissions	true
compliance_subsection_permissions	true
environment_subsection_permissions	true
reporting_subsection_permissions	true
test_email_ancillary_permission	true
ContentExtractionConfigDefaultTrueMigrationComplete	true
\.


--
-- Data for Name: termsofservices; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.termsofservices (id, createat, userid, text) FROM stdin;
\.


--
-- Data for Name: tokens; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.tokens (token, createat, type, extra) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.transactions (id, createat, paymentid, token, kind, issuccess, actionrequired, actionrequireddata, currency, amount, error, customerid, gatewayresponse, alreadyprocessed) FROM stdin;
\.


--
-- Data for Name: uploadsessions; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.uploadsessions (id, type, createat, userid, filename, path, filesize, fileoffset) FROM stdin;
\.


--
-- Data for Name: useraccesstokens; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.useraccesstokens (id, token, userid, description, isactive) FROM stdin;
\.


--
-- Data for Name: useraddresses; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.useraddresses (id, userid, addressid) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.users (id, email, username, firstname, lastname, defaultshippingaddressid, defaultbillingaddressid, password, authdata, authservice, emailverified, nickname, roles, props, notifyprops, lastpasswordupdate, lastpictureupdate, failedattempts, locale, timezone, mfaactive, mfasecret, createat, updateat, deleteat, isactive, note, jwttokenkey, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: usertermofservices; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.usertermofservices (userid, termsofserviceid, createat) FROM stdin;
\.


--
-- Data for Name: variantmedias; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.variantmedias (id, variantid, mediaid) FROM stdin;
\.


--
-- Data for Name: vouchercategories; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.vouchercategories (id, voucherid, categoryid, createat) FROM stdin;
\.


--
-- Data for Name: voucherchannellistings; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.voucherchannellistings (id, createat, voucherid, channelid, discountvalue, currency, minspenamount) FROM stdin;
\.


--
-- Data for Name: vouchercollections; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.vouchercollections (id, voucherid, collectionid) FROM stdin;
\.


--
-- Data for Name: vouchercustomers; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.vouchercustomers (id, voucherid, customeremail) FROM stdin;
\.


--
-- Data for Name: voucherproducts; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.voucherproducts (id, voucherid, productid) FROM stdin;
\.


--
-- Data for Name: voucherproductvariants; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.voucherproductvariants (id, voucherid, productvariantid, createat) FROM stdin;
\.


--
-- Data for Name: vouchers; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.vouchers (id, shopid, type, name, code, usagelimit, used, startdate, enddate, applyonceperorder, applyoncepercustomer, onlyforstaff, discountvaluetype, countries, mincheckoutitemsquantity, createat, updateat, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: vouchertranslations; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.vouchertranslations (id, languagecode, name, voucherid, createat) FROM stdin;
\.


--
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.warehouses (id, name, slug, addressid, email, clickandcollectoption, isprivate, metadata, privatemetadata) FROM stdin;
\.


--
-- Data for Name: warehouseshippingzones; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.warehouseshippingzones (id, warehouseid, shippingzoneid) FROM stdin;
\.


--
-- Data for Name: wishlistitemproductvariants; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.wishlistitemproductvariants (id, wishlistitemid, productvariantid) FROM stdin;
\.


--
-- Data for Name: wishlistitems; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.wishlistitems (id, wishlistid, productid, createat) FROM stdin;
\.


--
-- Data for Name: wishlists; Type: TABLE DATA; Schema: public; Owner: minh
--

COPY public.wishlists (id, token, userid, createat) FROM stdin;
\.


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: allocations allocations_orderlineid_stockid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT allocations_orderlineid_stockid_key UNIQUE (orderlineid, stockid);


--
-- Name: allocations allocations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT allocations_pkey PRIMARY KEY (id);


--
-- Name: apps apps_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.apps
    ADD CONSTRAINT apps_pkey PRIMARY KEY (id);


--
-- Name: apptokens apptokens_authtoken_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.apptokens
    ADD CONSTRAINT apptokens_authtoken_key UNIQUE (authtoken);


--
-- Name: apptokens apptokens_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.apptokens
    ADD CONSTRAINT apptokens_pkey PRIMARY KEY (id);


--
-- Name: assignedpageattributes assignedpageattributes_pageid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributes
    ADD CONSTRAINT assignedpageattributes_pageid_assignmentid_key UNIQUE (pageid, assignmentid);


--
-- Name: assignedpageattributes assignedpageattributes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributes
    ADD CONSTRAINT assignedpageattributes_pkey PRIMARY KEY (id);


--
-- Name: assignedpageattributevalues assignedpageattributevalues_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributevalues
    ADD CONSTRAINT assignedpageattributevalues_pkey PRIMARY KEY (id);


--
-- Name: assignedpageattributevalues assignedpageattributevalues_valueid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributevalues
    ADD CONSTRAINT assignedpageattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);


--
-- Name: assignedproductattributes assignedproductattributes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributes
    ADD CONSTRAINT assignedproductattributes_pkey PRIMARY KEY (id);


--
-- Name: assignedproductattributes assignedproductattributes_productid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributes
    ADD CONSTRAINT assignedproductattributes_productid_assignmentid_key UNIQUE (productid, assignmentid);


--
-- Name: assignedproductattributevalues assignedproductattributevalues_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributevalues
    ADD CONSTRAINT assignedproductattributevalues_pkey PRIMARY KEY (id);


--
-- Name: assignedproductattributevalues assignedproductattributevalues_valueid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributevalues
    ADD CONSTRAINT assignedproductattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);


--
-- Name: assignedvariantattributes assignedvariantattributes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributes
    ADD CONSTRAINT assignedvariantattributes_pkey PRIMARY KEY (id);


--
-- Name: assignedvariantattributes assignedvariantattributes_variantid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributes
    ADD CONSTRAINT assignedvariantattributes_variantid_assignmentid_key UNIQUE (variantid, assignmentid);


--
-- Name: assignedvariantattributevalues assignedvariantattributevalues_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributevalues
    ADD CONSTRAINT assignedvariantattributevalues_pkey PRIMARY KEY (id);


--
-- Name: assignedvariantattributevalues assignedvariantattributevalues_valueid_assignmentid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributevalues
    ADD CONSTRAINT assignedvariantattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);


--
-- Name: attributepages attributepages_attributeid_pagetypeid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributepages
    ADD CONSTRAINT attributepages_attributeid_pagetypeid_key UNIQUE (attributeid, pagetypeid);


--
-- Name: attributepages attributepages_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributepages
    ADD CONSTRAINT attributepages_pkey PRIMARY KEY (id);


--
-- Name: attributeproducts attributeproducts_attributeid_producttypeid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributeproducts
    ADD CONSTRAINT attributeproducts_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);


--
-- Name: attributeproducts attributeproducts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributeproducts
    ADD CONSTRAINT attributeproducts_pkey PRIMARY KEY (id);


--
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- Name: attributes attributes_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_slug_key UNIQUE (slug);


--
-- Name: attributetranslations attributetranslations_languagecode_attributeid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributetranslations
    ADD CONSTRAINT attributetranslations_languagecode_attributeid_key UNIQUE (languagecode, attributeid);


--
-- Name: attributetranslations attributetranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributetranslations
    ADD CONSTRAINT attributetranslations_pkey PRIMARY KEY (id);


--
-- Name: attributevalues attributevalues_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevalues
    ADD CONSTRAINT attributevalues_pkey PRIMARY KEY (id);


--
-- Name: attributevalues attributevalues_slug_attributeid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevalues
    ADD CONSTRAINT attributevalues_slug_attributeid_key UNIQUE (slug, attributeid);


--
-- Name: attributevaluetranslations attributevaluetranslations_languagecode_attributevalueid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevaluetranslations
    ADD CONSTRAINT attributevaluetranslations_languagecode_attributevalueid_key UNIQUE (languagecode, attributevalueid);


--
-- Name: attributevaluetranslations attributevaluetranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevaluetranslations
    ADD CONSTRAINT attributevaluetranslations_pkey PRIMARY KEY (id);


--
-- Name: attributevariants attributevariants_attributeid_producttypeid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevariants
    ADD CONSTRAINT attributevariants_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);


--
-- Name: attributevariants attributevariants_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevariants
    ADD CONSTRAINT attributevariants_pkey PRIMARY KEY (id);


--
-- Name: audits audits_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.audits
    ADD CONSTRAINT audits_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categories_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);


--
-- Name: categorytranslations categorytranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.categorytranslations
    ADD CONSTRAINT categorytranslations_pkey PRIMARY KEY (id);


--
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- Name: channels channels_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_slug_key UNIQUE (slug);


--
-- Name: checkoutlines checkoutlines_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkoutlines
    ADD CONSTRAINT checkoutlines_pkey PRIMARY KEY (id);


--
-- Name: checkouts checkouts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT checkouts_pkey PRIMARY KEY (token);


--
-- Name: clusterdiscovery clusterdiscovery_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.clusterdiscovery
    ADD CONSTRAINT clusterdiscovery_pkey PRIMARY KEY (id);


--
-- Name: collectionchannellistings collectionchannellistings_collectionid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectionchannellistings
    ADD CONSTRAINT collectionchannellistings_collectionid_channelid_key UNIQUE (collectionid, channelid);


--
-- Name: collectionchannellistings collectionchannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectionchannellistings
    ADD CONSTRAINT collectionchannellistings_pkey PRIMARY KEY (id);


--
-- Name: collections collections_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collections
    ADD CONSTRAINT collections_pkey PRIMARY KEY (id);


--
-- Name: collectiontranslations collectiontranslations_languagecode_collectionid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectiontranslations
    ADD CONSTRAINT collectiontranslations_languagecode_collectionid_key UNIQUE (languagecode, collectionid);


--
-- Name: collectiontranslations collectiontranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectiontranslations
    ADD CONSTRAINT collectiontranslations_pkey PRIMARY KEY (id);


--
-- Name: compliances compliances_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.compliances
    ADD CONSTRAINT compliances_pkey PRIMARY KEY (id);


--
-- Name: customerevents customerevents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.customerevents
    ADD CONSTRAINT customerevents_pkey PRIMARY KEY (id);


--
-- Name: db_lock db_lock_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.db_lock
    ADD CONSTRAINT db_lock_pkey PRIMARY KEY (id);


--
-- Name: db_migrations db_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.db_migrations
    ADD CONSTRAINT db_migrations_pkey PRIMARY KEY (version);


--
-- Name: digitalcontents digitalcontents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontents
    ADD CONSTRAINT digitalcontents_pkey PRIMARY KEY (id);


--
-- Name: digitalcontenturls digitalcontenturls_lineid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_lineid_key UNIQUE (lineid);


--
-- Name: digitalcontenturls digitalcontenturls_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_pkey PRIMARY KEY (id);


--
-- Name: digitalcontenturls digitalcontenturls_token_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontenturls
    ADD CONSTRAINT digitalcontenturls_token_key UNIQUE (token);


--
-- Name: exportevents exportevents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.exportevents
    ADD CONSTRAINT exportevents_pkey PRIMARY KEY (id);


--
-- Name: exportfiles exportfiles_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.exportfiles
    ADD CONSTRAINT exportfiles_pkey PRIMARY KEY (id);


--
-- Name: fileinfos fileinfos_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fileinfos
    ADD CONSTRAINT fileinfos_pkey PRIMARY KEY (id);


--
-- Name: fulfillmentlines fulfillmentlines_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillmentlines
    ADD CONSTRAINT fulfillmentlines_pkey PRIMARY KEY (id);


--
-- Name: fulfillments fulfillments_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillments
    ADD CONSTRAINT fulfillments_pkey PRIMARY KEY (id);


--
-- Name: giftcardcheckouts giftcardcheckouts_giftcardid_checkoutid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardcheckouts
    ADD CONSTRAINT giftcardcheckouts_giftcardid_checkoutid_key UNIQUE (giftcardid, checkoutid);


--
-- Name: giftcardcheckouts giftcardcheckouts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardcheckouts
    ADD CONSTRAINT giftcardcheckouts_pkey PRIMARY KEY (id);


--
-- Name: giftcardevents giftcardevents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardevents
    ADD CONSTRAINT giftcardevents_pkey PRIMARY KEY (id);


--
-- Name: giftcards giftcards_code_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcards
    ADD CONSTRAINT giftcards_code_key UNIQUE (code);


--
-- Name: giftcards giftcards_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcards
    ADD CONSTRAINT giftcards_pkey PRIMARY KEY (id);


--
-- Name: invoiceevents invoiceevents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoiceevents
    ADD CONSTRAINT invoiceevents_pkey PRIMARY KEY (id);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--
-- Name: menuitems menuitems_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT menuitems_pkey PRIMARY KEY (id);


--
-- Name: menuitemtranslations menuitemtranslations_languagecode_menuitemid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitemtranslations
    ADD CONSTRAINT menuitemtranslations_languagecode_menuitemid_key UNIQUE (languagecode, menuitemid);


--
-- Name: menuitemtranslations menuitemtranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitemtranslations
    ADD CONSTRAINT menuitemtranslations_pkey PRIMARY KEY (id);


--
-- Name: menus menus_name_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_name_key UNIQUE (name);


--
-- Name: menus menus_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_pkey PRIMARY KEY (id);


--
-- Name: menus menus_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_slug_key UNIQUE (slug);


--
-- Name: openexchangerates openexchangerates_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.openexchangerates
    ADD CONSTRAINT openexchangerates_pkey PRIMARY KEY (id);


--
-- Name: openexchangerates openexchangerates_tocurrency_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.openexchangerates
    ADD CONSTRAINT openexchangerates_tocurrency_key UNIQUE (tocurrency);


--
-- Name: orderdiscounts orderdiscounts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderdiscounts
    ADD CONSTRAINT orderdiscounts_pkey PRIMARY KEY (id);


--
-- Name: orderevents orderevents_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderevents
    ADD CONSTRAINT orderevents_pkey PRIMARY KEY (id);


--
-- Name: ordergiftcards ordergiftcards_giftcardid_orderid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.ordergiftcards
    ADD CONSTRAINT ordergiftcards_giftcardid_orderid_key UNIQUE (giftcardid, orderid);


--
-- Name: ordergiftcards ordergiftcards_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.ordergiftcards
    ADD CONSTRAINT ordergiftcards_pkey PRIMARY KEY (id);


--
-- Name: orderlines orderlines_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderlines
    ADD CONSTRAINT orderlines_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: orders orders_token_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_token_key UNIQUE (token);


--
-- Name: pages pages_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pages
    ADD CONSTRAINT pages_pkey PRIMARY KEY (id);


--
-- Name: pages pages_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pages
    ADD CONSTRAINT pages_slug_key UNIQUE (slug);


--
-- Name: pagetranslations pagetranslations_languagecode_pageid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pagetranslations
    ADD CONSTRAINT pagetranslations_languagecode_pageid_key UNIQUE (languagecode, pageid);


--
-- Name: pagetranslations pagetranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pagetranslations
    ADD CONSTRAINT pagetranslations_pkey PRIMARY KEY (id);


--
-- Name: pagetypes pagetypes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pagetypes
    ADD CONSTRAINT pagetypes_pkey PRIMARY KEY (id);


--
-- Name: pagetypes pagetypes_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pagetypes
    ADD CONSTRAINT pagetypes_slug_key UNIQUE (slug);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: pluginconfigurations pluginconfigurations_identifier_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pluginconfigurations
    ADD CONSTRAINT pluginconfigurations_identifier_channelid_key UNIQUE (identifier, channelid);


--
-- Name: pluginconfigurations pluginconfigurations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pluginconfigurations
    ADD CONSTRAINT pluginconfigurations_pkey PRIMARY KEY (id);


--
-- Name: pluginkeyvaluestore pluginkeyvaluestore_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pluginkeyvaluestore
    ADD CONSTRAINT pluginkeyvaluestore_pkey PRIMARY KEY (pluginid, pkey);


--
-- Name: preferences preferences_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.preferences
    ADD CONSTRAINT preferences_pkey PRIMARY KEY (userid, category, name);


--
-- Name: preorderallocations preorderallocations_orderlineid_productvariantchannellistin_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.preorderallocations
    ADD CONSTRAINT preorderallocations_orderlineid_productvariantchannellistin_key UNIQUE (orderlineid, productvariantchannellistingid);


--
-- Name: preorderallocations preorderallocations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.preorderallocations
    ADD CONSTRAINT preorderallocations_pkey PRIMARY KEY (id);


--
-- Name: productchannellistings productchannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productchannellistings
    ADD CONSTRAINT productchannellistings_pkey PRIMARY KEY (id);


--
-- Name: productchannellistings productchannellistings_productid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productchannellistings
    ADD CONSTRAINT productchannellistings_productid_channelid_key UNIQUE (productid, channelid);


--
-- Name: productcollections productcollections_collectionid_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productcollections
    ADD CONSTRAINT productcollections_collectionid_productid_key UNIQUE (collectionid, productid);


--
-- Name: productcollections productcollections_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productcollections
    ADD CONSTRAINT productcollections_pkey PRIMARY KEY (id);


--
-- Name: productmedias productmedias_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productmedias
    ADD CONSTRAINT productmedias_pkey PRIMARY KEY (id);


--
-- Name: products products_name_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_name_key UNIQUE (name);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: products products_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);


--
-- Name: producttranslations producttranslations_languagecode_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.producttranslations
    ADD CONSTRAINT producttranslations_languagecode_productid_key UNIQUE (languagecode, productid);


--
-- Name: producttranslations producttranslations_name_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.producttranslations
    ADD CONSTRAINT producttranslations_name_key UNIQUE (name);


--
-- Name: producttranslations producttranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.producttranslations
    ADD CONSTRAINT producttranslations_pkey PRIMARY KEY (id);


--
-- Name: producttypes producttypes_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.producttypes
    ADD CONSTRAINT producttypes_pkey PRIMARY KEY (id);


--
-- Name: productvariantchannellistings productvariantchannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariantchannellistings
    ADD CONSTRAINT productvariantchannellistings_pkey PRIMARY KEY (id);


--
-- Name: productvariantchannellistings productvariantchannellistings_variantid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariantchannellistings
    ADD CONSTRAINT productvariantchannellistings_variantid_channelid_key UNIQUE (variantid, channelid);


--
-- Name: productvariants productvariants_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariants
    ADD CONSTRAINT productvariants_pkey PRIMARY KEY (id);


--
-- Name: productvariants productvariants_sku_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariants
    ADD CONSTRAINT productvariants_sku_key UNIQUE (sku);


--
-- Name: productvarianttranslations productvarianttranslations_languagecode_productvariantid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvarianttranslations
    ADD CONSTRAINT productvarianttranslations_languagecode_productvariantid_key UNIQUE (languagecode, productvariantid);


--
-- Name: productvarianttranslations productvarianttranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvarianttranslations
    ADD CONSTRAINT productvarianttranslations_pkey PRIMARY KEY (id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: salecategories salecategories_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecategories
    ADD CONSTRAINT salecategories_pkey PRIMARY KEY (id);


--
-- Name: salecategories salecategories_saleid_categoryid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecategories
    ADD CONSTRAINT salecategories_saleid_categoryid_key UNIQUE (saleid, categoryid);


--
-- Name: salechannellistings salechannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salechannellistings
    ADD CONSTRAINT salechannellistings_pkey PRIMARY KEY (id);


--
-- Name: salechannellistings salechannellistings_saleid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salechannellistings
    ADD CONSTRAINT salechannellistings_saleid_channelid_key UNIQUE (saleid, channelid);


--
-- Name: salecollections salecollections_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecollections
    ADD CONSTRAINT salecollections_pkey PRIMARY KEY (id);


--
-- Name: salecollections salecollections_saleid_collectionid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecollections
    ADD CONSTRAINT salecollections_saleid_collectionid_key UNIQUE (saleid, collectionid);


--
-- Name: saleproducts saleproducts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproducts
    ADD CONSTRAINT saleproducts_pkey PRIMARY KEY (id);


--
-- Name: saleproducts saleproducts_saleid_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproducts
    ADD CONSTRAINT saleproducts_saleid_productid_key UNIQUE (saleid, productid);


--
-- Name: saleproductvariants saleproductvariants_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproductvariants
    ADD CONSTRAINT saleproductvariants_pkey PRIMARY KEY (id);


--
-- Name: saleproductvariants saleproductvariants_saleid_productvariantid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproductvariants
    ADD CONSTRAINT saleproductvariants_saleid_productvariantid_key UNIQUE (saleid, productvariantid);


--
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- Name: saletranslations saletranslations_languagecode_saleid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saletranslations
    ADD CONSTRAINT saletranslations_languagecode_saleid_key UNIQUE (languagecode, saleid);


--
-- Name: saletranslations saletranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saletranslations
    ADD CONSTRAINT saletranslations_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: shippingmethodchannellistings shippingmethodchannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodchannellistings
    ADD CONSTRAINT shippingmethodchannellistings_pkey PRIMARY KEY (id);


--
-- Name: shippingmethodchannellistings shippingmethodchannellistings_shippingmethodid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodchannellistings
    ADD CONSTRAINT shippingmethodchannellistings_shippingmethodid_channelid_key UNIQUE (shippingmethodid, channelid);


--
-- Name: shippingmethodexcludedproducts shippingmethodexcludedproducts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodexcludedproducts
    ADD CONSTRAINT shippingmethodexcludedproducts_pkey PRIMARY KEY (id);


--
-- Name: shippingmethodexcludedproducts shippingmethodexcludedproducts_shippingmethodid_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodexcludedproducts
    ADD CONSTRAINT shippingmethodexcludedproducts_shippingmethodid_productid_key UNIQUE (shippingmethodid, productid);


--
-- Name: shippingmethodpostalcoderules shippingmethodpostalcoderules_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodpostalcoderules
    ADD CONSTRAINT shippingmethodpostalcoderules_pkey PRIMARY KEY (id);


--
-- Name: shippingmethodpostalcoderules shippingmethodpostalcoderules_shippingmethodid_start_end_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodpostalcoderules
    ADD CONSTRAINT shippingmethodpostalcoderules_shippingmethodid_start_end_key UNIQUE (shippingmethodid, start, "end");


--
-- Name: shippingmethods shippingmethods_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethods
    ADD CONSTRAINT shippingmethods_pkey PRIMARY KEY (id);


--
-- Name: shippingmethodtranslations shippingmethodtranslations_languagecode_shippingmethodid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodtranslations
    ADD CONSTRAINT shippingmethodtranslations_languagecode_shippingmethodid_key UNIQUE (languagecode, shippingmethodid);


--
-- Name: shippingmethodtranslations shippingmethodtranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodtranslations
    ADD CONSTRAINT shippingmethodtranslations_pkey PRIMARY KEY (id);


--
-- Name: shippingzonechannels shippingzonechannels_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingzonechannels
    ADD CONSTRAINT shippingzonechannels_pkey PRIMARY KEY (id);


--
-- Name: shippingzonechannels shippingzonechannels_shippingzoneid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingzonechannels
    ADD CONSTRAINT shippingzonechannels_shippingzoneid_channelid_key UNIQUE (shippingzoneid, channelid);


--
-- Name: shippingzones shippingzones_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingzones
    ADD CONSTRAINT shippingzones_pkey PRIMARY KEY (id);


--
-- Name: shops shops_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT shops_pkey PRIMARY KEY (id);


--
-- Name: shopstaffs shopstaffs_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shopstaffs
    ADD CONSTRAINT shopstaffs_pkey PRIMARY KEY (id);


--
-- Name: shopstaffs shopstaffs_shopid_staffid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shopstaffs
    ADD CONSTRAINT shopstaffs_shopid_staffid_key UNIQUE (shopid, staffid);


--
-- Name: shoptranslations shoptranslations_languagecode_shopid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shoptranslations
    ADD CONSTRAINT shoptranslations_languagecode_shopid_key UNIQUE (languagecode, shopid);


--
-- Name: shoptranslations shoptranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shoptranslations
    ADD CONSTRAINT shoptranslations_pkey PRIMARY KEY (id);


--
-- Name: staffnotificationrecipients staffnotificationrecipients_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.staffnotificationrecipients
    ADD CONSTRAINT staffnotificationrecipients_pkey PRIMARY KEY (id);


--
-- Name: status status_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.status
    ADD CONSTRAINT status_pkey PRIMARY KEY (userid);


--
-- Name: stocks stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.stocks
    ADD CONSTRAINT stocks_pkey PRIMARY KEY (id);


--
-- Name: stocks stocks_warehouseid_productvariantid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.stocks
    ADD CONSTRAINT stocks_warehouseid_productvariantid_key UNIQUE (warehouseid, productvariantid);


--
-- Name: systems systems_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.systems
    ADD CONSTRAINT systems_pkey PRIMARY KEY (name);


--
-- Name: termsofservices termsofservices_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.termsofservices
    ADD CONSTRAINT termsofservices_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (token);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: uploadsessions uploadsessions_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.uploadsessions
    ADD CONSTRAINT uploadsessions_pkey PRIMARY KEY (id);


--
-- Name: useraccesstokens useraccesstokens_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraccesstokens
    ADD CONSTRAINT useraccesstokens_pkey PRIMARY KEY (id);


--
-- Name: useraccesstokens useraccesstokens_token_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraccesstokens
    ADD CONSTRAINT useraccesstokens_token_key UNIQUE (token);


--
-- Name: useraddresses useraddresses_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraddresses
    ADD CONSTRAINT useraddresses_pkey PRIMARY KEY (id);


--
-- Name: useraddresses useraddresses_userid_addressid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraddresses
    ADD CONSTRAINT useraddresses_userid_addressid_key UNIQUE (userid, addressid);


--
-- Name: users users_authdata_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_authdata_key UNIQUE (authdata);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: usertermofservices usertermofservices_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.usertermofservices
    ADD CONSTRAINT usertermofservices_pkey PRIMARY KEY (userid);


--
-- Name: variantmedias variantmedias_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.variantmedias
    ADD CONSTRAINT variantmedias_pkey PRIMARY KEY (id);


--
-- Name: variantmedias variantmedias_variantid_mediaid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.variantmedias
    ADD CONSTRAINT variantmedias_variantid_mediaid_key UNIQUE (variantid, mediaid);


--
-- Name: vouchercategories vouchercategories_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercategories
    ADD CONSTRAINT vouchercategories_pkey PRIMARY KEY (id);


--
-- Name: vouchercategories vouchercategories_voucherid_categoryid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercategories
    ADD CONSTRAINT vouchercategories_voucherid_categoryid_key UNIQUE (voucherid, categoryid);


--
-- Name: voucherchannellistings voucherchannellistings_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherchannellistings
    ADD CONSTRAINT voucherchannellistings_pkey PRIMARY KEY (id);


--
-- Name: voucherchannellistings voucherchannellistings_voucherid_channelid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherchannellistings
    ADD CONSTRAINT voucherchannellistings_voucherid_channelid_key UNIQUE (voucherid, channelid);


--
-- Name: vouchercollections vouchercollections_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercollections
    ADD CONSTRAINT vouchercollections_pkey PRIMARY KEY (id);


--
-- Name: vouchercollections vouchercollections_voucherid_collectionid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercollections
    ADD CONSTRAINT vouchercollections_voucherid_collectionid_key UNIQUE (voucherid, collectionid);


--
-- Name: vouchercustomers vouchercustomers_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercustomers
    ADD CONSTRAINT vouchercustomers_pkey PRIMARY KEY (id);


--
-- Name: vouchercustomers vouchercustomers_voucherid_customeremail_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercustomers
    ADD CONSTRAINT vouchercustomers_voucherid_customeremail_key UNIQUE (voucherid, customeremail);


--
-- Name: voucherproducts voucherproducts_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproducts
    ADD CONSTRAINT voucherproducts_pkey PRIMARY KEY (id);


--
-- Name: voucherproducts voucherproducts_voucherid_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproducts
    ADD CONSTRAINT voucherproducts_voucherid_productid_key UNIQUE (voucherid, productid);


--
-- Name: voucherproductvariants voucherproductvariants_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproductvariants
    ADD CONSTRAINT voucherproductvariants_pkey PRIMARY KEY (id);


--
-- Name: voucherproductvariants voucherproductvariants_voucherid_productvariantid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproductvariants
    ADD CONSTRAINT voucherproductvariants_voucherid_productvariantid_key UNIQUE (voucherid, productvariantid);


--
-- Name: vouchers vouchers_code_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchers
    ADD CONSTRAINT vouchers_code_key UNIQUE (code);


--
-- Name: vouchers vouchers_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchers
    ADD CONSTRAINT vouchers_pkey PRIMARY KEY (id);


--
-- Name: vouchertranslations vouchertranslations_languagecode_voucherid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchertranslations
    ADD CONSTRAINT vouchertranslations_languagecode_voucherid_key UNIQUE (languagecode, voucherid);


--
-- Name: vouchertranslations vouchertranslations_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchertranslations
    ADD CONSTRAINT vouchertranslations_pkey PRIMARY KEY (id);


--
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- Name: warehouses warehouses_slug_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_slug_key UNIQUE (slug);


--
-- Name: warehouseshippingzones warehouseshippingzones_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouseshippingzones
    ADD CONSTRAINT warehouseshippingzones_pkey PRIMARY KEY (id);


--
-- Name: warehouseshippingzones warehouseshippingzones_warehouseid_shippingzoneid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouseshippingzones
    ADD CONSTRAINT warehouseshippingzones_warehouseid_shippingzoneid_key UNIQUE (warehouseid, shippingzoneid);


--
-- Name: wishlistitemproductvariants wishlistitemproductvariants_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitemproductvariants
    ADD CONSTRAINT wishlistitemproductvariants_pkey PRIMARY KEY (id);


--
-- Name: wishlistitemproductvariants wishlistitemproductvariants_wishlistitemid_productvariantid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitemproductvariants
    ADD CONSTRAINT wishlistitemproductvariants_wishlistitemid_productvariantid_key UNIQUE (wishlistitemid, productvariantid);


--
-- Name: wishlistitems wishlistitems_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitems
    ADD CONSTRAINT wishlistitems_pkey PRIMARY KEY (id);


--
-- Name: wishlistitems wishlistitems_wishlistid_productid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitems
    ADD CONSTRAINT wishlistitems_wishlistid_productid_key UNIQUE (wishlistid, productid);


--
-- Name: wishlists wishlists_pkey; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_pkey PRIMARY KEY (id);


--
-- Name: wishlists wishlists_token_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_token_key UNIQUE (token);


--
-- Name: wishlists wishlists_userid_key; Type: CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_userid_key UNIQUE (userid);


--
-- Name: idx_address_city; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_city ON public.addresses USING btree (city);


--
-- Name: idx_address_country; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_country ON public.addresses USING btree (country);


--
-- Name: idx_address_firstname; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_firstname ON public.addresses USING btree (firstname);


--
-- Name: idx_address_firstname_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_firstname_lower_textpattern ON public.addresses USING btree (lower((firstname)::text) text_pattern_ops);


--
-- Name: idx_address_lastname; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_lastname ON public.addresses USING btree (lastname);


--
-- Name: idx_address_lastname_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_lastname_lower_textpattern ON public.addresses USING btree (lower((lastname)::text) text_pattern_ops);


--
-- Name: idx_address_phone; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_address_phone ON public.addresses USING btree (phone);


--
-- Name: idx_app_tokens_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_app_tokens_name ON public.apptokens USING btree (name);


--
-- Name: idx_app_tokens_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_app_tokens_name_lower_textpattern ON public.apptokens USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_apps_identifier; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_apps_identifier ON public.apps USING btree (identifier);


--
-- Name: idx_apps_identifier_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_apps_identifier_lower_textpattern ON public.apps USING btree (lower((identifier)::text) text_pattern_ops);


--
-- Name: idx_apps_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_apps_name ON public.apps USING btree (name);


--
-- Name: idx_apps_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_apps_name_lower_textpattern ON public.apps USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_attribute_value_translations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attribute_value_translations_name ON public.attributevaluetranslations USING btree (name);


--
-- Name: idx_attribute_value_translations_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attribute_value_translations_name_lower_textpattern ON public.attributevaluetranslations USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_attributes_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributes_name ON public.attributes USING btree (name);


--
-- Name: idx_attributes_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributes_name_lower_textpattern ON public.attributes USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_attributes_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributes_slug ON public.attributes USING btree (slug);


--
-- Name: idx_attributetranslations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributetranslations_name ON public.attributetranslations USING btree (name);


--
-- Name: idx_attributetranslations_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributetranslations_name_lower_textpattern ON public.attributetranslations USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_attributevalues_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributevalues_name ON public.attributevalues USING btree (name);


--
-- Name: idx_attributevalues_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributevalues_name_lower_textpattern ON public.attributevalues USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_attributevalues_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_attributevalues_slug ON public.attributevalues USING btree (slug);


--
-- Name: idx_audits_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_audits_user_id ON public.audits USING btree (userid);


--
-- Name: idx_category_translations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_category_translations_name ON public.categorytranslations USING btree (name);


--
-- Name: idx_category_translations_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_category_translations_name_lower_textpattern ON public.categorytranslations USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_channels_currency; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_channels_currency ON public.channels USING btree (currency);


--
-- Name: idx_channels_isactive; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_channels_isactive ON public.channels USING btree (isactive);


--
-- Name: idx_channels_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_channels_name ON public.channels USING btree (name);


--
-- Name: idx_channels_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_channels_name_lower_textpattern ON public.channels USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_channels_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_channels_slug ON public.channels USING btree (slug);


--
-- Name: idx_checkoutlines_checkout_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkoutlines_checkout_id ON public.checkoutlines USING btree (checkoutid);


--
-- Name: idx_checkoutlines_variant_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkoutlines_variant_id ON public.checkoutlines USING btree (variantid);


--
-- Name: idx_checkouts_billing_address_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_billing_address_id ON public.checkouts USING btree (billingaddressid);


--
-- Name: idx_checkouts_channelid; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_channelid ON public.checkouts USING btree (channelid);


--
-- Name: idx_checkouts_shipping_address_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_shipping_address_id ON public.checkouts USING btree (shippingaddressid);


--
-- Name: idx_checkouts_shipping_method_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_shipping_method_id ON public.checkouts USING btree (shippingmethodid);


--
-- Name: idx_checkouts_token; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_token ON public.checkouts USING btree (token);


--
-- Name: idx_checkouts_userid; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_checkouts_userid ON public.checkouts USING btree (userid);


--
-- Name: idx_collections_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_collections_name ON public.collections USING btree (name);


--
-- Name: idx_collections_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_collections_name_lower_textpattern ON public.collections USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_customer_notes_date; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_customer_notes_date ON public.customernotes USING btree (date);


--
-- Name: idx_fileinfo_content_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_content_txt ON public.fileinfos USING gin (to_tsvector('english'::regconfig, content));


--
-- Name: idx_fileinfo_create_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_create_at ON public.fileinfos USING btree (createat);


--
-- Name: idx_fileinfo_delete_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_delete_at ON public.fileinfos USING btree (deleteat);


--
-- Name: idx_fileinfo_extension_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_extension_at ON public.fileinfos USING btree (extension);


--
-- Name: idx_fileinfo_name_splitted; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_name_splitted ON public.fileinfos USING gin (to_tsvector('english'::regconfig, translate((name)::text, '.,-'::text, '   '::text)));


--
-- Name: idx_fileinfo_name_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_name_txt ON public.fileinfos USING gin (to_tsvector('english'::regconfig, (name)::text));


--
-- Name: idx_fileinfo_parent_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_parent_id ON public.fileinfos USING btree (parentid);


--
-- Name: idx_fileinfo_update_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fileinfo_update_at ON public.fileinfos USING btree (updateat);


--
-- Name: idx_fulfillments_status; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fulfillments_status ON public.fulfillments USING btree (status);


--
-- Name: idx_fulfillments_tracking_number; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_fulfillments_tracking_number ON public.fulfillments USING btree (trackingnumber);


--
-- Name: idx_giftcardevents_date; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_giftcardevents_date ON public.giftcardevents USING btree (date);


--
-- Name: idx_giftcards_code; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_giftcards_code ON public.giftcards USING btree (code);


--
-- Name: idx_giftcards_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_giftcards_metadata ON public.giftcards USING btree (metadata);


--
-- Name: idx_giftcards_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_giftcards_private_metadata ON public.giftcards USING btree (privatemetadata);


--
-- Name: idx_giftcards_tag; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_giftcards_tag ON public.giftcards USING btree (tag);


--
-- Name: idx_jobs_type; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_jobs_type ON public.jobs USING btree (type);


--
-- Name: idx_menus_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_menus_name ON public.menus USING btree (name);


--
-- Name: idx_menus_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_menus_name_lower_textpattern ON public.menus USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_menus_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_menus_slug ON public.menus USING btree (slug);


--
-- Name: idx_openexchange_to_currency; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_openexchange_to_currency ON public.openexchangerates USING btree (tocurrency);


--
-- Name: idx_order_discounts_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_discounts_name ON public.orderdiscounts USING btree (name);


--
-- Name: idx_order_discounts_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_discounts_name_lower_textpattern ON public.orderdiscounts USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_order_discounts_translated_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_discounts_translated_name ON public.orderdiscounts USING btree (translatedname);


--
-- Name: idx_order_discounts_translated_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_discounts_translated_name_lower_textpattern ON public.orderdiscounts USING btree (lower((translatedname)::text) text_pattern_ops);


--
-- Name: idx_order_lines_product_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_product_name ON public.orderlines USING btree (productname);


--
-- Name: idx_order_lines_product_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_product_name_lower_textpattern ON public.orderlines USING btree (lower((productname)::text) text_pattern_ops);


--
-- Name: idx_order_lines_translated_product_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_translated_product_name ON public.orderlines USING btree (translatedproductname);


--
-- Name: idx_order_lines_translated_variant_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_translated_variant_name ON public.orderlines USING btree (translatedvariantname);


--
-- Name: idx_order_lines_variant_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_variant_name ON public.orderlines USING btree (variantname);


--
-- Name: idx_order_lines_variant_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_order_lines_variant_name_lower_textpattern ON public.orderlines USING btree (lower((variantname)::text) text_pattern_ops);


--
-- Name: idx_orders_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_orders_metadata ON public.orders USING btree (metadata);


--
-- Name: idx_orders_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_orders_private_metadata ON public.orders USING btree (privatemetadata);


--
-- Name: idx_orders_user_email; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_orders_user_email ON public.orders USING btree (useremail);


--
-- Name: idx_orders_user_email_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_orders_user_email_lower_textpattern ON public.orders USING btree (lower((useremail)::text) text_pattern_ops);


--
-- Name: idx_page_types_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_page_types_name ON public.pagetypes USING btree (name);


--
-- Name: idx_page_types_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_page_types_name_lower_textpattern ON public.pagetypes USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_page_types_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_page_types_slug ON public.pagetypes USING btree (slug);


--
-- Name: idx_pages_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pages_metadata ON public.pages USING btree (metadata);


--
-- Name: idx_pages_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pages_private_metadata ON public.pages USING btree (privatemetadata);


--
-- Name: idx_pages_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pages_slug ON public.pages USING btree (slug);


--
-- Name: idx_pages_title; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pages_title ON public.pages USING btree (title);


--
-- Name: idx_pages_title_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pages_title_lower_textpattern ON public.pages USING btree (lower((title)::text) text_pattern_ops);


--
-- Name: idx_pagetypes_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pagetypes_metadata ON public.pagetypes USING btree (metadata);


--
-- Name: idx_pagetypes_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_pagetypes_private_metadata ON public.pagetypes USING btree (privatemetadata);


--
-- Name: idx_payments_charge_status; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_charge_status ON public.payments USING btree (chargestatus);


--
-- Name: idx_payments_is_active; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_is_active ON public.payments USING btree (isactive);


--
-- Name: idx_payments_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_metadata ON public.payments USING btree (metadata);


--
-- Name: idx_payments_order_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_order_id ON public.payments USING btree (orderid);


--
-- Name: idx_payments_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_private_metadata ON public.payments USING btree (privatemetadata);


--
-- Name: idx_payments_psp_reference; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_payments_psp_reference ON public.payments USING btree (pspreference);


--
-- Name: idx_plugin_configurations_identifier; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_plugin_configurations_identifier ON public.pluginconfigurations USING btree (identifier);


--
-- Name: idx_plugin_configurations_lower_textpattern_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_plugin_configurations_lower_textpattern_name ON public.pluginconfigurations USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_plugin_configurations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_plugin_configurations_name ON public.pluginconfigurations USING btree (name);


--
-- Name: idx_preferences_category; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_preferences_category ON public.preferences USING btree (category);


--
-- Name: idx_preferences_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_preferences_name ON public.preferences USING btree (name);


--
-- Name: idx_preferences_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_preferences_user_id ON public.preferences USING btree (userid);


--
-- Name: idx_product_types_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_product_types_name ON public.producttypes USING btree (name);


--
-- Name: idx_product_types_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_product_types_name_lower_textpattern ON public.producttypes USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_product_types_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_product_types_slug ON public.producttypes USING btree (slug);


--
-- Name: idx_product_variants_sku; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_product_variants_sku ON public.productvariants USING btree (sku);


--
-- Name: idx_productchannellistings_puplication_date; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_productchannellistings_puplication_date ON public.productchannellistings USING btree (publicationdate);


--
-- Name: idx_products_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_products_metadata ON public.products USING btree (metadata);


--
-- Name: idx_products_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_products_name ON public.products USING btree (name);


--
-- Name: idx_products_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_products_name_lower_textpattern ON public.products USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_products_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_products_private_metadata ON public.products USING btree (privatemetadata);


--
-- Name: idx_products_slug; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_products_slug ON public.products USING btree (slug);


--
-- Name: idx_sale_translations_language_code; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sale_translations_language_code ON public.saletranslations USING btree (languagecode);


--
-- Name: idx_sale_translations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sale_translations_name ON public.saletranslations USING btree (name);


--
-- Name: idx_sales_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sales_name ON public.sales USING btree (name);


--
-- Name: idx_sales_type; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sales_type ON public.sales USING btree (type);


--
-- Name: idx_sessions_create_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sessions_create_at ON public.sessions USING btree (createat);


--
-- Name: idx_sessions_expires_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sessions_expires_at ON public.sessions USING btree (expiresat);


--
-- Name: idx_sessions_last_activity_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sessions_last_activity_at ON public.sessions USING btree (lastactivityat);


--
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token);


--
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (userid);


--
-- Name: idx_shipping_method_translations_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_method_translations_name ON public.shippingmethodtranslations USING btree (name);


--
-- Name: idx_shipping_method_translations_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_method_translations_name_lower_textpattern ON public.shippingmethodtranslations USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_shipping_methods_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_methods_name ON public.shippingmethods USING btree (name);


--
-- Name: idx_shipping_methods_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_methods_name_lower_textpattern ON public.shippingmethods USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_shipping_zone_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_zone_name ON public.shippingzones USING btree (name);


--
-- Name: idx_shipping_zone_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shipping_zone_name_lower_textpattern ON public.shippingzones USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_shops_description; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shops_description ON public.shops USING btree (description);


--
-- Name: idx_shops_description_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shops_description_lower_textpattern ON public.shops USING btree (lower((description)::text) text_pattern_ops);


--
-- Name: idx_shops_name; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shops_name ON public.shops USING btree (name);


--
-- Name: idx_shops_name_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_shops_name_lower_textpattern ON public.shops USING btree (lower((name)::text) text_pattern_ops);


--
-- Name: idx_status_status; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_status_status ON public.status USING btree (status);


--
-- Name: idx_status_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_status_user_id ON public.status USING btree (userid);


--
-- Name: idx_uploadsessions_create_at; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_uploadsessions_create_at ON public.uploadsessions USING btree (createat);


--
-- Name: idx_uploadsessions_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_uploadsessions_user_id ON public.uploadsessions USING btree (type);


--
-- Name: idx_user_access_tokens_token; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_user_access_tokens_token ON public.useraccesstokens USING btree (token);


--
-- Name: idx_user_access_tokens_user_id; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_user_access_tokens_user_id ON public.useraccesstokens USING btree (userid);


--
-- Name: idx_users_all_no_full_name_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_all_no_full_name_txt ON public.users USING gin (to_tsvector('english'::regconfig, (((((username)::text || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));


--
-- Name: idx_users_all_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_all_txt ON public.users USING gin (to_tsvector('english'::regconfig, (((((((((username)::text || ' '::text) || (firstname)::text) || ' '::text) || (lastname)::text) || ' '::text) || (nickname)::text) || ' '::text) || (email)::text)));


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_email_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_email_lower_textpattern ON public.users USING btree (lower((email)::text) text_pattern_ops);


--
-- Name: idx_users_firstname_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_firstname_lower_textpattern ON public.users USING btree (lower((firstname)::text) text_pattern_ops);


--
-- Name: idx_users_lastname_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_lastname_lower_textpattern ON public.users USING btree (lower((lastname)::text) text_pattern_ops);


--
-- Name: idx_users_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_metadata ON public.users USING btree (metadata);


--
-- Name: idx_users_names_no_full_name_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_names_no_full_name_txt ON public.users USING gin (to_tsvector('english'::regconfig, (((username)::text || ' '::text) || (nickname)::text)));


--
-- Name: idx_users_names_txt; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_names_txt ON public.users USING gin (to_tsvector('english'::regconfig, (((((((username)::text || ' '::text) || (firstname)::text) || ' '::text) || (lastname)::text) || ' '::text) || (nickname)::text)));


--
-- Name: idx_users_nickname_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_nickname_lower_textpattern ON public.users USING btree (lower((nickname)::text) text_pattern_ops);


--
-- Name: idx_users_private_metadata; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_private_metadata ON public.users USING btree (privatemetadata);


--
-- Name: idx_users_username_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_users_username_lower_textpattern ON public.users USING btree (lower((username)::text) text_pattern_ops);


--
-- Name: idx_vouchers_code; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_vouchers_code ON public.vouchers USING btree (code);


--
-- Name: idx_warehouses_email; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_warehouses_email ON public.warehouses USING btree (email);


--
-- Name: idx_warehouses_email_lower_textpattern; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_warehouses_email_lower_textpattern ON public.warehouses USING btree (lower((email)::text) text_pattern_ops);


--
-- Name: idx_wishlist_items; Type: INDEX; Schema: public; Owner: minh
--

CREATE INDEX idx_wishlist_items ON public.wishlistitems USING btree (createat);


--
-- Name: allocations fk_allocations_orderlines; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT fk_allocations_orderlines FOREIGN KEY (stockid) REFERENCES public.orderlines(id) ON DELETE CASCADE;


--
-- Name: allocations fk_allocations_stocks; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.allocations
    ADD CONSTRAINT fk_allocations_stocks FOREIGN KEY (orderlineid) REFERENCES public.stocks(id) ON DELETE CASCADE;


--
-- Name: assignedpageattributes fk_assignedpageattributes_attributepages; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_attributepages FOREIGN KEY (assignmentid) REFERENCES public.attributepages(id) ON DELETE CASCADE;


--
-- Name: assignedpageattributes fk_assignedpageattributes_pages; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributes
    ADD CONSTRAINT fk_assignedpageattributes_pages FOREIGN KEY (pageid) REFERENCES public.pages(id) ON DELETE CASCADE;


--
-- Name: assignedpageattributevalues fk_assignedpageattributevalues_assignedpageattributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_assignedpageattributes FOREIGN KEY (assignmentid) REFERENCES public.assignedpageattributes(id) ON DELETE CASCADE;


--
-- Name: assignedpageattributevalues fk_assignedpageattributevalues_attributevalues; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES public.attributevalues(id) ON DELETE CASCADE;


--
-- Name: assignedproductattributes fk_assignedproductattributes_attributeproducts; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_attributeproducts FOREIGN KEY (assignmentid) REFERENCES public.attributeproducts(id) ON DELETE CASCADE;


--
-- Name: assignedproductattributes fk_assignedproductattributes_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: assignedproductattributevalues fk_assignedproductattributevalues_assignedproductattributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_assignedproductattributes FOREIGN KEY (assignmentid) REFERENCES public.assignedproductattributes(id) ON DELETE CASCADE;


--
-- Name: assignedproductattributevalues fk_assignedproductattributevalues_attributevalues; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES public.attributevalues(id) ON DELETE CASCADE;


--
-- Name: assignedvariantattributes fk_assignedvariantattributes_attributevariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_attributevariants FOREIGN KEY (assignmentid) REFERENCES public.attributevariants(id) ON DELETE CASCADE;


--
-- Name: assignedvariantattributes fk_assignedvariantattributes_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_productvariants FOREIGN KEY (variantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: assignedvariantattributevalues fk_assignedvariantattributevalues_assignedvariantattributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_assignedvariantattributes FOREIGN KEY (assignmentid) REFERENCES public.assignedvariantattributes(id) ON DELETE CASCADE;


--
-- Name: assignedvariantattributevalues fk_assignedvariantattributevalues_attributevalues; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES public.attributevalues(id) ON DELETE CASCADE;


--
-- Name: attributepages fk_attributepages_attributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributepages
    ADD CONSTRAINT fk_attributepages_attributes FOREIGN KEY (attributeid) REFERENCES public.attributes(id) ON DELETE CASCADE;


--
-- Name: attributepages fk_attributepages_pagetypes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributepages
    ADD CONSTRAINT fk_attributepages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES public.pagetypes(id) ON DELETE CASCADE;


--
-- Name: attributeproducts fk_attributeproducts_attributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributeproducts
    ADD CONSTRAINT fk_attributeproducts_attributes FOREIGN KEY (attributeid) REFERENCES public.attributes(id) ON DELETE CASCADE;


--
-- Name: attributeproducts fk_attributeproducts_producttypes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributeproducts
    ADD CONSTRAINT fk_attributeproducts_producttypes FOREIGN KEY (producttypeid) REFERENCES public.producttypes(id) ON DELETE CASCADE;


--
-- Name: attributevalues fk_attributevalues_attributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevalues
    ADD CONSTRAINT fk_attributevalues_attributes FOREIGN KEY (attributeid) REFERENCES public.attributes(id) ON DELETE CASCADE;


--
-- Name: attributevariants fk_attributevariants_attributes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevariants
    ADD CONSTRAINT fk_attributevariants_attributes FOREIGN KEY (attributeid) REFERENCES public.attributes(id) ON DELETE CASCADE;


--
-- Name: attributevariants fk_attributevariants_producttypes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.attributevariants
    ADD CONSTRAINT fk_attributevariants_producttypes FOREIGN KEY (producttypeid) REFERENCES public.producttypes(id) ON DELETE CASCADE;


--
-- Name: categories fk_categories_categories; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT fk_categories_categories FOREIGN KEY (parentid) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: checkoutlines fk_checkoutlines_checkouts; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkoutlines
    ADD CONSTRAINT fk_checkoutlines_checkouts FOREIGN KEY (checkoutid) REFERENCES public.checkouts(token) ON DELETE CASCADE;


--
-- Name: checkoutlines fk_checkoutlines_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkoutlines
    ADD CONSTRAINT fk_checkoutlines_productvariants FOREIGN KEY (variantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: checkouts fk_checkouts_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_addresses FOREIGN KEY (billingaddressid) REFERENCES public.addresses(id);


--
-- Name: checkouts fk_checkouts_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_channels FOREIGN KEY (channelid) REFERENCES public.channels(id);


--
-- Name: checkouts fk_checkouts_shippingmethods; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES public.shippingmethods(id);


--
-- Name: checkouts fk_checkouts_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_shops FOREIGN KEY (shopid) REFERENCES public.shops(id) ON DELETE CASCADE;


--
-- Name: checkouts fk_checkouts_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: checkouts fk_checkouts_warehouses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.checkouts
    ADD CONSTRAINT fk_checkouts_warehouses FOREIGN KEY (collectionpointid) REFERENCES public.warehouses(id);


--
-- Name: collectionchannellistings fk_collectionchannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: collectionchannellistings fk_collectionchannellistings_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- Name: collections fk_collections_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collections
    ADD CONSTRAINT fk_collections_shops FOREIGN KEY (shopid) REFERENCES public.shops(id) ON DELETE CASCADE;


--
-- Name: collectiontranslations fk_collectiontranslations_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.collectiontranslations
    ADD CONSTRAINT fk_collectiontranslations_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- Name: customerevents fk_customerevents_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.customerevents
    ADD CONSTRAINT fk_customerevents_orders FOREIGN KEY (orderid) REFERENCES public.orders(id);


--
-- Name: customerevents fk_customerevents_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.customerevents
    ADD CONSTRAINT fk_customerevents_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: customernotes fk_customernotes_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.customernotes
    ADD CONSTRAINT fk_customernotes_users FOREIGN KEY (userid) REFERENCES public.users(id);


--
-- Name: digitalcontents fk_digitalcontents_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontents
    ADD CONSTRAINT fk_digitalcontents_productvariants FOREIGN KEY (productvariantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: digitalcontents fk_digitalcontents_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontents
    ADD CONSTRAINT fk_digitalcontents_shops FOREIGN KEY (shopid) REFERENCES public.shops(id) ON DELETE CASCADE;


--
-- Name: digitalcontenturls fk_digitalcontenturls_digitalcontents; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_digitalcontents FOREIGN KEY (contentid) REFERENCES public.digitalcontents(id) ON DELETE CASCADE;


--
-- Name: digitalcontenturls fk_digitalcontenturls_orderlines; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.digitalcontenturls
    ADD CONSTRAINT fk_digitalcontenturls_orderlines FOREIGN KEY (lineid) REFERENCES public.orderlines(id) ON DELETE CASCADE;


--
-- Name: exportevents fk_exportevents_exportfiles; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.exportevents
    ADD CONSTRAINT fk_exportevents_exportfiles FOREIGN KEY (exportfileid) REFERENCES public.exportfiles(id);


--
-- Name: exportevents fk_exportevents_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.exportevents
    ADD CONSTRAINT fk_exportevents_users FOREIGN KEY (userid) REFERENCES public.users(id);


--
-- Name: exportfiles fk_exportfiles_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.exportfiles
    ADD CONSTRAINT fk_exportfiles_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: fulfillmentlines fk_fulfillmentlines_fulfillments; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_fulfillments FOREIGN KEY (fulfillmentid) REFERENCES public.fulfillments(id) ON DELETE CASCADE;


--
-- Name: fulfillmentlines fk_fulfillmentlines_orderlines; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_orderlines FOREIGN KEY (orderlineid) REFERENCES public.orderlines(id) ON DELETE CASCADE;


--
-- Name: fulfillmentlines fk_fulfillmentlines_stocks; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_stocks FOREIGN KEY (stockid) REFERENCES public.stocks(id);


--
-- Name: fulfillments fk_fulfillments_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.fulfillments
    ADD CONSTRAINT fk_fulfillments_orders FOREIGN KEY (orderid) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: giftcardcheckouts fk_giftcardcheckouts_checkouts; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_checkouts FOREIGN KEY (checkoutid) REFERENCES public.checkouts(token);


--
-- Name: giftcardcheckouts fk_giftcardcheckouts_giftcards; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardcheckouts
    ADD CONSTRAINT fk_giftcardcheckouts_giftcards FOREIGN KEY (giftcardid) REFERENCES public.giftcards(id);


--
-- Name: giftcardevents fk_giftcardevents_giftcards; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcardevents
    ADD CONSTRAINT fk_giftcardevents_giftcards FOREIGN KEY (giftcardid) REFERENCES public.giftcards(id) ON DELETE CASCADE;


--
-- Name: giftcards fk_giftcards_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcards
    ADD CONSTRAINT fk_giftcards_products FOREIGN KEY (productid) REFERENCES public.products(id);


--
-- Name: giftcards fk_giftcards_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.giftcards
    ADD CONSTRAINT fk_giftcards_users FOREIGN KEY (createdbyid) REFERENCES public.users(id);


--
-- Name: invoiceevents fk_invoiceevents_invoices; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoiceevents
    ADD CONSTRAINT fk_invoiceevents_invoices FOREIGN KEY (invoiceid) REFERENCES public.invoices(id);


--
-- Name: invoiceevents fk_invoiceevents_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoiceevents
    ADD CONSTRAINT fk_invoiceevents_orders FOREIGN KEY (orderid) REFERENCES public.orders(id);


--
-- Name: invoiceevents fk_invoiceevents_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoiceevents
    ADD CONSTRAINT fk_invoiceevents_users FOREIGN KEY (userid) REFERENCES public.users(id);


--
-- Name: invoices fk_invoices_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT fk_invoices_orders FOREIGN KEY (orderid) REFERENCES public.orders(id);


--
-- Name: menuitems fk_menuitems_categories; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT fk_menuitems_categories FOREIGN KEY (categoryid) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: menuitems fk_menuitems_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT fk_menuitems_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- Name: menuitems fk_menuitems_menuitems; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT fk_menuitems_menuitems FOREIGN KEY (parentid) REFERENCES public.menuitems(id) ON DELETE CASCADE;


--
-- Name: menuitems fk_menuitems_menus; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT fk_menuitems_menus FOREIGN KEY (menuid) REFERENCES public.menus(id) ON DELETE CASCADE;


--
-- Name: menuitems fk_menuitems_pages; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitems
    ADD CONSTRAINT fk_menuitems_pages FOREIGN KEY (pageid) REFERENCES public.pages(id) ON DELETE CASCADE;


--
-- Name: menuitemtranslations fk_menuitemtranslations_menuitems; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.menuitemtranslations
    ADD CONSTRAINT fk_menuitemtranslations_menuitems FOREIGN KEY (menuitemid) REFERENCES public.menuitems(id) ON DELETE CASCADE;


--
-- Name: orderdiscounts fk_orderdiscounts_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderdiscounts
    ADD CONSTRAINT fk_orderdiscounts_orders FOREIGN KEY (orderid) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orderevents fk_orderevents_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderevents
    ADD CONSTRAINT fk_orderevents_orders FOREIGN KEY (orderid) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orderevents fk_orderevents_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderevents
    ADD CONSTRAINT fk_orderevents_users FOREIGN KEY (userid) REFERENCES public.users(id);


--
-- Name: ordergiftcards fk_ordergiftcards_giftcards; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.ordergiftcards
    ADD CONSTRAINT fk_ordergiftcards_giftcards FOREIGN KEY (giftcardid) REFERENCES public.giftcards(id);


--
-- Name: ordergiftcards fk_ordergiftcards_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.ordergiftcards
    ADD CONSTRAINT fk_ordergiftcards_orders FOREIGN KEY (orderid) REFERENCES public.orders(id);


--
-- Name: orderlines fk_orderlines_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderlines
    ADD CONSTRAINT fk_orderlines_orders FOREIGN KEY (orderid) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orderlines fk_orderlines_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orderlines
    ADD CONSTRAINT fk_orderlines_productvariants FOREIGN KEY (variantid) REFERENCES public.productvariants(id);


--
-- Name: orders fk_orders_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_addresses FOREIGN KEY (billingaddressid) REFERENCES public.addresses(id);


--
-- Name: orders fk_orders_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_channels FOREIGN KEY (channelid) REFERENCES public.channels(id);


--
-- Name: orders fk_orders_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_orders FOREIGN KEY (originalid) REFERENCES public.orders(id);


--
-- Name: orders fk_orders_shippingmethods; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES public.shippingmethods(id);


--
-- Name: orders fk_orders_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_shops FOREIGN KEY (shopid) REFERENCES public.shops(id);


--
-- Name: orders fk_orders_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_users FOREIGN KEY (userid) REFERENCES public.users(id);


--
-- Name: orders fk_orders_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id);


--
-- Name: orders fk_orders_warehouses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_warehouses FOREIGN KEY (collectionpointid) REFERENCES public.warehouses(id);


--
-- Name: pages fk_pages_pagetypes; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pages
    ADD CONSTRAINT fk_pages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES public.pagetypes(id) ON DELETE CASCADE;


--
-- Name: pagetranslations fk_pagetranslations_pages; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.pagetranslations
    ADD CONSTRAINT fk_pagetranslations_pages FOREIGN KEY (pageid) REFERENCES public.pages(id) ON DELETE CASCADE;


--
-- Name: payments fk_payments_checkouts; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT fk_payments_checkouts FOREIGN KEY (checkoutid) REFERENCES public.checkouts(token);


--
-- Name: payments fk_payments_orders; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT fk_payments_orders FOREIGN KEY (orderid) REFERENCES public.orders(id);


--
-- Name: productchannellistings fk_productchannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productchannellistings
    ADD CONSTRAINT fk_productchannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: productchannellistings fk_productchannellistings_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productchannellistings
    ADD CONSTRAINT fk_productchannellistings_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: productcollections fk_productcollections_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productcollections
    ADD CONSTRAINT fk_productcollections_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- Name: productcollections fk_productcollections_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productcollections
    ADD CONSTRAINT fk_productcollections_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: productmedias fk_productmedias_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productmedias
    ADD CONSTRAINT fk_productmedias_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: producttranslations fk_producttranslations_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.producttranslations
    ADD CONSTRAINT fk_producttranslations_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: productvariantchannellistings fk_productvariantchannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariantchannellistings
    ADD CONSTRAINT fk_productvariantchannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: productvariantchannellistings fk_productvariantchannellistings_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariantchannellistings
    ADD CONSTRAINT fk_productvariantchannellistings_productvariants FOREIGN KEY (variantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: productvariants fk_productvariants_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvariants
    ADD CONSTRAINT fk_productvariants_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: productvarianttranslations fk_productvarianttranslations_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.productvarianttranslations
    ADD CONSTRAINT fk_productvarianttranslations_productvariants FOREIGN KEY (productvariantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: salecategories fk_salecategories_categories; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecategories
    ADD CONSTRAINT fk_salecategories_categories FOREIGN KEY (categoryid) REFERENCES public.categories(id);


--
-- Name: salecategories fk_salecategories_sales; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecategories
    ADD CONSTRAINT fk_salecategories_sales FOREIGN KEY (saleid) REFERENCES public.sales(id);


--
-- Name: salechannellistings fk_salechannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salechannellistings
    ADD CONSTRAINT fk_salechannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: salechannellistings fk_salechannellistings_sales; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salechannellistings
    ADD CONSTRAINT fk_salechannellistings_sales FOREIGN KEY (saleid) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- Name: salecollections fk_salecollections_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecollections
    ADD CONSTRAINT fk_salecollections_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id);


--
-- Name: salecollections fk_salecollections_sales; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.salecollections
    ADD CONSTRAINT fk_salecollections_sales FOREIGN KEY (saleid) REFERENCES public.sales(id);


--
-- Name: saleproducts fk_saleproducts_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproducts
    ADD CONSTRAINT fk_saleproducts_products FOREIGN KEY (productid) REFERENCES public.products(id);


--
-- Name: saleproducts fk_saleproducts_sales; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saleproducts
    ADD CONSTRAINT fk_saleproducts_sales FOREIGN KEY (saleid) REFERENCES public.sales(id);


--
-- Name: sales fk_sales_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT fk_sales_shops FOREIGN KEY (shopid) REFERENCES public.shops(id);


--
-- Name: saletranslations fk_saletranslations_sales; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.saletranslations
    ADD CONSTRAINT fk_saletranslations_sales FOREIGN KEY (saleid) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- Name: shippingmethodchannellistings fk_shippingmethodchannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: shippingmethodchannellistings fk_shippingmethodchannellistings_shippingmethods; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES public.shippingmethods(id) ON DELETE CASCADE;


--
-- Name: shippingmethodexcludedproducts fk_shippingmethodexcludedproducts_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodexcludedproducts
    ADD CONSTRAINT fk_shippingmethodexcludedproducts_products FOREIGN KEY (productid) REFERENCES public.products(id);


--
-- Name: shippingmethodexcludedproducts fk_shippingmethodexcludedproducts_shippingmethods; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodexcludedproducts
    ADD CONSTRAINT fk_shippingmethodexcludedproducts_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES public.shippingmethods(id);


--
-- Name: shippingmethodpostalcoderules fk_shippingmethodpostalcoderules_shippingmethods; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethodpostalcoderules
    ADD CONSTRAINT fk_shippingmethodpostalcoderules_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES public.shippingmethods(id) ON DELETE CASCADE;


--
-- Name: shippingmethods fk_shippingmethods_shippingzones; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingmethods
    ADD CONSTRAINT fk_shippingmethods_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES public.shippingzones(id) ON DELETE CASCADE;


--
-- Name: shippingzonechannels fk_shippingzonechannels_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_channels FOREIGN KEY (channelid) REFERENCES public.channels(id);


--
-- Name: shippingzonechannels fk_shippingzonechannels_shippingzones; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES public.shippingzones(id);


--
-- Name: shops fk_shops_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT fk_shops_addresses FOREIGN KEY (addressid) REFERENCES public.addresses(id);


--
-- Name: shops fk_shops_menus; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT fk_shops_menus FOREIGN KEY (topmenuid) REFERENCES public.menus(id);


--
-- Name: shops fk_shops_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT fk_shops_users FOREIGN KEY (ownerid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: shopstaffs fk_shopstaffs_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shopstaffs
    ADD CONSTRAINT fk_shopstaffs_shops FOREIGN KEY (shopid) REFERENCES public.shops(id);


--
-- Name: shopstaffs fk_shopstaffs_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shopstaffs
    ADD CONSTRAINT fk_shopstaffs_users FOREIGN KEY (staffid) REFERENCES public.users(id);


--
-- Name: shoptranslations fk_shoptranslations_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.shoptranslations
    ADD CONSTRAINT fk_shoptranslations_shops FOREIGN KEY (shopid) REFERENCES public.shops(id) ON DELETE CASCADE;


--
-- Name: staffnotificationrecipients fk_staffnotificationrecipients_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.staffnotificationrecipients
    ADD CONSTRAINT fk_staffnotificationrecipients_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: stocks fk_stocks_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.stocks
    ADD CONSTRAINT fk_stocks_productvariants FOREIGN KEY (productvariantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: stocks fk_stocks_warehouses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.stocks
    ADD CONSTRAINT fk_stocks_warehouses FOREIGN KEY (warehouseid) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: transactions fk_transactions_payments; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT fk_transactions_payments FOREIGN KEY (paymentid) REFERENCES public.payments(id);


--
-- Name: useraddresses fk_useraddresses_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraddresses
    ADD CONSTRAINT fk_useraddresses_addresses FOREIGN KEY (addressid) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: useraddresses fk_useraddresses_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.useraddresses
    ADD CONSTRAINT fk_useraddresses_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users fk_users_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_addresses FOREIGN KEY (defaultshippingaddressid) REFERENCES public.addresses(id);


--
-- Name: variantmedias fk_variantmedias_productmedias; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.variantmedias
    ADD CONSTRAINT fk_variantmedias_productmedias FOREIGN KEY (mediaid) REFERENCES public.productmedias(id) ON DELETE CASCADE;


--
-- Name: variantmedias fk_variantmedias_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.variantmedias
    ADD CONSTRAINT fk_variantmedias_productvariants FOREIGN KEY (variantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: vouchercategories fk_vouchercategories_categories; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercategories
    ADD CONSTRAINT fk_vouchercategories_categories FOREIGN KEY (categoryid) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: vouchercategories fk_vouchercategories_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercategories
    ADD CONSTRAINT fk_vouchercategories_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: voucherchannellistings fk_voucherchannellistings_channels; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_channels FOREIGN KEY (channelid) REFERENCES public.channels(id) ON DELETE CASCADE;


--
-- Name: voucherchannellistings fk_voucherchannellistings_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherchannellistings
    ADD CONSTRAINT fk_voucherchannellistings_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: vouchercollections fk_vouchercollections_collections; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercollections
    ADD CONSTRAINT fk_vouchercollections_collections FOREIGN KEY (collectionid) REFERENCES public.collections(id) ON DELETE CASCADE;


--
-- Name: vouchercollections fk_vouchercollections_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercollections
    ADD CONSTRAINT fk_vouchercollections_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: vouchercustomers fk_vouchercustomers_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchercustomers
    ADD CONSTRAINT fk_vouchercustomers_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: voucherproducts fk_voucherproducts_products; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproducts
    ADD CONSTRAINT fk_voucherproducts_products FOREIGN KEY (productid) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: voucherproducts fk_voucherproducts_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.voucherproducts
    ADD CONSTRAINT fk_voucherproducts_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: vouchers fk_vouchers_shops; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchers
    ADD CONSTRAINT fk_vouchers_shops FOREIGN KEY (shopid) REFERENCES public.shops(id) ON DELETE CASCADE;


--
-- Name: vouchertranslations fk_vouchertranslations_vouchers; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.vouchertranslations
    ADD CONSTRAINT fk_vouchertranslations_vouchers FOREIGN KEY (voucherid) REFERENCES public.vouchers(id) ON DELETE CASCADE;


--
-- Name: warehouses fk_warehouses_addresses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT fk_warehouses_addresses FOREIGN KEY (addressid) REFERENCES public.addresses(id);


--
-- Name: warehouseshippingzones fk_warehouseshippingzones_shippingzones; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES public.shippingzones(id);


--
-- Name: warehouseshippingzones fk_warehouseshippingzones_warehouses; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.warehouseshippingzones
    ADD CONSTRAINT fk_warehouseshippingzones_warehouses FOREIGN KEY (warehouseid) REFERENCES public.warehouses(id);


--
-- Name: wishlistitemproductvariants fk_wishlistitemproductvariants_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_productvariants FOREIGN KEY (productvariantid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: wishlistitemproductvariants fk_wishlistitemproductvariants_wishlistitems; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitemproductvariants
    ADD CONSTRAINT fk_wishlistitemproductvariants_wishlistitems FOREIGN KEY (wishlistitemid) REFERENCES public.wishlistitems(id) ON DELETE CASCADE;


--
-- Name: wishlistitems fk_wishlistitems_productvariants; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitems
    ADD CONSTRAINT fk_wishlistitems_productvariants FOREIGN KEY (productid) REFERENCES public.productvariants(id) ON DELETE CASCADE;


--
-- Name: wishlistitems fk_wishlistitems_wishlists; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlistitems
    ADD CONSTRAINT fk_wishlistitems_wishlists FOREIGN KEY (wishlistid) REFERENCES public.wishlists(id) ON DELETE CASCADE;


--
-- Name: wishlists fk_wishlists_users; Type: FK CONSTRAINT; Schema: public; Owner: minh
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT fk_wishlists_users FOREIGN KEY (userid) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

