CREATE TABLE appartement(
    appartement_id INT AUTO_INCREMENT,
    area INT,
    capacity INT,
    streetNumber INT,
    streetName VARCHAR(255),
    city VARCHAR(255),
    disponibility BOOL NOT NULL DEFAULT 1,
    PRIMARY KEY(appartement_id)
);

CREATE TABLE appUser(
    user_id INT AUTO_INCREMENT,
    admin BOOL NOT NULL DEFAULT 0,
    api_key VARCHAR(255) NOT NULL,
    PRIMARY KEY(user_id)
);

CREATE TABLE own(
    appartement_id INT,
    user_id INT,
    PRIMARY KEY(appartement_id, user_id),
    FOREIGN KEY(appartement_id) REFERENCES appartement(appartement_id),
    FOREIGN KEY(user_id) REFERENCES appUser(user_id)
);

CREATE TABLE rent(
     appartement_id INT,
     user_id INT,
     date_begin DATE NOT NULL,
     date_end DATE NOT NULL,
     price DECIMAL(12, 2) NOT NULL,
     PRIMARY KEY(appartement_id, user_id),
     FOREIGN KEY(appartement_id) REFERENCES appartement(appartement_id),
     FOREIGN KEY(user_id) REFERENCES appUser(user_id)
);

INSERT INTO appUser(user_id, admin, api_key) VALUES (0, 1, "admin");
