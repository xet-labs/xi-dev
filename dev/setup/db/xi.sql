SET foreign_key_checks = 0;
-- MariaDB dump 10.19-11.8.1-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: xi
-- ------------------------------------------------------
-- Server version	11.8.1-MariaDB-2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `blog_media`
--

DROP TABLE IF EXISTS `blog_media`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_media` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `blog_id` bigint(20) unsigned NOT NULL,
  `media_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `blog_id` (`blog_id`),
  KEY `media_id` (`media_id`),
  CONSTRAINT `blog_media_ibfk_1` FOREIGN KEY (`blog_id`) REFERENCES `blogs` (`id`) ON DELETE CASCADE,
  CONSTRAINT `blog_media_ibfk_2` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `blog_media`
--

LOCK TABLES `blog_media` WRITE;
/*!40000 ALTER TABLE `blog_media` DISABLE KEYS */;
set autocommit=0;
/*!40000 ALTER TABLE `blog_media` ENABLE KEYS */;
UNLOCK TABLES;
commit;

--
-- Table structure for table `blogs`
--

DROP TABLE IF EXISTS `blogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `blogs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) unsigned NOT NULL,
  `status` enum('draft','published','published_hidden','archived') NOT NULL DEFAULT 'draft',
  `tags` varchar(255) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `short_title` varchar(255) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `featured_img` varchar(255) DEFAULT NULL,
  `slug` varchar(255) DEFAULT NULL,
  `path` varchar(255) NOT NULL,
  `created_at` timestamp(3) NULL DEFAULT current_timestamp(3),
  `updated_at` timestamp(3) NULL DEFAULT current_timestamp(3) ON UPDATE current_timestamp(3),
  `content` longtext DEFAULT NULL,
  `meta_keywords` text DEFAULT NULL,
  `meta_og` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`meta_og`)),
  `meta_ld` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`meta_ld`)),
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  CONSTRAINT `blogs_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `blogs`
--

LOCK TABLES `blogs` WRITE;
/*!40000 ALTER TABLE `blogs` DISABLE KEYS */;
set autocommit=0;
INSERT INTO `blogs` VALUES
(1,4,'published','[\"SEO\", \"Marketing Insight\"]','What is SEO? A beginner’s guide to search engine optimization','What is SEO ?','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results...','[\"/media/4/img/what-is-seo.webp\"]','what-is-seo','','2024-05-15 06:35:14.000','2025-03-08 05:39:14.401','','what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(2,2,'published','[\"Passive Income\", \"Internet Monetization\"]','Passive income by sharing internet',NULL,'Transform your internet connection into a Passive Income Stream by sharing your internet','[\"/media/2/img/passive-income-by-sharing-internet.webp\"]','passive-income-by-sharing-internet','','2024-05-16 06:36:14.000','2025-02-05 10:38:43.565',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(3,1,'published','[\"mysteriumnode\", \"mysteriumvpn\"]','Mysterium node backup & restore on linux or windows',NULL,'Backup & restore your mysterium node on both windows and linux platforms','[\"/media/1/img/mysterium-node-backup-and-restore.webp\"]','mysterium-node-backup-and-restore','','2024-05-01 06:36:14.000','2025-03-17 16:25:22.137','<p>Backing up your Mysterium node is crucial to ensure that your configs and data are safe in case of system failures or data corruption. This guide will walk you through the steps to back up your Mysterium node on both Windows and Linux platforms.</p>\n\n<h2> Myst node config & data </h2>\n<p> Mysterium node\'s config includes active services, account API, current theme etc, while data includes configs, node’s public address, private key and its encryption details, Message Authentication Code, file id, version etc.</p>\n<p>\n  <b>Config:</b> <code>/etc/mysterium-node/</code>\n  <br>\n  <b>Data:</b> <code>/var/lib/mysterium-node/</code>\n</p>\n\n\n<h2> Backup Myst node on Linux </h2>\n<p>Create a compressed backup with necessary files from <code>/var/lib/mysterium-node/*</code>.</p>\n\n<pre>\n<code class=\"language-bash\"># creates &ltmyst-date.tar.gz&gt compressed backup \nsudo tar --exclude={\"mainnet/*\",\"*/logs/*\"} -cvzf ~/myst-$(date +\"%Y%m%d%H%M%S\").tar.gz -C /var/lib/mysterium-node/</code>\n</pre>\n\n<p>The <code>--exclude={\"mainnet/*\",\"*/logs/*\"}</code> excludes unnecessary files.</p>\n\n\n<h2> Restore Myst node on Linux </h2>\n<p>Replace <code>~/myst-date.tar.gz</code> with actual zip backup of MystNode.</p>\n<pre>\n<code class=\"language-bash\"># extracts &ltmyst-date.tar.gz&gt to \'/var/lib/mysterium-node/\'\nsudo tar -xvzf ~/myst-date.tar.gz -C /var/lib/mysterium-node/ --strip-components=2</code>\n</pre>','how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(4,1,'published_hidden',NULL,'','Does my business need a website? 9 reasons why it does','The Internet powers virtually every aspect of the modern economy, transform your Internet connection into a passive income stream by merely sharing your Internet connection','[\"https://static.wixstatic.com/media/1f6616_88503e0f296f49afa70a53e3812fd92e~mv2.png/v1/fill/w_740,h_489,al_c,q_90,usm_0.66_1.00_0.01,enc_avif,quality_auto/1f6616_88503e0f296f49afa70a53e3812fd92e~mv2.png\"]','passive-income-by-sharing-internet','','2024-05-17 06:36:14.000','2025-02-28 18:25:07.565',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(5,1,'published_hidden','[\"Mysterium Node\"]','','What is a Domain name and why it matters','Backup & restore your mysterium node on windows and linux','[\"https://static.wixstatic.com/media/84b06e_74efe24ffe3f4cccbfb4d23e904498eb~mv2.png/v1/fill/w_740,h_423,al_c,q_85,usm_0.66_1.00_0.01,enc_avif,quality_auto/84b06e_74efe24ffe3f4cccbfb4d23e904498eb~mv2.png\"]','mysterium-node-backup-and-restore','','2024-05-01 01:06:14.000','2025-03-01 07:29:15.379',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(6,2,'published_hidden','[\"Passive Income\", \"Internet Monetization\"]','','10 outstanging site examples','Transform your internet connection into a Passive Income Stream by sharing your internet','[\"https://static.wixstatic.com/media/46e2e0_0b737cde10004f8c8a4deade3e45e47e~mv2.png/v1/fill/w_740,h_489,al_c,q_90,usm_0.66_1.00_0.01,enc_avif,quality_auto/46e2e0_0b737cde10004f8c8a4deade3e45e47e~mv2.png\"]','passive-income-by-sharing-internet','','2024-05-16 01:06:14.000','2025-03-01 07:25:15.861',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(7,4,'published_hidden','[\"SEO\", \"Marketing Insight\"]','','Worlds shortest guide on Ui Ux','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results...','[\"https://static.wixstatic.com/media/5af200_7fe19f9ddee04ecfbaafd51a304008e2~mv2.png/v1/fill/w_740,h_423,al_c,q_85,usm_0.66_1.00_0.01,enc_avif,quality_auto/5af200_7fe19f9ddee04ecfbaafd51a304008e2~mv2.png\"]','what-is-seo','','2024-05-15 01:05:14.000','2025-03-01 07:25:15.862',NULL,'what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(8,4,'published_hidden','[\"SEO\", \"Marketing Insight\"]','','13 vintage websites that showcase timeless retro web design','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results...','[\"/media/local/sites_example.avif\"]','what-is-seo','','2024-05-15 01:05:14.000','2025-03-01 09:21:13.142',NULL,'what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(9,2,'published_hidden','[\"Passive Income\", \"Internet Monetization\"]','','How to build a software engineering, coding and development portfolio','Transform your internet connection into a Passive Income Stream by sharing your internet','[\"/media/local/case_study.avif\"]','passive-income-by-sharing-internet','','2024-05-16 01:06:14.000','2025-03-01 09:24:30.036',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(10,1,'published_hidden','[\"Mysterium Node\"]','','21 social media content ideas, plus best posts that go viral','21 social media content ideas, plus best posts that go viral','[\"/media/local/site_building_guide.avif\"]','mysterium-node-backup-and-restore','','2024-05-01 01:06:14.000','2025-03-01 09:27:04.351',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(11,4,'published_hidden','[\"SEO\", \"Marketing Insight\"]','What is SEO? A beginner’s guide to search engine optimization','What is SEO ?','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results...','[\"/media/4/img/what-is-seo.webp\"]','what-is-seo','','2024-05-15 06:35:14.000','2025-03-01 09:30:23.623',NULL,'what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(12,2,'published_hidden','[\"Passive Income\", \"Internet Monetization\"]','Passive income by sharing internet',NULL,'Transform your internet connection into a Passive Income Stream by sharing your internet','[\"/media/2/img/passive-income-by-sharing-internet.webp\"]','passive-income-by-sharing-internet','','2024-05-16 06:36:14.000','2025-03-01 09:30:23.626',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(13,1,'published_hidden','[\"Mysterium Node\"]','Mysterium node backup & restore',NULL,'Backup & restore your mysterium node on windows and linux','[\"/media/1/img/mysterium-node-backup-and-restore.webp\"]','mysterium-node-backup-and-restore','','2024-05-01 06:36:14.000','2025-03-01 09:30:23.627',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(14,1,'published_hidden',NULL,'','Does my business need a website? 9 reasons why it does','The Internet powers virtually every aspect of the modern economy, transform your Internet connection into a passive income stream by merely sharing your Internet connection ','[\"https://static.wixstatic.com/media/1f6616_88503e0f296f49afa70a53e3812fd92e~mv2.png/v1/fill/w_740,h_489,al_c,q_90,usm_0.66_1.00_0.01,enc_avif,quality_auto/1f6616_88503e0f296f49afa70a53e3812fd92e~mv2.png\"]','passive-income-by-sharing-internet','','2024-05-17 06:36:14.000','2025-03-01 09:30:59.101',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(15,1,'published_hidden','[\"Mysterium Node\"]','','What is a Domain name and why it matters','Backup & restore your mysterium node on windows and linux ','[\"https://static.wixstatic.com/media/84b06e_74efe24ffe3f4cccbfb4d23e904498eb~mv2.png/v1/fill/w_740,h_423,al_c,q_85,usm_0.66_1.00_0.01,enc_avif,quality_auto/84b06e_74efe24ffe3f4cccbfb4d23e904498eb~mv2.png\"]','mysterium-node-backup-and-restore','','2024-05-01 01:06:14.000','2025-03-01 09:31:02.508',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(16,2,'published_hidden','[\"Passive Income\", \"Internet Monetization\"]','','10 outstanging site examples','Transform your internet connection into a Passive Income Stream by sharing your internet ','[\"https://static.wixstatic.com/media/46e2e0_0b737cde10004f8c8a4deade3e45e47e~mv2.png/v1/fill/w_740,h_489,al_c,q_90,usm_0.66_1.00_0.01,enc_avif,quality_auto/46e2e0_0b737cde10004f8c8a4deade3e45e47e~mv2.png\"]','passive-income-by-sharing-internet','','2024-05-16 01:06:14.000','2025-03-01 09:31:05.612',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(17,4,'published_hidden','[\"SEO\", \"Marketing Insight\"]','','Worlds shortest guide on Ui Ux','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results... ','[\"https://static.wixstatic.com/media/5af200_7fe19f9ddee04ecfbaafd51a304008e2~mv2.png/v1/fill/w_740,h_423,al_c,q_85,usm_0.66_1.00_0.01,enc_avif,quality_auto/5af200_7fe19f9ddee04ecfbaafd51a304008e2~mv2.png\"]','what-is-seo','','2024-05-15 01:05:14.000','2025-03-01 09:31:09.997',NULL,'what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(18,4,'published_hidden','[\"SEO\", \"Marketing Insight\"]','','13 vintage websites that showcase timeless retro web design','SEO is the strategic art of enhancing a website\'s visibility on search engines like Google, Bing, and Yahoo. It involves optimizing various elements to rank higher in search results... ','[\"/media/local/sites_example.avif\"]','what-is-seo','','2024-05-15 01:05:14.000','2025-03-01 09:31:12.499',NULL,'what is SEO, Search Engine Optimization, what is website optimization, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(19,2,'published_hidden','[\"Passive Income\", \"Internet Monetization\"]','','How to build a software engineering, coding and development portfolio','Transform your internet connection into a Passive Income Stream by sharing your internet ','[\"/media/local/case_study.avif\"]','passive-income-by-sharing-internet','','2024-05-16 01:06:14.000','2025-03-01 09:31:15.419',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL),
(20,1,'published_hidden','[\"Mysterium Node\"]','','21 social media content ideas, plus best posts that go viral','21 social media content ideas, plus best posts that go viral ','[\"/media/local/site_building_guide.avif\"]','mysterium-node-backup-and-restore','','2024-05-01 01:06:14.000','2025-03-01 09:31:18.924',NULL,'how to make money online, ways to make maney online, xet, xet industries, XetIndustries, XetIndustries blog, Rishikesh Prasad',NULL,NULL);
/*!40000 ALTER TABLE `blogs` ENABLE KEYS */;
UNLOCK TABLES;
commit;

--
-- Table structure for table `media`
--

DROP TABLE IF EXISTS `media`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `media` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` enum('audio','image','video','document') NOT NULL,
  `hash` varchar(255) NOT NULL,
  `filename` varchar(255) NOT NULL,
  `path` varchar(255) NOT NULL,
  `size` bigint(20) unsigned NOT NULL,
  `format` varchar(50) NOT NULL,
  `width` int(10) unsigned DEFAULT NULL,
  `height` int(10) unsigned DEFAULT NULL,
  `duration` int(10) unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `media`
