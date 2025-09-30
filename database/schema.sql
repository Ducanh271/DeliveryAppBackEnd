-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: DeliveryAppDB
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Images`
--

DROP TABLE IF EXISTS `Images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Images` (
  `id` int NOT NULL AUTO_INCREMENT,
  `url` varchar(500) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `public_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=69 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Images`
--

LOCK TABLES `Images` WRITE;
/*!40000 ALTER TABLE `Images` DISABLE KEYS */;
INSERT INTO `Images` VALUES (28,'https://www.coca-cola.com/content/dam/brands/us/coca-cola/coca-cola-logo.png','2025-09-24 09:46:55',''),(29,'https://upload.wikimedia.org/wikipedia/commons/2/27/Coca_Cola_Flasche_-_Original_Taste.jpg','2025-09-24 09:46:55',''),(30,'https://product.hstatic.net/200000534989/product/dsc08341-enhanced-nr_1_e6d5d0a13c8f42c2bd7cea59e03ce199_master.jpg','2025-09-24 09:46:55',''),(31,'https://www.pepsi.com/en-us/uploads/images/twil-can.png','2025-09-24 09:46:55',''),(32,'https://www.oreo.com/images/hero/oreo-original.png','2025-09-24 09:46:55',''),(33,'https://cdn.lottemart.vn/media/catalog/product/cache/9b5f86ccf0bb6d794da3fb554015eb8c/s/t/sting_dau_330ml.jpg','2025-09-24 09:46:55',''),(34,'https://cdn.thtrue.vn/wp-content/uploads/2022/04/tra-xanh-0-do.jpg','2025-09-24 09:46:55',''),(35,'https://www.oreo.com/images/hero/oreo-original.png','2025-09-24 09:46:55',''),(36,'https://cdn.lottemart.vn/media/catalog/product/cache/9b5f86ccf0bb6d794da3fb554015eb8c/s/t/sting_dau_330ml.jpg','2025-09-24 09:46:55',''),(37,'https://bavifoods.com/thumbs/740x740x1/upload/product/cam-ep-5018.jpg','2025-09-24 09:46:55',''),(38,'https://bavifoods.com/thumbs/740x740x1/upload/product/cam-ep-5018.jpg','2025-09-24 09:46:55',''),(39,'https://cdn-i.vtcnews.vn/files/news/2019/01/22/-145625.jpg','2025-09-24 09:46:55',''),(40,'https://media.baobinhphuoc.com.vn/upload/news/5_2023/img_8476_06413001052023.jpeg','2025-09-24 09:46:55',''),(41,'https://media.vov.vn/sites/default/files/styles/large/public/2023-06/nuoc_chanh_5.jpg','2025-09-24 09:46:55',''),(42,'https://suckhoedoisong.qltns.mediacdn.vn/324455921873985536/2022/4/18/uong-nuoc-chanh-moi-ngay-co-tot-khong-va-uong-khi-nao-chanh1-1592466583-666-width1024height768-16502694973802093569436.jpg','2025-09-24 09:46:55',''),(43,'Chi-em-thi-nhau-lung-mua-dua-xiem-ve-uong-sau-tiem-phong-cua-hang-moi-ngay-ban-5000-qua-1-1631524190-570-width650height431.jpg (650×431)','2025-09-24 09:46:55',''),(44,'medium_20200513_094458_574364_nuoc_dua_max_1800x1800_jpg_095dc5e7ad.jpg (750×563)','2025-09-24 09:46:55',''),(45,'coconut-water-benefits-17218412875751213756362.jpg (800×562)','2025-09-24 09:46:55',''),(46,'https://file.hstatic.net/1000199715/file/uong-sua-sau-sinh-1_90f6b928e6084e7e87c4e7a89e1b1be3_grande.jpg','2025-09-24 09:46:55',''),(47,'khi_nao_nen_cho_be_uong_sua_1_4401cf044a.jpg (800×600)','2025-09-24 09:46:55',''),(48,'https://suckhoedoisong.qltns.mediacdn.vn/324455921873985536/2025/2/21/dau-nanh-1-1740125251401155246723.jpg','2025-09-24 09:46:55',''),(49,'glass-soy-milk_20dc83bb32164c49bd11a7d7b60b717b_grande.jpg (600×377)','2025-09-24 09:46:55',''),(50,'may-lam-sua-dau-nanh-1-1412734006024.jpg (500×455)','2025-09-24 09:46:55',''),(51,'https://baothainguyen.vn/file/e7837c027f6ecd14017ffa4e5f2a0e34/032023/1-boba-tea-recipe-using-fresh-tapioca-pearls-1024x1024-1677809524112848165864_20230305161118.jpeg','2025-09-24 09:46:55',''),(52,'https://www.cet.edu.vn/wp-content/uploads/2018/04/tra-sua-tu-lam.jpg','2025-09-24 09:46:55',''),(53,'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRXtYkBiYM_AEWth56eq5VeEwFGlnh_cYm7cw&s','2025-09-28 18:34:00',''),(54,'https://billballcoffeetea.com/upload/product/img3589-3619-8202.jpg','2025-09-28 18:34:00',''),(55,'https://vcdn1-suckhoe.vnecdn.net/2023/02/01/iced-coffee-table-jpeg-1675223-7169-5352-1675223880.jpg?w=1200&h=0&q=100&dpr=1&fit=crop&s=ExRy7hbEHS2p2f2oOgWCjA','2025-09-28 18:34:00',''),(56,'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQie4Hw_GApMLyRn34hNhlXh46_33_56ZcfMA&s','2025-09-28 18:34:00',''),(57,'https://cdn.tgdd.vn/2020/07/CookProductThumb/59-620x620-3.jpg','2025-09-28 18:34:00',''),(58,'https://product.hstatic.net/200000848723/product/combo_2_burger_e7394736b32d499e8b9482c04030f5f5_master.jpg','2025-09-28 18:34:00',''),(59,'https://www.bluestone.com.vn/blogs/vao-bep/chien-khoai-tay-bang-noi-chien-khong-dau?srsltid=AfmBOoqJ8H9yAh7LTcp--k4VJ6EP2X2OuqePxFBm3Y8K1ygpG_FpQDcN','2025-09-28 18:34:00',''),(60,'https://img.giftpop.vn/brand/LOTTERIA/1PEMP2010312501_BASIC_origin.jpg','2025-09-28 18:34:00',''),(61,'https://cdn.tgdd.vn/Files/2020/03/02/1239549/2-cong-thuc-lam-banh-hotdog-xuc-xich-hotdog-pho-mai-han-quoc-gay-nghien-14-760x367.png','2025-09-28 18:34:00',''),(62,'https://checkinvietnam.vtc.vn/media/20211221/files/pizza-xuc-xich-pho-mai-vuong.jpg','2025-09-28 18:34:00',''),(63,'https://i.ytimg.com/vi/ng3vo1RmeyQ/maxresdefault.jpg','2025-09-28 18:34:00',''),(64,'https://storage.googleapis.com/onelife-public/blog.onelife.vn/2021/10/cach-lam-banh-mi-sandwich-trung-jambon-mon-an-sang-349515833958.jpg','2025-09-28 18:34:00',''),(65,'https://cdnv2.tgdd.vn/bhx-static/bhx/Products/Images/7259/332717/bhx/frame-3475095-2-1_202412022211137044.jpg','2025-09-28 18:34:00',''),(66,'https://cdn.tgdd.vn/Files/2022/03/07/1418886/9-cach-lam-salad-tron-mayonnaise-giam-can-tai-nha-hieu-qua-202203071357195806.jpg','2025-09-28 18:34:00',''),(68,'https://res.cloudinary.com/dytix5ybu/image/upload/v1759138848/product/TraDaoCamSa.jpg.jpg','2025-09-29 17:03:26','41');
/*!40000 ALTER TABLE `Images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ProductImages`
--

DROP TABLE IF EXISTS `ProductImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ProductImages` (
  `product_id` int NOT NULL,
  `image_id` int NOT NULL,
  `is_main` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`product_id`,`image_id`),
  KEY `fk_image` (`image_id`),
  CONSTRAINT `fk_image` FOREIGN KEY (`image_id`) REFERENCES `Images` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_product` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ProductImages`
--

LOCK TABLES `ProductImages` WRITE;
/*!40000 ALTER TABLE `ProductImages` DISABLE KEYS */;
INSERT INTO `ProductImages` VALUES (16,37,1),(16,38,0),(16,39,0),(17,40,1),(17,41,0),(17,42,0),(19,46,1),(19,47,0),(20,48,1),(20,49,0),(20,50,0),(21,51,1),(21,52,0),(22,53,1),(23,54,1),(24,55,1),(25,56,1),(26,57,1),(27,58,1),(29,60,1),(30,61,1),(31,62,1),(32,63,1),(33,64,1),(34,65,1),(35,66,1),(41,68,1);
/*!40000 ALTER TABLE `ProductImages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Products`
--

DROP TABLE IF EXISTS `Products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` decimal(10,2) NOT NULL,
  `qty_initial` int DEFAULT '0',
  `qty_sold` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `Products_chk_1` CHECK ((`price` >= 0)),
  CONSTRAINT `Products_chk_2` CHECK ((`qty_initial` >= 0)),
  CONSTRAINT `Products_chk_3` CHECK ((`qty_sold` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Products`
