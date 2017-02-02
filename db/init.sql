# User model service

CREATE DATABASE service_model_user CHARACTER SET utf8 COLLATE utf8_unicode_ci;
CREATE USER "service_model_user"@"%" IDENTIFIED BY "service_model_user";
GRANT SELECT, INSERT, UPDATE ON service_model_user.* TO "service_model_user"@"%";
