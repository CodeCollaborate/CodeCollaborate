package dbfs

import "testing"

func TestOpenMySQLConn(t *testing.T)  {
	configSetup()

	my, err := openMySQLConn("cc")
	defer CloseMySQL()

	if err != nil {
		t.Fatal(err)
	}

	err = my.db.Ping()

	if err != nil {
		t.Fatal(err)
	}

}
