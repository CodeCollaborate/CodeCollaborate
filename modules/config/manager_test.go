package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tmpDir := createTmpDir(t, ".", "test-config-files")
	defer os.RemoveAll(tmpDir)

	err := LoadConfig()
	if err == nil {
		t.Fatal("Config dir not set yet; ./config does not exist. Should have failed.")
	}

	SetConfigDir(tmpDir)
	err = LoadConfig()
	if err == nil {
		t.Fatal("Empty config dir, should have failed.")
	}

	tmpServerConfigFileName := filepath.Join(tmpDir, "server.cfg")
	serverContent := "{\"Name\": \"CodeCollaborate\",\"Port\": 80}"
	err = ioutil.WriteFile(tmpServerConfigFileName, []byte(serverContent), 0777)
	if err != nil {
		t.Fatal(err)
	}

	tmpConnConfigFileName := filepath.Join(tmpDir, "conn.cfg")
	connContent := "{\"MySQL\": {\"Host\": \"mysqlHost\",\"Port\": 3306,\"Username\": \"user1\",\"Password\": \"pw1\"},\"Couchbase\": {\"Host\": \"couchbaseHost\",\"Port\": 8092,\"Username\": \"user2\",\"Password\": \"pw2\"}}"
	err = ioutil.WriteFile(tmpConnConfigFileName, []byte(connContent), 0777)
	if err != nil {
		t.Fatal(err)
	}

	tmpDatastoreConfigFileName := filepath.Join(tmpDir, "datastore.cfg")
	datastoreContent := "{\"RelationalStoreName\":\"MySQL\",\"RelationalStoreCfg\":{\"Host\":\"mysqlHost\",\"Port\":3306,\"Username\":\"user1\",\"Password\":\"pw1\",\"Timeout\":5,\"NumRetries\":3,\"Schema\":\"cc\"},\"DocumentStoreName\":\"Couchbase\",\"DocumentStoreCfg\":{\"Host\":\"couchbase://couchbaseHost\",\"Port\":8092,\"Username\":\"user2\",\"Password\":\"pw2\",\"Timeout\":5,\"NumRetries\":3,\"Schema\":\"cc\"}}"
	err = ioutil.WriteFile(tmpDatastoreConfigFileName, []byte(datastoreContent), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = LoadConfig()
	data := GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		ServerConfig: &ServerCfg{
			Name: "CodeCollaborate",
			Port: 80,
		},
		DataStoreConfig: &DataStoreCfg{
			RelationalStoreName: "MySQL",
			RelationalStoreCfg: &ConnCfg{
				Host:       "mysqlHost",
				Port:       3306,
				Username:   "user1",
				Password:   "pw1",
				Timeout:    5,
				NumRetries: 3,
				Schema:     "cc",
			},
			DocumentStoreName: "Couchbase",
			DocumentStoreCfg: &ConnCfg{
				Host:       "couchbase://couchbaseHost",
				Port:       8092,
				Username:   "user2",
				Password:   "pw2",
				Timeout:    5,
				NumRetries: 3,
				Schema:     "cc",
			},
		},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Parsed data incorrect. Expected: \n%v\n Actual: \n%v\n", data, expected)
	}
}
