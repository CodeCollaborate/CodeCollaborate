CREATE DATABASE  IF NOT EXISTS `cc` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
# CREATE DATABASE  IF NOT EXISTS `testing` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;

GRANT EXECUTE ON cc.* to username@'%' IDENTIFIED BY 'pass';
# GRANT ALL PRIVILEGES ON testing.* to username@localhost IDENTIFIED BY 'pass';
