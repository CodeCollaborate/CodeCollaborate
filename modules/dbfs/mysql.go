package dbfs

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	// required to load into local namespace to initialize sql driver mapping
	// in sql.Open("mysql", ...)
	"github.com/CodeCollaborate/Server/modules/config"
	"strconv"
	"github.com/CodeCollaborate/Server/utils"
	"errors"
)

var mysqldb *mysqlConn

type mysqlConn struct {
	config config.ConnCfg
	db *sql.DB
}

func openMySQLConn(dbName string) (*mysqlConn, error) {
	if mysqldb != nil && mysqldb.db != nil {
		err := mysqldb.db.Ping()
		if err == nil {
			return mysqldb, nil
		}
	}

	if (mysqldb == nil || mysqldb.config == (config.ConnCfg{})) {
		mysqldb = new(mysqlConn)
		configMap := config.GetConfig()
		mysqldb.config = configMap.ConnectionConfig["MySQL"]
	}

	db, err := sql.Open("mysql", mysqldb.config.Username + ":" + mysqldb.config.Password + "@tcp(" + mysqldb.config.Host + ":" + strconv.Itoa(int(mysqldb.config.Port)) + ")/" + dbName + "?timeout=" + strconv.Itoa(int(mysqldb.config.Timeout)) + "s")
	if err != nil {
		utils.LogOnError(err, "Unable to connect to MySQL")
		return mysqldb, err
	}

	mysqldb.db = db
	return mysqldb, nil
}

// CloseMySQL closes the MySQL db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func CloseMySQL() error {
	if mysqldb.db != nil {
		err := mysqldb.db.Close()
		mysqldb = nil
		return err
	}
	return errors.New("Bucket not created")
}
