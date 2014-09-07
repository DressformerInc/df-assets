package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var AppConfig *ConfigScheme
var path_exists map[string]bool

func init() {
	AppConfig = &ConfigScheme{}
	path_exists = make(map[string]bool, 0)
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
	if !path_exists[path] {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(fmt.Sprintf("Unable to create directory %v, error: %v", path, err))
		}

		path_exists[path] = true
	}

	return path
}

func (this *ConfigScheme) StorageFor(id string) string {
	return this.PathCreate(this.StoragePath(id)) + id
}
