package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseServerConfig(t *testing.T) {
	tmpDir := createTmpDir(t, ".", "test-config-files")
	defer os.RemoveAll(tmpDir)

	data, err := parseServerConfig(tmpDir)
	if err == nil {
		t.Fatal("Empty config dir, should have failed.")
	}

	tmpConfigFileName := filepath.Join(tmpDir, "server.cfg")
	content := fmt.Sprint("{\"Name\": \"CodeCollaborate\",\"Port\": 80}")
	err = ioutil.WriteFile(tmpConfigFileName, []byte(content), 0777)
	if err != nil {
		t.Fatal(err)
	}

	data, err = parseServerConfig(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	expected := &ServerCfg{
		Name: "CodeCollaborate",
		Port: 80,
	}

	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Parsed data incorrect. Expected: \n%v\n Actual: \n%v\n", data, expected)
	}
}

func TestParseServerConfigInvalidJSON(t *testing.T) {
	tmpDir := createTmpDir(t, ".", "test-config-files")
	defer os.RemoveAll(tmpDir)

	tmpConfigFileName := filepath.Join(tmpDir, "server.cfg")
	content := fmt.Sprint("{\"InvalidJson\"}")
	err := ioutil.WriteFile(tmpConfigFileName, []byte(content), 0777)
	if err != nil {
		t.Fatal(err)
	}

	_, err = parseServerConfig(tmpDir)
	if err == nil {
		t.Fatal("Invalid JSON input. Should have failed.")
	}
}

func TestParseConfig(t *testing.T) {
	tmpDir := createTmpDir(t, ".", "test-config-files")
	defer os.RemoveAll(tmpDir)

	data, err := parseConfig(tmpDir)
	if err == nil {
		t.Fatal("Empty config dir, should have failed.")
	}

	tmpServerConfigFileName := filepath.Join(tmpDir, "server.cfg")
	serverContent := fmt.Sprint("{\"Name\": \"CodeCollaborate\",\"Port\": 80}")
	err = ioutil.WriteFile(tmpServerConfigFileName, []byte(serverContent), 0777)
	if err != nil {
		t.Fatal(err)
	}

	data, err = parseConfig(tmpDir)
	if err == nil {
		t.Fatal("Conn config file does not exist, should have failed.")
	}

	tmpConnConfigFileName := filepath.Join(tmpDir, "conn.cfg")
	connContent := fmt.Sprint("{\"MySQL\": {\"Host\": \"mysqlHost\",\"Port\": 3306,\"Username\": \"user1\",\"Password\": \"pw1\"},\"Couchbase\": {\"Host\": \"couchbaseHost\",\"Port\": 8092,\"Username\": \"user2\",\"Password\": \"pw2\"}}")
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

	data, err = parseConfig(tmpDir)
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

func createTmpDir(t *testing.T, dir string, prefix string) string {
	relPath, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		t.Fatal(err)
	}
	return relPath

}
