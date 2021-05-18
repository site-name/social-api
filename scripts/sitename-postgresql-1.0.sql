--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;


BEGIN;
--
-- Create model User
--
CREATE TABLE "userprofile_user" ("id" serial NOT NULL PRIMARY KEY, "is_superuser" boolean NOT NULL, "email" varchar(254) NOT NULL UNIQUE, "is_staff" boolean NOT NULL, "is_active" boolean NOT NULL, "password" varchar(128) NOT NULL, "date_joined" timestamp with time zone NOT NULL, "last_login" timestamp with time zone NOT NULL);
--
-- Create model Address
--
CREATE TABLE "userprofile_address" ("id" serial NOT NULL PRIMARY KEY, "first_name" varchar(256) NOT NULL, "last_name" varchar(256) NOT NULL, "company_name" varchar(256) NOT NULL, "street_address_1" varchar(256) NOT NULL, "street_address_2" varchar(256) NOT NULL, "city" varchar(256) NOT NULL, "postal_code" varchar(20) NOT NULL, "country" varchar(2) NOT NULL, "country_area" varchar(128) NOT NULL, "phone" varchar(30) NOT NULL);
--
-- Add field addresses to user
--
CREATE TABLE "userprofile_user_addresses" ("id" serial NOT NULL PRIMARY KEY, "user_id" integer NOT NULL, "address_id" integer NOT NULL);
--
-- Add field default_billing_address to user
--
ALTER TABLE "userprofile_user" ADD COLUMN "default_billing_address_id" integer NULL CONSTRAINT "userprofile_user_default_billing_addr_0489abf1_fk_userprofi" REFERENCES "userprofile_address"("id") DEFERRABLE INITIALLY DEFERRED; SET CONSTRAINTS "userprofile_user_default_billing_addr_0489abf1_fk_userprofi" IMMEDIATE;
--
-- Add field default_shipping_address to user
--
ALTER TABLE "userprofile_user" ADD COLUMN "default_shipping_address_id" integer NULL CONSTRAINT "userprofile_user_default_shipping_add_aae7a203_fk_userprofi" REFERENCES "userprofile_address"("id") DEFERRABLE INITIALLY DEFERRED; SET CONSTRAINTS "userprofile_user_default_shipping_add_aae7a203_fk_userprofi" IMMEDIATE;
--
-- Add field groups to user
--
CREATE TABLE "userprofile_user_groups" ("id" serial NOT NULL PRIMARY KEY, "user_id" integer NOT NULL, "group_id" integer NOT NULL);
--
-- Add field user_permissions to user
--
CREATE TABLE "userprofile_user_user_permissions" ("id" serial NOT NULL PRIMARY KEY, "user_id" integer NOT NULL, "permission_id" integer NOT NULL);
CREATE INDEX "userprofile_user_email_b0fb0137_like" ON "userprofile_user" ("email" varchar_pattern_ops);
ALTER TABLE "userprofile_user_addresses" ADD CONSTRAINT "userprofile_user_addresses_user_id_address_id_6cb87bcc_uniq" UNIQUE ("user_id", "address_id");
ALTER TABLE "userprofile_user_addresses" ADD CONSTRAINT "userprofile_user_add_user_id_bb5aa55e_fk_userprofi" FOREIGN KEY ("user_id") REFERENCES "userprofile_user" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "userprofile_user_addresses" ADD CONSTRAINT "userprofile_user_add_address_id_ad7646b4_fk_userprofi" FOREIGN KEY ("address_id") REFERENCES "userprofile_address" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "userprofile_user_addresses_user_id_bb5aa55e" ON "userprofile_user_addresses" ("user_id");
CREATE INDEX "userprofile_user_addresses_address_id_ad7646b4" ON "userprofile_user_addresses" ("address_id");
CREATE INDEX "userprofile_user_default_billing_address_id_0489abf1" ON "userprofile_user" ("default_billing_address_id");
CREATE INDEX "userprofile_user_default_shipping_address_id_aae7a203" ON "userprofile_user" ("default_shipping_address_id");
ALTER TABLE "userprofile_user_groups" ADD CONSTRAINT "userprofile_user_groups_user_id_group_id_90ce1781_uniq" UNIQUE ("user_id", "group_id");
ALTER TABLE "userprofile_user_groups" ADD CONSTRAINT "userprofile_user_groups_user_id_5e712a24_fk_userprofile_user_id" FOREIGN KEY ("user_id") REFERENCES "userprofile_user" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "userprofile_user_groups" ADD CONSTRAINT "userprofile_user_groups_group_id_c7eec74e_fk_auth_group_id" FOREIGN KEY ("group_id") REFERENCES "auth_group" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "userprofile_user_groups_user_id_5e712a24" ON "userprofile_user_groups" ("user_id");
CREATE INDEX "userprofile_user_groups_group_id_c7eec74e" ON "userprofile_user_groups" ("group_id");
ALTER TABLE "userprofile_user_user_permissions" ADD CONSTRAINT "userprofile_user_user_pe_user_id_permission_id_706d65c8_uniq" UNIQUE ("user_id", "permission_id");
ALTER TABLE "userprofile_user_user_permissions" ADD CONSTRAINT "userprofile_user_use_user_id_6d654469_fk_userprofi" FOREIGN KEY ("user_id") REFERENCES "userprofile_user" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "userprofile_user_user_permissions" ADD CONSTRAINT "userprofile_user_use_permission_id_1caa8a71_fk_auth_perm" FOREIGN KEY ("permission_id") REFERENCES "auth_permission" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "userprofile_user_user_permissions_user_id_6d654469" ON "userprofile_user_user_permissions" ("user_id");
CREATE INDEX "userprofile_user_user_permissions_permission_id_1caa8a71" ON "userprofile_user_user_permissions" ("permission_id");
COMMIT;

