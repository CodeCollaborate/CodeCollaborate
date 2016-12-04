package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	prefix     = flag.String("docker_prefix", "CodeCollaborate", "Prefix for all Docker images and volumes")
	username   = flag.String("username", "root", "Username for all required authentication")
	password   = flag.String("password", "password", "Password for all required authentication")
	schemaName = flag.String("schema_name", "cc", "Name of created schemas for all databases")
)

func volumeName(name string) string {
	return *prefix + "_" + name + "_Data"
}

func imageName(name string) string {
	return *prefix + "_" + name
}

func createAndStartDockerImage(iName string, args ...string) *exec.Cmd {
	args = append([]string{"run", "--name", imageName(iName)}, args...) // prepend docker-command and name to argument set

	c := exec.Command("docker", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = os.Stdin
	c.Run()

	return c
}

func createDockerVolume(vName string) *exec.Cmd {
	c := exec.Command("docker", "volume", "create", "--name", volumeName(vName))
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = os.Stdin
	c.Run()

	return c
}

func main() {
	flag.Parse()

	// Create volumes
	fmt.Println("Creating MySQL Data Volume:")
	createDockerVolume("MySQL")

	fmt.Println("\nCreating MySQL Data Volume:")
	createDockerVolume("Couchbase")

	//Create containers
	//TODO(wongb): Create a new user, and delete root
	fmt.Println("\nStarting MySQL Container:")
	createAndStartDockerImage("MySQL", "-d", "-v", volumeName("MySQL")+":/var/lib/mysql", "-p", "3306:3306", "-e", "MYSQL_ROOT_PASSWORD="+*password, "mysql")

	fmt.Println("\nStarting RabbitMQ Container:")
	//TODO(wongb): Create a new user, and delete guest
	createAndStartDockerImage("RabbitMQ", "-d", "-p", "5672:5672", "-p", "15672:15672", "--hostname", "CodeCollaborate-RabbitMQ", "rabbitmq:management")

	fmt.Println("\nStarting Couchbase Container:")
	createAndStartDockerImage("Couchbase", "-d", "-v", volumeName("Couchbase")+":/opt/couchbase/var/lib/couchbase/data", "-p", "8091-8094:8091-8094", "-p", "11210-11300:11210-11300", "couchbase")

	fmt.Println("\nWaiting for containers to start up:")
	time.Sleep(15 * time.Second)

	fmt.Println("\nSetting up MySQL:")
	f, err := os.Open("../config/defaults/mysql_setup.sql")
	if err != nil {
		f, err = os.Open("config/defaults/mysql_setup.sql")
		if err != nil {
			log.Fatal("Exiting: Failed to open mysql_setup.sql", err)
		}
	}
	c := exec.Command("docker", "exec", "-i", imageName("MySQL"), "mysql", "--protocol=tcp", "-uroot", "-p"+*password)
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = f
	c.Run()

	f, err = os.Open("../config/defaults/mysql_schema_setup.sql")
	if err != nil {
		f, err = os.Open("config/defaults/mysql_schema_setup.sql")
		if err != nil {
			log.Fatal("Exiting: Failed to open mysql_schema_setup.sql", err)
		}
	}
	c = exec.Command("docker", "exec", "-i", imageName("MySQL"), "mysql", "--protocol=tcp", "-uroot", "-p"+*password)
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = f
	c.Run()

	fmt.Println("\nSetting up Couchbase:")
	c = exec.Command("docker", "exec", imageName("Couchbase"), "couchbase-cli", "cluster-init", "-c", "localhost:8091", "--cluster-username="+*username, "--cluster-password="+*password, "--cluster-ramsize=512", "--wait")
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = os.Stdin
	c.Run()
	c = exec.Command("docker", "exec", imageName("Couchbase"), "couchbase-cli", "bucket-create", "-c", "localhost:8091", "-u", *username, "-p", *password, "--bucket="+*schemaName, "--bucket-type=couchbase", "--bucket-password="+*password, "--bucket-ramsize=200", "--enable-index-replica=0", "--bucket-replica=0", "--enable-flush=1", "--wait")
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Stdin = os.Stdin
	c.Run()
}