--

LOCK TABLES `media` WRITE;
/*!40000 ALTER TABLE `media` DISABLE KEYS */;
set autocommit=0;
/*!40000 ALTER TABLE `media` ENABLE KEYS */;
UNLOCK TABLES;
commit;

--
-- Table structure for table `user_media`
--

DROP TABLE IF EXISTS `user_media`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_media` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) unsigned NOT NULL,
  `media_id` bigint(20) unsigned NOT NULL,
  `tags` varchar(255) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`,`media_id`),
  KEY `media_id` (`media_id`),
  CONSTRAINT `user_media_ibfk_1` FOREIGN KEY (`uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE,
  CONSTRAINT `user_media_ibfk_2` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_media`
--

LOCK TABLES `user_media` WRITE;
/*!40000 ALTER TABLE `user_media` DISABLE KEYS */;
set autocommit=0;
/*!40000 ALTER TABLE `user_media` ENABLE KEYS */;
UNLOCK TABLES;
commit;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `uid` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `email` varchar(50) NOT NULL,
  `verified` tinyint(1) NOT NULL DEFAULT 0,
  `password` varchar(255) NOT NULL,
  `role` enum('admin','dev','user') DEFAULT 'user',
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `last_login` datetime DEFAULT NULL,
  `status` enum('active','inactive','suspended','deleted') DEFAULT 'active',
  `profile_img` varchar(255) DEFAULT NULL,
  `address` text DEFAULT NULL,
  `phone_no` varchar(20) DEFAULT NULL,
  `dob` date DEFAULT NULL,
  `config` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`config`)),
  `remember_token` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
