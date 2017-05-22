USE pufferpanel;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `_meta` (
  `metaId` INT(11) NOT NULL AUTO_INCREMENT,
  `metaKey` VARCHAR(20) NOT NULL,
  `metaValue` VARCHAR(200) NOT NULL,
  PRIMARY KEY (`metaId`),
  UNIQUE INDEX `UK_metaKey` (`metaKey`)
);

INSERT INTO _meta (metaKey, metaValue) VALUES
  ('version', '1.1.2'),
  ('updateDate', CURRENT_TIMESTAMP)
  ON DUPLICATE KEY UPDATE
  metaKey=VALUES(metaKey),
  metaValue=VALUES(metaValue);

SET FOREIGN_KEY_CHECKS = 1;
