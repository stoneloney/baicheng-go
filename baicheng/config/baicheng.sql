CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `password` varchar(255) NOT NULL,
  `logintime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

/*
adminitor
$2a$10$5LeQNWuUSvS.n8dleeD7pu3FZ3hOOGl/xWSbU47yvHWW6omBiadGe
insert into users set name='adminitor', password='$2a$10$5LeQNWuUSvS.n8dleeD7pu3FZ3hOOGl/xWSbU47yvHWW6omBiadGe';
*/


CREATE TABLE IF NOT EXISTS `channels` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`name` varchar(32) NOT NULL,
	`pid` int(11) DEFAULT 0,
	`status` tinyint(1) DEFAULT 1,
	`weight` int(11) DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;


CREATE TABLE IF NOT EXISTS `articles` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`title` varchar(255) NOT NULL,
	`desc` varchar(255) DEFAULT NULL,
	`content` text NOT NULL,
	`channel` int(11) DEFAULT 0,
	`author` varchar(255) DEFAULT NULL,
	`thumburl` varchar(255) DEFAULT NULL,
	`status` tinyint(1) DEFAULT 1,
	`createtime` datetime DEFAULT NULL,
	`modifytime` datetime DEFAULT NULL, 
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;