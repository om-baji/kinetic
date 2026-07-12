package shared

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		HandleErr(err)
	}

	return db
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&WorkflowRecord{}, &TaskRecord{}, &TaskDependency{}, &Graph{}, &GraphNode{}, &GraphEdge{}, &LogEntry{})
}
