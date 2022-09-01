ALTER TABLE ONLY shippingmethodexcludedproducts
    ADD CONSTRAINT fk_shippingmethodexcludedproducts_products FOREIGN KEY (productid) REFERENCES products(id);
ALTER TABLE ONLY shippingmethodexcludedproducts
    ADD CONSTRAINT fk_shippingmethodexcludedproducts_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id);