--

LOCK TABLES `Products` WRITE;
/*!40000 ALTER TABLE `Products` DISABLE KEYS */;
INSERT INTO `Products` VALUES (16,'Nước Cam Ép','Nước cam tươi nguyên chất',15000.00,100,35,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(17,'Nước Chanh Tươi','Nước chanh mát lạnh giải khát',12000.00,90,28,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(19,'Sữa Tươi','Sữa tươi tiệt trùng nguyên chất',18000.00,80,30,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(20,'Sữa Đậu Nành','Thức uống từ đậu nành bổ dưỡng',12000.00,90,22,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(21,'Trà Sữa Trân Châu','Trà sữa kèm trân châu dai ngon',35000.00,100,50,'2025-09-24 02:46:55','2025-09-24 09:46:55'),(22,'Trà Đào Cam Sả','Trà đào cam sả mát lạnh',30000.00,60,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(23,'Cà Phê Đen','Cà phê đen nguyên chất',20000.00,80,35,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(24,'Cà Phê Sữa','Cà phê sữa đá truyền thống',25000.00,90,40,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(25,'Sinh Tố Bơ','Sinh tố bơ béo ngậy',40000.00,50,18,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(26,'Burger Bò Phô Mai','Bánh burger bò kèm phô mai tan chảy',45000.00,50,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(27,'Burger Gà Giòn','Bánh burger gà chiên giòn',40000.00,60,25,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(29,'Gà Rán 2 Miếng','Gà rán giòn rụm, hương vị đặc trưng',60000.00,80,30,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(30,'Hotdog Xúc Xích','Bánh mì kẹp xúc xích và tương cà',30000.00,70,20,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(31,'Pizza Phô Mai','Pizza nhỏ phủ phô mai mozzarella',70000.00,40,15,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(32,'Pizza Hải Sản','Pizza hải sản tươi ngon',85000.00,35,10,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(33,'Sandwich Thịt Nguội','Bánh sandwich kẹp thịt nguội và rau',35000.00,60,22,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(34,'Mì Ý Sốt Bò Bằm','Mì Ý sốt cà chua bò bằm',65000.00,45,18,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(35,'Salad Rau Trộn','Salad rau củ tươi mát',30000.00,50,12,'2025-09-28 11:34:01','2025-09-28 18:34:00'),(41,'Trà Đào Cam Sả','Trà đào cam xả 100% làm từ thiên nhiên',25000.00,100,0,'2025-09-29 10:03:24','2025-09-29 17:03:24');
/*!40000 ALTER TABLE `Products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ReviewImages`
--

DROP TABLE IF EXISTS `ReviewImages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ReviewImages` (
  `review_id` int NOT NULL,
  `image_id` int NOT NULL,
  UNIQUE KEY `uq_review_image` (`review_id`,`image_id`),
  KEY `fk_reviewimages_image` (`image_id`),
  CONSTRAINT `fk_reviewimages_image` FOREIGN KEY (`image_id`) REFERENCES `Images` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reviewimages_review` FOREIGN KEY (`review_id`) REFERENCES `Reviews` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ReviewImages`
--

LOCK TABLES `ReviewImages` WRITE;
/*!40000 ALTER TABLE `ReviewImages` DISABLE KEYS */;
/*!40000 ALTER TABLE `ReviewImages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Reviews`
--

DROP TABLE IF EXISTS `Reviews`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Reviews` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_id` int NOT NULL,
  `user_id` int NOT NULL,
  `order_id` int NOT NULL,
  `rate` int NOT NULL,
  `content` varchar(1000) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_product_order` (`user_id`,`product_id`,`order_id`),
  KEY `fk_review_product` (`product_id`),
  KEY `fk_review_order` (`order_id`),
  CONSTRAINT `fk_review_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_review_product` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`),
  CONSTRAINT `fk_review_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `Reviews_chk_1` CHECK ((`rate` between 1 and 5))
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Reviews`
--

LOCK TABLES `Reviews` WRITE;
/*!40000 ALTER TABLE `Reviews` DISABLE KEYS */;
/*!40000 ALTER TABLE `Reviews` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `sender_id` int NOT NULL,
  `receiver_id` int NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `fk_order` (`order_id`),
  KEY `fk_sender` (`sender_id`),
  KEY `fk_receiver` (`receiver_id`),
  CONSTRAINT `fk_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `fk_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (5,6,39,37,'Đơn hàng đã sẵn sàng giao','2025-09-30 11:15:25',0),(6,6,39,37,'Đơn hàng đã sẵn sàng giao','2025-09-30 16:33:21',0),(7,6,39,37,'chao lan 1','2025-09-30 16:33:38',0),(8,6,39,37,'chao lan 2','2025-09-30 16:33:42',0),(9,6,39,37,'chao lan 2','2025-09-30 16:33:45',0),(10,6,37,39,'chao lai lan 1','2025-09-30 16:34:07',0),(11,6,39,37,'hello customer','2025-09-30 16:46:58',0),(12,6,37,39,'hi shipper','2025-09-30 16:47:20',0);
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_items`
--

DROP TABLE IF EXISTS `order_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `order_id` int NOT NULL,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `order_id` (`order_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE,
  CONSTRAINT `order_items_chk_1` CHECK ((`quantity` > 0)),
  CONSTRAINT `order_items_chk_2` CHECK ((`price` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_items`
--

LOCK TABLES `order_items` WRITE;
/*!40000 ALTER TABLE `order_items` DISABLE KEYS */;
INSERT INTO `order_items` VALUES (8,7,16,7,15000.00);
/*!40000 ALTER TABLE `order_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `shipper_id` int DEFAULT NULL,
  `payment_status` enum('unpaid','paid','refunded') DEFAULT 'unpaid',
  `order_status` enum('pending','processing','shipping','delivered','cancelled') DEFAULT 'pending',
  `latitude` decimal(10,8) NOT NULL,
  `longitude` decimal(11,8) NOT NULL,
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `thumbnail_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_orders_thumbnail` (`thumbnail_id`),
  CONSTRAINT `fk_orders_thumbnail` FOREIGN KEY (`thumbnail_id`) REFERENCES `Images` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (6,37,39,'unpaid','shipping',21.02851100,105.80481700,149000.00,28,'2025-09-27 02:01:50','2025-09-28 11:46:30'),(7,37,39,'unpaid','shipping',21.02851100,105.80481700,165000.00,34,'2025-09-27 03:21:13','2025-09-28 11:46:44');
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refresh_tokens`
--

DROP TABLE IF EXISTS `refresh_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refresh_tokens` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `token` text NOT NULL,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `refresh_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
INSERT INTO `refresh_tokens` VALUES (10,37,'geBwj4xvzuErapojE5NZCjizXGsC7jWasGHYoZWUss8=','2025-10-03 00:42:50','2025-09-26 07:42:49','2025-09-26 07:42:49'),(11,37,'rv3jLza9cK-oFrl83jcVpqIvNDT6hZwH01Ec-g_uJ_s=','2025-10-04 01:50:06','2025-09-27 08:50:06','2025-09-27 08:50:06'),(12,35,'Gdu2gjSJRWxZZRWGSOsjkiZDDDpQ9XpUOAUDiEtZbA4=','2025-10-04 01:55:40','2025-09-27 08:55:40','2025-09-27 08:55:40'),(13,37,'kNofYfHuCFaMZ0K5FAodqn5PATSqINRUO4kJNFBbxAQ=','2025-10-04 01:58:31','2025-09-27 08:58:30','2025-09-27 08:58:30'),(14,37,'wMTUjtdha0VVaT_Y3k-lks5IP5ZSbjfvPwcKULmJEo4=','2025-10-04 03:20:21','2025-09-27 10:20:20','2025-09-27 10:20:20'),(15,36,'pBLST5x-fBHU6AkDLWGURDU68Cd3POF0vYDXoLYKd9A=','2025-10-04 03:25:05','2025-09-27 10:25:05','2025-09-27 10:25:05'),(16,36,'gn3KB5w6b8Rv2NxG3xVqqDtGgkUK288ivoOmpgxzCm8=','2025-10-04 10:23:02','2025-09-27 17:23:01','2025-09-27 17:23:01'),(17,39,'RWQVXjtKbUO40Coc0Tdd8klwy_tU9cBdelFDMq-50Lk=','2025-10-05 04:12:59','2025-09-28 11:12:59','2025-09-28 11:12:59'),(18,36,'p2vqjFNFfuLcr-NydJvCry28Ajhry14nnyKaly2Zgm8=','2025-10-05 07:42:04','2025-09-28 14:42:03','2025-09-28 14:42:03'),(19,36,'GoCJywlU1pcJqtkloeH1XL12oBCX204UkJWrnftUXFI=','2025-10-05 08:13:17','2025-09-28 15:13:17','2025-09-28 15:30:27'),(20,36,'BCE2qwZsOMi4dUBQ-4424KqOb8iw450PyJt2DUNvqaI=','2025-10-05 08:36:02','2025-09-28 15:36:01','2025-09-28 15:38:27'),(21,36,'JvPW-n5swM9VnPenXYuF1JG-Va1gHAT83I7QZNx6xjo=','2025-10-05 08:45:51','2025-09-28 15:45:50','2025-09-28 15:45:50'),(22,36,'9xkZNF-dZBDmQMcBHX68VvQk_sv-_IZLhrrHLdDEIHU=','2025-10-05 08:49:16','2025-09-28 15:49:15','2025-09-28 15:49:15'),(23,36,'JV00Q9bsSEExXZcPpkUJUpDNxMoiWRa-RDvEebPLsmI=','2025-10-05 08:56:57','2025-09-28 15:56:57','2025-09-28 16:03:12'),(24,36,'rKknYfShgf5J5tAR2C_yxRCz51xUAuBSCyEi1KjL1_w=','2025-10-05 09:03:56','2025-09-28 16:03:56','2025-09-28 16:05:05'),(25,36,'Kt1f2CUnQq87I0aQF5xt6YulCj2_bXikq7NYoPRHKKE=','2025-10-05 09:05:43','2025-09-28 16:05:42','2025-09-28 16:05:52'),(26,36,'0iZKc2rwMCsvekClA-Ni7fonGd0iKHC25K5GOdhbjAE=','2025-10-05 09:06:20','2025-09-28 16:06:19','2025-09-28 16:06:30'),(27,36,'5dt26J4APOBLCtKwRSmQ67Dg85slWeL7xuguseY0xo0=','2025-10-05 09:08:52','2025-09-28 16:08:51','2025-09-28 16:36:15'),(28,36,'Ye7AN6HJkeOzRMWkaTQO21fv2Hqb1gYr1D-no8eWZfI=','2025-10-05 09:42:16','2025-09-28 16:42:15','2025-09-28 16:45:39'),(29,36,'nXDENnH3ZpeFpm5YuPc358-GL9JamxTnAHkM0_kGDl8=','2025-10-05 11:34:54','2025-09-28 18:34:54','2025-09-28 18:37:33'),(30,36,'70yiRGssKNrOigjkZwYx8V0yfBW7B9Z_CJaQAggyu7c=','2025-10-06 09:33:32','2025-09-29 16:33:31','2025-09-29 16:40:46'),(31,36,'zOOUswQ944JIgk7M3G0pxVcVHI5gJMo-dYPpUr-WZwk=','2025-10-06 09:43:33','2025-09-29 16:43:32','2025-09-29 16:43:32'),(32,36,'O-WCA1lofuf4BXY0CyKMMp0HoEnaeLqvvAMFvwoEIFg=','2025-10-06 09:52:54','2025-09-29 16:52:54','2025-09-29 16:52:54'),(33,36,'mQgqI2r9DQe6-LJ2K-so6gj2Sms2RDe2d6xnmeVmXLw=','2025-10-06 10:01:03','2025-09-29 17:01:02','2025-09-29 17:01:02'),(34,36,'3d1-EODJgNUNEZPceVdniIugc3JFtUPHhaXbYEtWJ4I=','2025-10-06 10:25:09','2025-09-29 17:25:08','2025-09-29 17:25:08'),(35,36,'thSAQzJ2yYbUzCvgdSTvzT8Jw7UDCBaHQd-z3FvRX3Y=','2025-10-06 10:25:19','2025-09-29 17:25:19','2025-09-29 17:25:19'),(36,36,'vGVZp_YtN38buciK_wzcIijBlyDFP57U_rbjyGgHlyQ=','2025-10-06 10:25:22','2025-09-29 17:25:22','2025-09-29 17:25:22'),(37,39,'dCXyOAps2p6qNvInm-0WeVnj0lMJA3j2YH0JUPE3HuU=','2025-10-06 08:41:52','2025-09-29 15:41:52','2025-09-29 15:41:52'),(38,39,'XSxKlpsSqD5Jtx5uGpKS4OYqN_I2Gwtlqw9TSHQTMpw=','2025-10-06 09:13:55','2025-09-29 16:13:54','2025-09-29 16:13:54'),(39,39,'3Mgp2Is0WRRi5rjyQCmOPonxGNusRlBMv0nCxz8jXQM=','2025-10-06 09:41:45','2025-09-29 16:41:44','2025-09-29 16:41:44'),(40,39,'T04acPhglJaPQ1fCrblB2UTjmwvXf-6BtA9wxacTkp0=','2025-10-06 09:45:49','2025-09-29 16:45:48','2025-09-29 16:45:48'),(41,37,'zfpfs2dRyh3VZ9NoOcntnfQKtDGv0oLZ1-dbLZOBYsk=','2025-10-06 09:51:27','2025-09-29 16:51:26','2025-09-29 16:51:26'),(42,37,'uBUOn3Ed5e9Z0q3fAbUDSMTyaZf7cMA2LsZzd1MEdGE=','2025-10-06 10:10:27','2025-09-29 17:10:27','2025-09-29 17:10:27'),(43,37,'JKS58CembU0X8ZqYRu7iLRRbWjUgRq-3xfXESVuC93c=','2025-10-06 10:11:18','2025-09-29 17:11:18','2025-09-29 17:11:18'),(44,37,'p9479QSmBp_xOVQR2RYj1CW0HyVO8zqrI02cimpMTwM=','2025-10-06 10:13:01','2025-09-29 17:13:01','2025-09-29 17:13:01'),(45,37,'BMMAtOgkQ7_TjGz2Mbx8Fsqpxgb81UCa-fLU6Y_kG3Y=','2025-10-06 10:16:56','2025-09-29 17:16:56','2025-09-29 17:16:56'),(46,37,'2bxFyPwggEV42bXsyqDDqPQF4Vn1-WUoq-GUYnjg9Ds=','2025-10-06 10:17:40','2025-09-29 17:17:40','2025-09-29 17:17:40'),(47,37,'8RXgPsmuCqXu3jieVbVPdtXeOV5N72TPbJsScBeBODk=','2025-10-06 10:18:55','2025-09-29 17:18:54','2025-09-29 17:18:54'),(48,37,'qqLRyeS8FFhz0c2ieSsTkIGGQaxoOz5ty74Lxd4q0ak=','2025-10-06 10:19:06','2025-09-29 17:19:05','2025-09-29 17:19:05'),(49,37,'46rrbdAsJPv499hC9LieADWwyIdpylv7CB-dokNQ78E=','2025-10-06 10:26:58','2025-09-29 17:26:58','2025-09-29 17:26:58'),(50,37,'oQRF5shqjMeFqPD3ZLzhekvZWrM1AKGzzx0_klprpEg=','2025-10-06 10:39:59','2025-09-29 17:39:59','2025-09-29 17:39:59'),(51,39,'8lbVoAz2RxzBVimcN9n8hXOONGfW0m-fBjir-K2IWUA=','2025-10-06 10:54:18','2025-09-29 17:54:18','2025-09-29 17:54:18'),(52,39,'GUaUxi9hpBC_PkbrrHWgVzgOu60CpABmB7aoG20QyEA=','2025-10-07 03:20:39','2025-09-30 10:20:39','2025-09-30 10:20:39'),(53,37,'TZcNzWGAlTQaOYyO3cTvmmDl3eCA_TPKSAmPAidt0-I=','2025-10-07 03:21:14','2025-09-30 10:21:14','2025-09-30 10:21:14'),(54,37,'oc7jdnTRD-ScB627T2nuaYmPY_KnEF6piuc9Dvha_0M=','2025-10-07 09:27:20','2025-09-30 16:27:19','2025-09-30 16:27:19'),(55,39,'xqyVxJZ1C5x7_yiT2N_wAR4yRMUeTBUACrrbUQyLPNI=','2025-10-07 09:27:57','2025-09-30 16:27:56','2025-09-30 16:27:56');
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` text NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `role` enum('customer','shipper','supplier','admin') NOT NULL DEFAULT 'customer',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `otp_code` varchar(10) DEFAULT NULL,
  `otp_expires_at` timestamp NULL DEFAULT NULL,
  `is_verified` tinyint(1) DEFAULT '0',
  `reset_otp` varchar(10) DEFAULT NULL,
  `reset_otp_expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (35,'Admin','admin@example.com','$2a$10$OeGOOr1OBOm6VmxnCJdbOej0325iqiiSQpMtQdQaGSXF7DbRRcmQe','0000000000','Admin Address','admin','2025-09-24 09:36:35',NULL,NULL,1,NULL,NULL),(36,'Admin','admin@gmail.com','$2a$10$KG.aMLHJSAHMb2QwneUVwOZ2zISBuLm5K.15hTT9O6DC.o6CQmbHy','','','admin','2025-09-24 02:37:46',NULL,NULL,1,NULL,NULL),(37,'User1','user1@gmail.com','$2a$10$s50wUtykgFKNELyjFf2z9.PIO82CteRVDu9ATLmOgj0ns0doMOMqW','0987777777','Ha Noi','customer','2025-09-26 00:40:40','529084','2025-09-26 00:50:40',1,NULL,NULL),(38,'admin01','admin1@gmail.com','$2a$10$o8m4W5wobYd10hBL/uemL.F0OtqY0JKGZk6.0cb3dkiy7M4tpOlNu','0987777777','Ha Noi','admin','2025-09-27 01:53:09',NULL,NULL,1,NULL,NULL),(39,'Shipper1','shipper1@gmail.com','$2a$10$flof0WS10vwudJj7394emO4v4ZgIp9HAxi3Y6Gs0edU4mP11tPs3i','','','shipper','2025-09-27 10:23:26','563546','2025-09-28 04:21:28',1,NULL,NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-30 23:49:55
