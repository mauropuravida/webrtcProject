

--INSERTS FOR TESTS

--users
INSERT INTO USERS (name, surname, age, email, created, password) VALUES ("Magali", "Boulanger", 22, "maga.boulanger8gmail.com", CURDATE(), "1234");


--cameras
INSERT INTO cameras(users_id, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (1, true, curdate(), "urquiza 315", "", "", 1);
INSERT INTO cameras(users_id, active, created, loc, token_session_camera, token_session_consumer, id_camera) VALUES (1, true, curdate(), "Pinto 1765", "", "", 1)
