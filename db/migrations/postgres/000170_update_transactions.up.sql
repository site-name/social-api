ALTER TABLE ONLY transactions
    ADD CONSTRAINT fk_transactions_payments FOREIGN KEY (payment_id) REFERENCES payments(id);