set autocommit=0;
INSERT INTO `users` VALUES
(1,'xet','Rishikesh Prasad','rishikeshprasad@xetindustries.com',1,'$2y$12$3zXdoKP91LvYctZUrN9cYOfBV8TyAEeoQK5ADAVQuBINSBt1nJRTK','user','2024-07-27 04:05:01','2025-02-01 13:37:18',NULL,'active','/media/1/profile/xet.jpg',NULL,NULL,NULL,NULL,NULL),
(2,'zet','Zet Ohio','zet@g.com',1,'$2y$12$3zXdoKP91LvYctZUrN9cYOfBV8TyAEeoQK5ADAVQuBINSBt1nJRTK','user','2024-07-27 04:13:44','2025-02-01 13:37:18',NULL,'active','/media/2/profile/zet.jpg',NULL,NULL,NULL,NULL,NULL),
(4,'cristine','Cristine Lepcha','cr@g.com',0,'$2y$12$3zXdoKP91LvYctZUrN9cYOfBV8TyAEeoQK5ADAVQuBINSBt1nJRTK','user','2024-07-27 04:13:44','2025-02-01 13:37:18',NULL,'active','/media/4/profile/cristine.jpg',NULL,NULL,NULL,NULL,NULL),
(5,'t1','t1','t1@g.com',0,'$2y$12$3zXdoKP91LvYctZUrN9cYOfBV8TyAEeoQK5ADAVQuBINSBt1nJRTK','user','2024-07-27 04:13:44','2024-11-29 21:22:51',NULL,'active','/media/3/profile/zet.jpg',NULL,NULL,NULL,NULL,NULL),
(7,'wbeiliq41','UyXhRcENW','wbeiliq41@gmail.com',0,'$2y$12$pTONJnvKco8ksZaLcXvWKuXkzSlMksBvipuq8xueLxfq15RlYDUne','user','2025-03-16 06:02:15','2025-03-16 06:02:15',NULL,'active',NULL,NULL,NULL,NULL,NULL,NULL),
(9,'hey','hey','hey@gmail.com',0,'$2y$12$QMOBTM62Y2YNUkv.wGdBseZuF07BEBZMxr.3PzT5nTm6QDP/hd/z6','user','2025-03-22 14:47:46','2025-03-22 14:47:46',NULL,'active',NULL,NULL,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
commit;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-04-02  2:11:58
