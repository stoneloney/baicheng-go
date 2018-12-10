/* 用户表 */
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `password` varchar(255) NOT NULL,
  `logintime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;


CREATE TABLE IF NOT EXISTS `data` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` varchar(32) NOT NULL,
  `type` tinyint(1) DEFAULT 0,
  `duration` tinyint(1) DEFAULT 0,
  `director` tinyint(1) DEFAULT 0,
  `model` tinyint(1) DEFAULT 0,
  `effect` tinyint(1) DEFAULT 0,
  `dubbed` tinyint(1) DEFAULT 0,
  `price` int(11) DEFAULT 0,
  `createtime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;


CREATE TABLE IF NOT EXISTS `made` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` varchar(32) NOT NULL,
  `type` tinyint(1) DEFAULT 0,
  `duration` tinyint(1) DEFAULT 0,
  `city` varchar(32) DEFAULT NULL,
  `company` varchar(255) DEFAULT NULL,
  `createtime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

/* 短信表 */
CREATE TABLE IF NOT EXISTS `sms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` varchar(32) NOT NULL,
  `type` tinyint(1) DEFAULT 0,  /* 1:验证  2:提醒  */
  `ip` varchar(64) DEFAULT NULL,
  `number` int(11) NOT NULL,
  `status` tinyint(1) DEFAULT 0,  /* 1:发送中  2:成功  3:失败  */
  `createtime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;
