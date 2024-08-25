// Package grepo privides official plugin registry support for gonebot.
package grepo

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

var baseUrl string
var entry string

// SetProxy allows you to change base-url of the plugin repository for better performance.
//
// Default base-url "https://raw.githubusercontent.com/gonebot-dev/gonebot-plugin-repo/main"
// includes the path to our official repository and branch name.
//
// Then the url afterwards base-url is string like fmt.Sprintf("/plugins/%c/%s/%s.json", name[0], name, name)
func SetProxy(url string) {
	baseUrl = url
}

// Require allows you to add a new plugin to your gonebot from our official plugin repository.
func Require(name, version string) {
	if version == "" {
		version = "latest"
	}
	log.Printf("grepo: Trying to update plugin %s ...\n", name)

	resp, err := http.Get(fmt.Sprintf("%s/plugins/%c/%s/%s.json", baseUrl, name[0], name, name))
	if err != nil {
		log.Printf("grepo: Failed to fetch plugin %s - Network error!\n", name)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("grepo: Failed to fetch plugin %s - %s!\n", name, resp.Status)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !gjson.ValidBytes(body) {
		log.Printf("grepo: Failed to fetch plugin %s - Content is not JSON!\n", name)
		return
	}

	value := gjson.GetBytes(body, version).String()
	if value == gjson.Null.String() {
		log.Printf("grepo: Failed to fetch plugin %s - Version %s not found!\n", name, version)
		return
	}

	log.Printf("grepo: Plugin %s update comlete!\n", name)
}

func init() {
	baseUrl = "https://raw.githubusercontent.com/gonebot-dev/gonebot-plugin-repo/main"
	entry = ""

	log.Println("grepo: Looking for entry file ...")
}
