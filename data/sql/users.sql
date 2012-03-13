-- 
-- Table structure for table `users`
-- 

CREATE TABLE `users` (
  `id` varchar(255) character set utf8 collate utf8_roman_ci NOT NULL,
  `username` varchar(255) character set utf8 collate utf8_roman_ci NOT NULL,
  `emailAddress` varchar(255) character set utf8 collate utf8_roman_ci NOT NULL,
  `password` varchar(255) character set utf8 collate utf8_roman_ci NOT NULL,
  `role` enum('superadmin','user') NOT NULL default 'user',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `emailAddress` (`emailAddress`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8;

