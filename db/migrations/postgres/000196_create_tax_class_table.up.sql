CREATE TABLE IF NOT EXISTS tax_classes (
    id varchar(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    metatata JSONB,
    private_metadata JSONB
);

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname ILIKE 'tax_calculation_strategy')
THEN
CREATE TYPE tax_calculation_strategy AS ENUM (
	'flat_rates',
	'tax_app'
);
END IF;
END $$;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname ILIKE 'taxable_object_discount_type')
THEN
CREATE TYPE taxable_object_discount_type AS ENUM (
	'subtotal',
	'shipping'
);
END IF;
END $$;

CREATE TABLE IF NOT EXISTS tax_class_country_rates (
    id varchar(36) NOT NULL PRIMARY KEY,
    tax_class_id varchar(36), -- can NULL
    country country_code NOT NULL,
    rate DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (tax_class_id) REFERENCES tax_classes(id)
);

CREATE INDEX IF NOT EXISTS idx_tax_class_country_rate ON tax_class_country_rates (country, tax_class_id);
ALTER TABLE tax_class_country_rates ADD CONSTRAINT fk_tax_class_id FOREIGN KEY (tax_class_id) REFERENCES tax_classes(id) ON DELETE CASCADE;
-- ALTER TABLE tax_class_country_rates ADD CONSTRAINT unique_country_without_tax_class UNIQUE (country) WHERE tax_class_id IS NULL;
ALTER TABLE tax_class_country_rates ADD CONSTRAINT unique_country_tax_class UNIQUE (country, tax_class_id);

CREATE TABLE IF NOT EXISTS tax_configurations (
    id varchar(36) NOT NULL PRIMARY KEY,
    channel_id varchar(36) NOT NULL,
    charge_taxes boolean not null default true,
    tax_calculation_strategy tax_calculation_strategy,
    display_gross_price boolean not null default true,
    prices_entered_with_tax boolean not null default true,
    tax_app_id varchar(256),
    metatata JSONB,
    private_metadata JSONB
);

ALTER TABLE tax_configurations ADD CONSTRAINT fk_channel_id FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS tax_configutation_per_countries (
    id varchar(36) NOT NULL PRIMARY KEY,
    tax_configuration_id varchar(36) NOT NULL,
    country country_code NOT NULL,
    charge_taxes boolean not null default true,
    tax_calculation_strategy tax_calculation_strategy,
    display_gross_price boolean not null default true,
    tax_app_id varchar(256)
);

ALTER TABLE tax_configutation_per_countries ADD CONSTRAINT fk_tax_configuration_id FOREIGN KEY (tax_configuration_id) REFERENCES tax_configurations(id) ON DELETE CASCADE;
ALTER TABLE tax_configutation_per_countries ADD CONSTRAINT unique_country_tax_configuration_id UNIQUE (country, tax_configuration_id);
