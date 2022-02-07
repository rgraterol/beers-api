CREATE TABLE beers (
   id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
   name VARCHAR(50) not null,
   brewery VARCHAR(95),
   country VARCHAR(50),
   price FLOAT,
   currency VARCHAR(3),
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP
);

CREATE INDEX beer_id_index ON beers (id);