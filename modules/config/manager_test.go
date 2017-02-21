package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tmpDir := createTmpDir(t, ".", "test-config-files")
	defer os.RemoveAll(tmpDir)

	// ensure configDir is set to the default. other tests could have run before and set it elsewhere
	SetConfigDir("./config")
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
	err = LoadConfig()
	data := GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, data.ServerConfig.rsaKey, "Ensure rsaKey lazily loaded")
	privateKey := data.ServerConfig.RSAKey()
	assert.NotNil(t, data.ServerConfig.rsaKey, "Ensure rsaKey lazily loaded")
	assert.ObjectsAreEqual(privateKey, data.ServerConfig.rsaKey)

	expected := &Config{
		ServerConfig: ServerCfg{
			Name:   "CodeCollaborate",
			Port:   80,
			rsaKey: privateKey, // cheating
		},
		ConnectionConfig: ConnCfgMap{
			"MySQL": ConnCfg{
				Host:     "mysqlHost",
				Port:     3306,
				Username: "user1",
				Password: "pw1",
			},
			"Couchbase": ConnCfg{
				Host:     "couchbaseHost",
				Port:     8092,
				Username: "user2",
				Password: "pw2",
			},
		},
	}

	assert.ObjectsAreEqualValues(expected, data)
}
