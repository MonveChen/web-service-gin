##### 建表

```sql
CREATE TABLE [IF NOT EXISTS] public.contact_token (
  	id SERIAL PRIMARY KEY,
   	chainId VARCHAR(20) NOT NULL,
		token VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NULL,
    decimals SMALLINT NULL
);

CREATE UNIQUE INDEX idx_main ON public.contact_token (chainId, token);



```

