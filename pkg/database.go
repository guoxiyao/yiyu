/*package pkg

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

// Connect 连接到数据库并返回 *gorm.DB 和 error
func Connect(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	// 构建数据库连接字符串（DSN）
	dsn := dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 配置 GORM 并打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志模式
	})
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	// 执行数据库迁移
	if err := db.Migrator().AutoMigrate(
		&models.User{},
		&models.Diary{},
		&models.Tag{},      // 迁移 Tag 模型
		&models.DiaryTag{}, // 迁移 DiaryTag 模型
	); err != nil {
		log.Printf("Error migrating the database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully")
	return db, nil

}
*/
package pkg

import (
	"awesomeProject1/config"
	"awesomeProject1/models" // 替换为你的项目路径
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Connect 连接到数据库并执行迁移
func Connect(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	// 构建数据库连接字符串（DSN）
	dsn := dbConfig.FormatDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	// 执行数据库迁移
	if err := models.MigrateAll(db); err != nil {
		log.Printf("Error migrating the database: %v", err)
		return nil, err
	}

	log.Println("Database connection established and migration completed successfully")
	return db, nil
}

// ... 其他函数和逻辑 ...
