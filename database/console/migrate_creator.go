package console

import (
	"os"
	"strings"
	"time"

	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/support/file"
)

type MigrateCreator struct {
	driver string
}

// Create a new migration
func (receiver MigrateCreator) Create(name string) {
	// First we will get the stub file for the migration, which serves as a type
	// of template for the migration. Once we have those we will populate the
	// various place-holders, save the file, and run the post create event.
	table, upStub, downStub := receiver.getStub(name)

	//Create the up.sql file.
	file.Create(receiver.getPath(name, "up"), receiver.populateStub(upStub, table))

	//Create the down.sql file.
	file.Create(receiver.getPath(name, "down"), receiver.populateStub(downStub, table))
}

// getStub Get the migration stub file.
func (receiver MigrateCreator) getStub(name string) (string, string, string) {
	return receiver.smartMigration(name)
}

// populateStub Populate the place-holders in the migration stub.
func (receiver MigrateCreator) populateStub(stub string, table string) string {
	stub = strings.ReplaceAll(stub, "DummyDatabaseCharset", facades.Config.GetString("database.connections."+facades.Config.GetString("database.default")+".charset"))

	if table != "" {
		stub = strings.ReplaceAll(stub, "DummyTable", table)
	}

	return stub
}

// getPath Get the full path to the migration.
func (receiver MigrateCreator) getPath(name string, category string) string {
	pwd, _ := os.Getwd()

	return pwd + "/database/migrations/" + time.Now().Format("20060102150405") + "_" + name + "." + category + ".sql"
}

func (receiver MigrateCreator) smartMigration(migrationName string) (string, string, string) {
	nameParts := strings.Split(migrationName, `_`)
	upQuery := ""
	downQuery := ""
	tableName := ""
	if nameParts[len(nameParts)-1] == "table" {
		switch nameParts[0] {
		case "create":
			if receiver.driver == "postgres" {
				tableName = strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
				createSequence := "CREATE SEQUENCE IF NOT EXISTS " + tableName + "_id_seq;\n"
				upQuery = createSequence + "CREATE TABLE IF NOT EXISTS " + tableName + `
(
	id int8 NOT NULL DEFAULT nextval('` + tableName + `_id_seq'::regclass) PRIMARY KEY, 
	is_active bool default false,
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz
)` + ";"
				dropSequenceQuery := "DROP SEQUENCE IF EXISTS " + tableName + "_seq;\n"
				downQuery = dropSequenceQuery + "DROP TABLE IF EXISTS " + tableName + ";"
			} else if receiver.driver == "mysql" {
				tableName = strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
				upQuery = "CREATE TABLE IF NOT EXISTS " + tableName + `
(
	id BIGINT AUTO_INCREMENT PRIMARY KEY, 
	is_active bool default false,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME Null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime Null
)` + ";"
				downQuery = "DROP TABLE IF EXISTS " + tableName + ";"
			}
			break
		case "drop":
			if receiver.driver == "postgres" {
				tableName = strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
				dropSequenceQuery := "DROP SEQUENCE IF EXISTS " + tableName + "_seq;\n"
				createSequence := "CREATE SEQUENCE IF NOT EXISTS " + tableName + "_id_seq;\n"
				upQuery = dropSequenceQuery + "DROP TABLE IF EXISTS " + tableName + ";"
				downQuery = createSequence + "CREATE TABLE IF NOT EXISTS " + tableName + `
(
	id int8 NOT NULL DEFAULT nextval('` + tableName + `_id_seq'::regclass) PRIMARY KEY, 
	is_active bool default false,
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz
)` + ";"
			} else if receiver.driver == "mysql" {
				tableName = strings.Join(nameParts[1:(len(nameParts)-1)], `_`)
				upQuery = "DROP TABLE IF EXISTS " + tableName + ";"
				downQuery = "CREATE TABLE IF NOT EXISTS " + tableName + `
(
	id BIGINT AUTO_INCREMENT PRIMARY KEY, 
	is_active bool default false,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME Null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime Null
)` + ";"
			}
			break
		case "add":
			for i, part := range nameParts {
				if part == "in" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName = strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " ADD COLUMN " + field + " VARCHAR(200)" + ";"
					downQuery = "ALTER TABLE " + tableName + " DROP COLUMN " + field + ";"
					break
				}
			}
		case "remove":
			for i, part := range nameParts {
				if part == "from" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName = strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " DROP COLUMN " + field + ";"
					downQuery = "ALTER TABLE " + tableName + " ADD COLUMN " + field + " VARCHAR(200)" + ";"
					break
				}
			}
		case "rename":
			for i, part := range nameParts {
				if part == "in" {
					oldTableName := strings.Join(nameParts[1:i], `_`)
					newTableName := strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + oldTableName + " RENAME TO " + newTableName + ";"
					downQuery = "ALTER TABLE " + newTableName + " RENAME TO " + oldTableName + ";"
					break
				}
			}
		case "alter", "change":
			for i, part := range nameParts {
				if part == "in" {
					field := strings.Join(nameParts[1:i], `_`)
					tableName = strings.Join(nameParts[(i+1):(len(nameParts)-1)], `_`)
					upQuery = "ALTER TABLE " + tableName + " ALTER COLUMN " + field + " VARCHAR(200)" + ";"
					downQuery = "ALTER TABLE " + tableName + " ALTER COLUMN " + field + " VARCHAR(200)" + ";"
					break
				}
			}
		}
	}
	return tableName, upQuery, downQuery
}
