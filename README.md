### QIM即时通讯(IM)客户端-简单版的类QQ微信

配置文件地址：./config/app.yml

```
mysql:
  dns: "root:xxx@tcp(xxx:xxx)/?charset=utf8mb4&parseTime=True&loc=Local"
redis:
  addr: "xxx:xxx"
  password: "xxx"
  db: 0
port:
  server: ":xxx"
  udp: xxx
cdn:
  path: "xxx"
  url: "xxx"

```

mysql数据库：
```

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `apply`;

CREATE TABLE `apply` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `from_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID [主]',
  `to_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID  [从]',
  `type` int unsigned NOT NULL DEFAULT '0' COMMENT '联系人类型 1用户 2群',
  `reason` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '原因',
  `remark` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户备注',
  `info` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群欢迎语',
  `friend_group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户组ID 0 默认分组',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0默认 1同意 2拒绝',
  `operate_time` int unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='申请联系人表';



# Dump of table contact_friend
# ------------------------------------------------------------

DROP TABLE IF EXISTS `contact_friend`;

CREATE TABLE `contact_friend` (
  `from_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID [主]',
  `to_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID  [从]',
  `friend_group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户组ID 0 默认分组',
  `level` int unsigned NOT NULL DEFAULT '0' COMMENT '用户亲密度',
  `remark` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `phone` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `desc` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `is_top` int unsigned NOT NULL DEFAULT '0' COMMENT '是否置顶 0否1是',
  `is_hidden` int unsigned NOT NULL DEFAULT '0' COMMENT '是否隐藏 0否1是',
  `is_quiet` int unsigned NOT NULL DEFAULT '0' COMMENT '是否免打扰 0否1是',
  `join_time` int unsigned NOT NULL DEFAULT '0' COMMENT '加好友时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  KEY `from_id` (`from_id`),
  KEY `to_id` (`to_id`),
  KEY `friend_group_id` (`friend_group_id`),
  KEY `join_time` (`join_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='联系人表';



# Dump of table contact_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `contact_group`;

CREATE TABLE `contact_group` (
  `from_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID [主]',
  `to_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'ID  [从]',
  `group_power` int unsigned NOT NULL DEFAULT '0' COMMENT '群权限（0 普通 1管理员 2创始人）',
  `level` int unsigned NOT NULL DEFAULT '0' COMMENT '我在本群等级',
  `remark` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群聊备注',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '我在本群昵称',
  `is_top` int unsigned NOT NULL DEFAULT '0' COMMENT '是否置顶 0否1是',
  `is_hidden` int unsigned NOT NULL DEFAULT '0' COMMENT '是否隐藏 0否1是',
  `is_quiet` int unsigned NOT NULL DEFAULT '0' COMMENT '是否免打扰 0否1是',
  `join_time` int unsigned NOT NULL DEFAULT '0' COMMENT '入群时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  KEY `from_id` (`from_id`),
  KEY `to_id` (`to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='联系人表';



# Dump of table devicetoken
# ------------------------------------------------------------

DROP TABLE IF EXISTS `devicetoken`;

CREATE TABLE `devicetoken` (
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT 'UID',
  `token` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '消息token',
  `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型1为ios,2为android',
  `last_login_time` int unsigned NOT NULL DEFAULT '0' COMMENT '最近登陆时间',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备token';



# Dump of table friend_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `friend_group`;

CREATE TABLE `friend_group` (
  `friend_group_id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户组ID',
  `owner_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '拥有者',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户组名',
  `sort` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `is_default` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '默认分组，0否 1是',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`friend_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友组表';



# Dump of table group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `group`;

CREATE TABLE `group` (
  `group_id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `type` int unsigned NOT NULL DEFAULT '0' COMMENT '群类型',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名称',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '图标',
  `info` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `num` int unsigned NOT NULL DEFAULT '0' COMMENT '群人数',
  `exp` int unsigned NOT NULL DEFAULT '0' COMMENT '群经验',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群管理表';



# Dump of table group_tips
# ------------------------------------------------------------

DROP TABLE IF EXISTS `group_tips`;

CREATE TABLE `group_tips` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '群ID',
  `content` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群公告',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群公告';



# Dump of table message
# ------------------------------------------------------------

DROP TABLE IF EXISTS `message`;

CREATE TABLE `message` (
  `id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_id` bigint unsigned DEFAULT '0',
  `to_id` bigint unsigned DEFAULT '0',
  `msg_type` int unsigned DEFAULT '0',
  `msg_media` int unsigned DEFAULT '0',
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `create_time` bigint DEFAULT '0',
  `status` int unsigned DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# Dump of table message_unread
# ------------------------------------------------------------

DROP TABLE IF EXISTS `message_unread`;

CREATE TABLE `message_unread` (
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '目标UID',
  `msg_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '消息',
  `create_time` int NOT NULL DEFAULT '0' COMMENT '创建时间',
  KEY `uid` (`uid`),
  KEY `msg_id` (`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='未读消息表';



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `uid` int unsigned NOT NULL COMMENT 'UID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `sex` int unsigned NOT NULL DEFAULT '0' COMMENT '性别： 0 未知 1男 2女',
  `birthday` int unsigned NOT NULL DEFAULT '0' COMMENT '生日',
  `info` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '简介',
  `exp` int unsigned NOT NULL DEFAULT '0' COMMENT '用户经验',
  `devname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '设备名称',
  `deviceid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '设备ID',
  `reg_time` int unsigned NOT NULL DEFAULT '0' COMMENT '注册时间',
  `login_time` int unsigned NOT NULL DEFAULT '0' COMMENT '最近登录时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';



# Dump of table usermap
# ------------------------------------------------------------

DROP TABLE IF EXISTS `usermap`;

CREATE TABLE `usermap` (
  `uid` int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户UID',
  `siteuid` varchar(64) NOT NULL DEFAULT '' COMMENT '平台UID',
  `sid` int unsigned NOT NULL DEFAULT '0' COMMENT '平台配置id| 0游客、1账号、2google 3fb',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `siteuid` (`siteuid`,`sid`)
) ENGINE=InnoDB DEFAULT DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='玩家平台ID与UID映射表';



# Dump of table usermap_bind
# ------------------------------------------------------------

DROP TABLE IF EXISTS `usermap_bind`;

CREATE TABLE `usermap_bind` (
  `uid` int unsigned NOT NULL COMMENT '用户UID',
  `siteuid` varchar(64) NOT NULL COMMENT '平台UID',
  `sid` int unsigned NOT NULL COMMENT '平台配置id| 0游客、1账号、2google 3fb',
  PRIMARY KEY (`siteuid`,`sid`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='玩家平台ID与UID映射表-绑定';



# Dump of table usermap_device
# ------------------------------------------------------------

DROP TABLE IF EXISTS `usermap_device`;

CREATE TABLE `usermap_device` (
  `deviceid` varchar(50) NOT NULL DEFAULT '' COMMENT '设备号',
  `siteuid` varchar(64) NOT NULL COMMENT '平台UID',
  PRIMARY KEY (`deviceid`),
  UNIQUE KEY `siteuid` (`siteuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备映射表';



# Dump of table usermap_sso
# ------------------------------------------------------------

DROP TABLE IF EXISTS `usermap_sso`;

CREATE TABLE `usermap_sso` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `siteuid` varchar(64) NOT NULL DEFAULT '' COMMENT '平台UID',
  `phone` varchar(32) NOT NULL DEFAULT '' COMMENT '电话号码',
  `email` varchar(32) NOT NULL DEFAULT '' COMMENT '邮箱',
  `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `siteuid` (`siteuid`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `phone` (`phone`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='玩家平台ID与UID映射表-账号';


```
