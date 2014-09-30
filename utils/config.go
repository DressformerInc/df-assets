package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var AppConfig *ConfigScheme

var pathExists map[string]bool

func init() {
	AppConfig = &ConfigScheme{}
}

type ConfigScheme struct {
	App struct {
		ListenOn string `json:"listen_on"`
		HttpsOn  string `json:"https_on"`
		SSLCert  string `json:"ssl_cert"`
		SSLKey   string `json:"ssl_key"`
	} `json:"application"`

	Node struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		StorageRoot string `json:"storage_root"`
	} `json:"node"`

	Connections struct {
		Rethink struct {
			Spec   string `json:"spec"`
			DbName string `json:"db_name"`
		} `json:"rethink"`
	} `json:"connections"`
}

func InitConfigFrom(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("Unable to read", file, "Error:", err)
		return
	}

	err = json.Unmarshal(data, AppConfig)
	if err != nil {
		log.Println("Unable to read config.", err)
	}
}

func (this *ConfigScheme) ListenOn() string {
	return this.App.ListenOn
}

func (this *ConfigScheme) HttpsOn() string {
	return this.App.HttpsOn
}

func (this *ConfigScheme) SSLCert() string {
	return this.App.SSLCert
}

func (this *ConfigScheme) SSLKey() string {
	return this.App.SSLKey
}

func (this *ConfigScheme) NodeId() int {
	return this.Node.Id
}

func (this *ConfigScheme) StorageRoot() string {
	return this.Node.StorageRoot
}

func (this *ConfigScheme) StoragePath(id string) string {
	if len(id) != 24 {
		log.Println(fmt.Sprintf("Wrong id length. %v got, len: %d", id, len(id)))
		return ""
	}

	parts := []string{id[0:2], id[2:4], id[4:6]}

	return this.StorageRoot() + "/" + strings.Join(parts, "/") + "/"
}

func (this *ConfigScheme) StorageFilePath(id string) string {
	return this.StoragePath(id) + id
}

func (this *ConfigScheme) PathCreate(path string) string {
	if !pathExists[path] {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(fmt.Sprintf("Unable to create directory %v, error: %v", path, err))
		}

		pathExists[path] = true
	}

	return path
}

func (this *ConfigScheme) StorageFor(id string) string {
	return this.PathCreate(this.StoragePath(id)) + id
}

func (this *ConfigScheme) RethinkAddress() string {
	return this.Connections.Rethink.Spec
}

func (this *ConfigScheme) RethinkDbName() string {
	return this.Connections.Rethink.DbName
}
