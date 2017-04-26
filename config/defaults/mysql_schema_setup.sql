CREATE DATABASE  IF NOT EXISTS `cc` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `cc`;
-- MySQL dump 10.13  Distrib 5.7.17, for Linux (x86_64)
--
-- Host: localhost    Database: cc
-- ------------------------------------------------------
-- Server version	5.7.17

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `File`
--

DROP TABLE IF EXISTS `File`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `File` (
  `FileID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Creator` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  `CreationDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `RelativePath` varchar(2083) COLLATE utf8_unicode_ci NOT NULL,
  `ProjectID` bigint(20) NOT NULL,
  `Filename` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`FileID`),
  UNIQUE KEY `FileID_UNIQUE` (`FileID`),
  KEY `fk_File_Username_idx` (`Creator`),
  KEY `fk_File_ProjectID_idx` (`ProjectID`),
  CONSTRAINT `fk_File_ProjectID` FOREIGN KEY (`ProjectID`) REFERENCES `Project` (`ProjectID`) ON DELETE NO ACTION ON UPDATE CASCADE,
  CONSTRAINT `fk_File_Username` FOREIGN KEY (`Creator`) REFERENCES `User` (`Username`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `File`
--

LOCK TABLES `File` WRITE;
/*!40000 ALTER TABLE `File` DISABLE KEYS */;
/*!40000 ALTER TABLE `File` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Permissions`
--

DROP TABLE IF EXISTS `Permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Permissions` (
  `Username` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  `ProjectID` bigint(20) NOT NULL,
  `PermissionLevel` tinyint(1) NOT NULL DEFAULT '0',
  `GrantedBy` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  `GrantedDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ProjectID`,`Username`),
  KEY `fk_ProjectID_idx` (`ProjectID`),
  KEY `fk_Permissions_Username_idx` (`Username`),
  CONSTRAINT `fk_Permissions_ProjectID` FOREIGN KEY (`ProjectID`) REFERENCES `Project` (`ProjectID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_Permissions_Username` FOREIGN KEY (`Username`) REFERENCES `User` (`Username`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Permissions`
--

LOCK TABLES `Permissions` WRITE;
/*!40000 ALTER TABLE `Permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `Permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Project`
--

DROP TABLE IF EXISTS `Project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Project` (
  `ProjectID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `Owner` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`ProjectID`),
  UNIQUE KEY `ProjectID_UNIQUE` (`ProjectID`),
  UNIQUE KEY `NameOwner_UNIQUE` (`Name`,`Owner`),
  KEY `fk_Project_Username` (`Owner`),
  CONSTRAINT `fk_Project_Username` FOREIGN KEY (`Owner`) REFERENCES `User` (`Username`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Project`
--

LOCK TABLES `Project` WRITE;
/*!40000 ALTER TABLE `Project` DISABLE KEYS */;
/*!40000 ALTER TABLE `Project` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `cc`.`Project_BEFORE_DELETE` BEFORE DELETE ON `Project` FOR EACH ROW
  BEGIN
    DELETE FROM Permissions
    WHERE Permissions.ProjectID = OLD.ProjectID;
    DELETE FROM `File`
    WHERE `File`.ProjectID = OLD.ProjectID;
  END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `User` (
  `Username` varchar(25) COLLATE utf8_unicode_ci NOT NULL,
  `Password` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `Email` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `FirstName` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `LastName` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`Username`),
  UNIQUE KEY `Email_UNIQUE` (`Email`),
  KEY `Email_INDEX` (`Email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
INSERT INTO `User` VALUES ('test','test','test','test','test');
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'cc'
--
/*!50003 DROP PROCEDURE IF EXISTS `file_create` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `file_create`(IN username varchar(25), IN filename varchar(50), IN relativePath varchar(2083), IN projectID bigint(20))
  BEGIN
    IF (SELECT COUNT(*) FROM File WHERE File.ProjectID = projectID AND File.RelativePath = relativePath AND File.Filename = filename)>0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "Project already contains file at the given location" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        INSERT INTO `File`
        (Creator, RelativePath, ProjectID, Filename)
        VALUES (username, relativePath, projectID, filename);
        SELECT LAST_INSERT_ID();
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `file_delete` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `file_delete`(IN fileID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM File WHERE File.FileID = fileID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such fileID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        DELETE FROM `File`
        WHERE `File`.FileID = fileID;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `file_get_info` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `file_get_info`(IN fileID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM File WHERE File.FileID = fileID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such fileID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT `File`.`Creator`, `File`.`CreationDate`, `File`.`RelativePath`, `File`.`ProjectID`, `File`.`Filename`
        FROM File
        WHERE `File`.`FileID` = fileID;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `file_move` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `file_move`(IN fileID bigint(20), IN newPath varchar(2083), IN newName varchar(50))
  BEGIN
    IF (SELECT count(*) FROM File WHERE File.FileID = fileID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such fileID found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM File WHERE File.ProjectID = (SELECT ProjectID FROM File WHERE File.FileID = fileID) && File.RelativePath = newPath && File.Filename = (SELECT Filename FROM File WHERE File.FileID = fileID))>0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "Project already contains file at the given location" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        UPDATE `File`
        SET `File`.RelativePath = newPath, `File`.Filename = newName
        WHERE `File`.FileID = fileID;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_create` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_create`(IN projectName varchar(50), IN username varchar(25))
  BEGIN
    IF (SELECT COUNT(*) FROM Project WHERE Project.Name = projectName AND Project.Owner = username)<=0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "Owner already has a project with the given name" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        INSERT INTO Project (`Name`, `Owner`)
        VALUES (projectName, username);
        SELECT LAST_INSERT_ID();
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_delete` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_delete`(IN projectID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        DELETE FROM Project
        WHERE Project.ProjectID = projectID;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_get_files` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_get_files`(IN projectID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT *
        FROM File
        WHERE File.ProjectID = projectID;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_grant_permissions` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_grant_permissions`(IN projectID bigint(20),
                                                                        IN grantUsername varchar(25),
                                                                        IN permissionLevel tinyint(1),
                                                                        IN grantedByUsername varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = grantUsername)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM User WHERE User.Username = grantedByUsername)<=0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "No such grantedByUsername found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 3 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        INSERT INTO `Permissions`
        (Username, ProjectID, PermissionLevel, GrantedBy)
        VALUES (grantUsername, projectID, permissionLevel, grantedByUsername)
        ON DUPLICATE KEY UPDATE
          PermissionLevel = permissionLevel,
          GrantedBy = grantedByUsername;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_lookup` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_lookup`(IN projectID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT `Project`.`Name`, `Permissions`.`Username`, `Permissions`.`PermissionLevel`, `Permissions`.`GrantedBy`, `Permissions`.`GrantedDate`
        FROM Project LEFT JOIN Permissions
            ON Project.ProjectID = Permissions.ProjectID
        WHERE Project.ProjectID = projectID
        UNION
        SELECT `Project`.`Name`, `Project`.`Owner`, 10, `Project`.`Owner`, 0
        FROM `Project`
        WHERE `Project`.`ProjectID` = projectID;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_rename` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_rename`(IN username VARCHAR(25), IN projectID bigint(20), IN newName varchar(50))
  BEGIN
    IF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM Project WHERE Project.ProjectID != projectID && Project.Owner = (SELECT Owner FROM Project WHERE Project.ProjectID = projectID) && Project.Name = newName)>0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "Owner already has a project with the given name" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        UPDATE Project
        SET Project.Name = newName
        WHERE Project.ProjectID = projectID;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `project_revoke_permissions` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `project_revoke_permissions`(IN projectID bigint(20),
                                                                         IN revokeUsername varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = revokeUsername)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        DELETE FROM Permissions
        WHERE Permissions.ProjectID = projectID AND Permissions.Username = revokeUsername;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_delete` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_delete`(IN username varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        DELETE FROM `User`
        WHERE `User`.Username = username;
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_get_password` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_get_password`(IN username varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT Password
        FROM User where User.Username = username;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_get_projectids` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_get_projectids`(IN username varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT `Project`.`ProjectID` FROM `Project`
        WHERE `Project`.`Owner` = username;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_lookup` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_lookup`(IN username varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT FirstName, LastName, Email, Username
        FROM User where User.Username = username;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_projects` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_projects`(IN username varchar(25))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT `Project`.`ProjectID`, `Project`.`Name`, `Permissions`.`PermissionLevel`
        FROM (Permissions LEFT JOIN Project ON Permissions.ProjectID = Project.ProjectID)
        WHERE Permissions.Username = username
        UNION
        SELECT `Project`.`ProjectID`, `Project`.`Name`, 10
        FROM `Project`
        WHERE `Project`.`Owner` = username;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_project_permission` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_project_permission`(username varchar(25), projectID bigint(20))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)<=0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "No such username found" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM Project WHERE Project.ProjectID = projectID)<=0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "No such projectID found" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        SELECT Permissions.PermissionLevel
        FROM Permissions
        WHERE Permissions.Username = username and Permissions.ProjectID = projectID
        UNION
        SELECT 10
        FROM Project
        WHERE Project.ProjectID = projectID and Project.Owner = username;
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `user_register` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `user_register`(IN username varchar(25),
                                                            IN pass varchar(100),
                                                            IN email varchar(50),
                                                            IN firstName varchar(30),
                                                            IN lastName varchar(30))
  BEGIN
    IF (SELECT count(*) FROM User WHERE User.Username = username)>0 THEN
      BEGIN
        SELECT 1 as 'ERROR_CODE', "Username already taken" AS 'ERROR_MSG';
      END;
    ELSEIF (SELECT count(*) FROM User WHERE User.Email = email)>0 THEN
      BEGIN
        SELECT 2 as 'ERROR_CODE', "Email already registered" AS 'ERROR_MSG';
      END;
    ELSE
      BEGIN
        INSERT INTO User VALUES (username, pass, email, firstName, lastName);
        SELECT 0 AS 'ERROR_CODE', '' AS 'ERROR_MSG';
      END;
    END IF;
  END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-04-26 13:35:55