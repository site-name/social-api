ALTER TABLE ONLY transactions
    ADD CONSTRAINT fk_transactions_payments FOREIGN KEY (paymentid) REFERENCES payments(id);
