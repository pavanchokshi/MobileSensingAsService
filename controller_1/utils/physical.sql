-- MySQL dump 10.13  Distrib 5.6.27, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: mobisense_db
-- ------------------------------------------------------
-- Server version	5.6.27-0ubuntu0.14.04.1

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
-- Table structure for table `physical_sensor_metadata`
--

DROP TABLE IF EXISTS `physical_sensor_metadata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `physical_sensor_metadata` (
  `sensor_id` varchar(50) NOT NULL,
  `owner_id` varchar(50) DEFAULT NULL,
  `sensor_name` varchar(100) NOT NULL,
  `sensor_desc` varchar(200) DEFAULT NULL,
  `sensor_type` varchar(100) NOT NULL,
  `latitude` varchar(20) DEFAULT NULL,
  `longitude` varchar(20) DEFAULT NULL,
  `is_shared` tinyint(1) DEFAULT NULL,
  `is_activated` tinyint(1) DEFAULT NULL,
  `owner_enabled` tinyint(1) DEFAULT NULL,
  `created` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`sensor_id`),
  UNIQUE KEY `sensor_name` (`sensor_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `physical_sensor_metadata`
--

LOCK TABLES `physical_sensor_metadata` WRITE;
/*!40000 ALTER TABLE `physical_sensor_metadata` DISABLE KEYS */;
INSERT INTO `physical_sensor_metadata` VALUES ('02cb9ac7-9ac1-4e82-af2b-1b255bfde7cd\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','accelerometer123','This is a accelero sensor1','Axcel','50','150',1,0,0,'2015-11-24 19:14:15'),('074e7976-40d7-418e-a54c-0836d0a8727c\n','','','','',NULL,NULL,0,0,0,'2015-12-01 15:01:58'),('0cb3c939-0b48-43e5-8036-de345e65c861\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','Sense1','this is sensor of pavan','Proximity',NULL,NULL,0,0,0,'2015-12-01 15:32:21'),('253e9e60-fb25-49f0-8cb2-72d8d8b6b15b\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','Dhruv','This is a dhruv\'s personal sensor','Secret','23.02','-89.23',1,1,1,'2015-11-25 17:15:02'),('3ac77a01-3335-4a8f-be3f-67d52cae83b0\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','humisense','dssvjkl','Humidity',NULL,NULL,0,0,0,'2015-12-01 15:28:48'),('731df8f6-79b3-4c4b-9ea2-0d7045f45a5b\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','fscsv12r13','sdsadas','Temperature','','',0,0,0,'2015-12-01 15:17:21'),('88b654e1-2204-45f7-b6f8-b342e331cc43\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','Hello','Sense','Proximity',NULL,NULL,1,0,0,'2015-12-01 15:31:37'),('d5145196-c3dc-44a8-a806-8fa0a73b93cb\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','poolj','kljklj','Gyroscope',NULL,NULL,1,0,0,'2015-12-01 15:23:59'),('e5facd0f-ac5e-422b-8a91-058189e5fb19\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','gyroscope','This is a humidity sensor','Gyro','100','200',0,1,0,'2015-11-24 19:13:31'),('ee469854-d1e0-4cf7-bd34-69b9010212a7\n','81eb168e-7c5a-4ffd-86da-6ca8cfd0ac22','mobisensor','This is a updated sensor','mobisense','121.23','56.25',1,1,0,'2015-11-25 17:16:15');
/*!40000 ALTER TABLE `physical_sensor_metadata` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-12-01  3:07:04
