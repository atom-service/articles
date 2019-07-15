package provider

import (
	"bytes"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql 驱动
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var (
	createArticleTableStmt  *sqlx.Stmt
	createCategoryTableStmt *sqlx.Stmt
	createLabelTableStmt    *sqlx.Stmt

	insertArticleStmt  *sqlx.NamedStmt // 插入一条文章
	insertLabelStmt    *sqlx.NamedStmt // 插入一条标签
	insertCategoryStmt *sqlx.NamedStmt // 插入一条分类

	updateArticleByIDStmt  *sqlx.NamedStmt // 根据 ID 更新文章
	updateLabelByIDStmt    *sqlx.NamedStmt // 使用 ID 更新 标签
	updateCategoryByIDStmt *sqlx.NamedStmt //  通过 ID 更新分类信息

	countArticleByIDStmt             *sqlx.NamedStmt // 根据 ID 统计文章
	countArticleByOwnerStmt          *sqlx.NamedStmt
	countCategoryByIDStmt            *sqlx.NamedStmt // 使用 ID 统计分类
	countLabelByIDStmt               *sqlx.NamedStmt // 使用 ID 统计 标签
	countArticleByOwnerCategoryStmt  *sqlx.NamedStmt
	countLabelByOwnerStmt            *sqlx.NamedStmt
	countCategoryByOwnerStmt         *sqlx.NamedStmt
	countCategoryByOwnerCategoryStmt *sqlx.NamedStmt

	queryArticleByIDStmt             *sqlx.NamedStmt // 使用 ID 查询文章
	queryArticleByOwnerStmt          *sqlx.NamedStmt
	queryArticleByOwnerCategoryStmt  *sqlx.NamedStmt // 查询分类下的所有文章
	queryLabelByIDStmt               *sqlx.NamedStmt // 使用 ID 查询 标签
	queryLabelByOwnerStmt            *sqlx.NamedStmt // 插入一条标签
	queryCategoryByOwnerStmt         *sqlx.NamedStmt // 通过 Owner 查询
	queryCategoryByOwnerCategoryStmt *sqlx.NamedStmt // 通过 OwnerCategory 查询

	deleteArticleByIDStmt  *sqlx.NamedStmt // 删除文章
	deleteLabelByIDStmt    *sqlx.NamedStmt // 使用 ID 删除 标签
	deleteCategoryByIDStmt *sqlx.NamedStmt //  通过 ID 删除分类信息
)

func init() {
	var err error
	godotenv.Load()

	database, err := sqlx.Connect("mysql", os.Getenv("MYSQL_DB_URI"))
	if err != nil {
		panic(err)
	}

	// 设置 Name 映射方法
	database.MapperFunc(func(field string) string { return field })

	// 创建文章表
	createArticleTableStmt = MustPreparex(
		database,
		"	CREATE TABLE IF NOT EXISTS `articles` (",
		" `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一ID',",
		" `Type` varchar(128) NOT NULL COMMENT '类型',",
		" `Title` varchar(512) NOT NULL COMMENT '标题',",
		" `Owner` int(11) NOT NULL COMMENT '所属作者',",
		" `State` varchar(128) DEFAULT '' COMMENT '状态',",
		" `Cover` varchar(512) DEFAULT '' COMMENT '封面',",
		" `Summary` varchar(1024) DEFAULT '' COMMENT '简介',",
		" `Context` text NOT NULL COMMENT '正文',",
		" `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
		" `UpdateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
		" `OwnerCategory` int(11) DEFAULT '0' COMMENT '所属分类',",
		" PRIMARY KEY (`ID`)",
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
	)
	_, err = createArticleTableStmt.Exec()
	if err != nil {
		panic(err)
	}

	// 创建分类表
	createCategoryTableStmt = MustPreparex(
		database,
		" CREATE TABLE IF NOT EXISTS `categorys` (",
		" `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一ID',",
		" `Type` varchar(128) NOT NULL COMMENT '类型',",
		" `Name` varchar(512) NOT NULL COMMENT '名称',",
		" `Owner` int(11) NOT NULL COMMENT '所属',",
		" `State` varchar(128) DEFAULT '' COMMENT '状态',",
		" `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
		" `UpdateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
		" `OwnerCategory` INT(11) NOT NULL COMMENT ' 所属父类',",
		"  PRIMARY KEY (`ID`)",
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
	)
	_, err = createCategoryTableStmt.Exec()
	if err != nil {
		panic(err)
	}

	// 创建标签表
	createLabelTableStmt = MustPreparex(
		database,
		" CREATE TABLE IF NOT EXISTS `labels` (",
		" `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',",
		" `Type` varchar(128) NOT NULL COMMENT '类型',",
		" `State` varchar(128) DEFAULT '' COMMENT '状态',",
		" `Value` varchar(512) DEFAULT '' COMMENT '值',",
		" `Owner` int(11) NOT NULL COMMENT '所属',",
		" `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
		" `UpdateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
		" PRIMARY KEY (`ID`,`Type`)",
		" ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	)
	_, err = createLabelTableStmt.Exec()
	if err != nil {
		panic(err)
	}

	// 插入一条文章
	insertArticleStmt = MustPreparexNamed(
		database,
		"INSERT INTO `articles`",
		" (`Type`, `Title`,	`Owner`,	`State`,	`Cover`,	`Summary`,	`Context`,	`OwnerCategory`)",
		" VALUES",
		" (:Type,	:Title,	:Owner,	:State,	:Cover,	:Summary,	:Context,	:OwnerCategory)",
		" ;",
	)
	// 通过 ID 查询指定文章
	queryArticleByIDStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `articles`",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 根据 id 统计文章
	countArticleByIDStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `articles`",
		" WHERE `ID` = :ID",
		" ;",
	)

	// 通过 Owner 查询指定文章
	queryArticleByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `articles`",
		" WHERE `Owner` = :Owner",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	)
	// 根据 Owner 统计文章
	countArticleByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `articles`",
		" WHERE `Owner` = :Owner",
		" ;",
	)

	// 删除文章
	deleteArticleByIDStmt = MustPreparexNamed(
		database,
		" DELETE FROM `articles`",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 根据 ID 更新文章
	updateArticleByIDStmt = MustPreparexNamed(
		database,
		" UPDATE `articles` SET",
		" `Type` = :Type,",
		" `Title` = :Title,",
		" `Owner` = :Owner,",
		" `State` = :State,",
		" `Cover` = :Cover,",
		" `Summary` = :Summary,",
		" `Context` = :Context,",
		" `OwnerCategory` = :OwnerCategory",
		"  WHERE `ID` = :ID",
		" ;",
	)

	countArticleByOwnerCategoryStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `articles`",
		" WHERE `OwnerCategory` = :OwnerCategory",
		" ;",
	)

	// 查询分类下的所有文章
	queryArticleByOwnerCategoryStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `articles`",
		" WHERE `OwnerCategory` = :OwnerCategory",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	)
	// 使用 ID 统计分类
	countCategoryByIDStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `categorys`",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 通过 ID 查询标签
	queryLabelByIDStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `labels`",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 通过 ID 更新标签
	queryLabelByIDStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `labels`",
		" WHERE `ID` = :ID",
		" ;",
	)

	countLabelByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `labels`",
		" WHERE `Owner` = :Owner",
		" ;",
	)

	// 通过 Owner 更新标签
	queryLabelByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `labels`",
		" WHERE `Owner` = :Owner",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	)
	// 通过 ID 更新标签
	updateLabelByIDStmt = MustPreparexNamed(
		database,
		" UPDATE `labels` SET",
		" `Type` = :Type,",
		" `State` = :State,",
		" `Value` = :Value,",
		" `Owner` = :Owner",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 使用 ID 统计标签
	countLabelByIDStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `labels`",
		" WHERE `ID` = :ID",
		" ;",
	)
	// 使用 ID 删除标签
	deleteLabelByIDStmt = MustPreparexNamed(
		database,
		" DELETE FROM `labels`",
		" WHERE`ID` = :ID",
		" ;",
	)
	// 插入一条 label
	insertLabelStmt = MustPreparexNamed(
		database,
		" INSERT INTO `labels`",
		" (`Type`, `State`, `Value`, `Owner`)",
		" VALUES",
		" (:Type,	:State,	:Value,	:Owner)",
		" ;",
	)
	// 插入一条分类
	insertCategoryStmt = MustPreparexNamed(
		database,
		" INSERT INTO `categorys`",
		" (`Type`, `Name`, `Owner`, `State`, `OwnerCategory`)",
		" VALUES",
		" (:Type, :Name, :Owner,	:State,	:OwnerCategory)",
		" ;",
	) // 插入一条分类

	// 通过 所属 查询 分类
	countCategoryByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `categorys`",
		" WHERE `Owner` = :Owner",
		" ;",
	) // 通过 Owner 查询

	// 通过 所属 查询 分类
	queryCategoryByOwnerStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `categorys`",
		" WHERE `Owner` = :Owner",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	) // 通过 Owner 查询

	countCategoryByOwnerCategoryStmt = MustPreparexNamed(
		database,
		" SELECT COUNT(*) FROM `categorys`",
		" WHERE `OwnerCategory` = :OwnerCategory",
		" ;",
	) // 通过 OwnerCategory 统计

	// 通过 所属分类 查询 分类 分类的父分类
	queryCategoryByOwnerCategoryStmt = MustPreparexNamed(
		database,
		" SELECT * FROM `categorys`",
		" WHERE `OwnerCategory` = :OwnerCategory",
		" LIMIT :Limit",
		" OFFSET :Offset",
		" ;",
	) // 通过 OwnerCategory 查询
	// 通过 ID 更新分类的信息
	updateCategoryByIDStmt = MustPreparexNamed(
		database,
		" UPDATE `categorys` SET",
		" `Type` = :Type,",
		" `Name` = :Name,",
		" `Owner` = :Owner,",
		" `State` = :State,",
		" `OwnerCategory` = :OwnerCategory ",
		" WHERE `ID` = :ID ;",
	)
	// 使用 ID 删除标签
	deleteCategoryByIDStmt = MustPreparexNamed(
		database,
		" DELETE FROM `categorys`",
		" WHERE `ID` = :ID ;",
	)
}

// MustPreparex 解析 query
func MustPreparex(database *sqlx.DB, querys ...string) *sqlx.Stmt {
	var queryBuf bytes.Buffer

	for _, s := range querys {
		queryBuf.WriteString(s)
	}

	stmp, err := database.Preparex(queryBuf.String())
	if err != nil {
		panic(err)
	}
	return stmp
}

// MustPreparexNamed 解析 query
func MustPreparexNamed(database *sqlx.DB, querys ...string) *sqlx.NamedStmt {
	var queryBuf bytes.Buffer

	for _, s := range querys {
		queryBuf.WriteString(s)
	}

	stmp, err := database.PrepareNamed(queryBuf.String())
	if err != nil {
		panic(err)
	}
	return stmp
}
