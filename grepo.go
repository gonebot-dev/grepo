// Package grepo privides official plugin registry support for gonebot.
//
// Update process is quite slow, remember to call Disable first if you don't want it to automatically update.
package grepo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// SetEntry allows you to specify entry file path.
func SetEntry(path string) {
	entry, _ = filepath.Abs(path)
	log.Printf("grepo: [INFO] Specified entry file: %s\n", entry)
}

// Disable will let grepo skip all the update process.
func Disable() {
	entry = ""
}

// Require allows you to add a new plugin to your gonebot from plugin repository.
func Require(name, version string) {
	if entry == "" {
		return
	}
	if version == "" {
		version = "latest"
	}
	log.Printf("grepo: [INFO] Trying to update plugin %s ...\n", name)

	// Get json data from remote repository
	resp, err := http.Get(fmt.Sprintf("%s/plugins/%c/%s/%s.json", baseUrl, name[0], name, name))
	if err != nil {
		log.Printf("grepo: [ERROR] Failed to fetch plugin %s - Network error!\n", name)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("grepo: [ERROR] Failed to fetch plugin %s - %s!\n", name, resp.Status)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !gjson.ValidBytes(body) {
		log.Printf("grepo: [ERROR] Failed to fetch plugin %s - Content is not JSON!\n", name)
		return
	}

	// Find version and latest url from remote data
	dependency := gjson.GetBytes(body, "latest").String()
	if dependency == gjson.Null.String() {
		log.Printf("grepo: [ERROR] Failed to fetch plugin %s - Version latest not found!\n", name)
		return
	}
	dependentVersion := dependency
	if version != "latest" {
		dependentVersion = gjson.GetBytes(body, strings.ReplaceAll(version, ".", `\.`)).String()
		if dependentVersion == gjson.Null.String() {
			log.Printf("grepo: [ERROR] Failed to fetch plugin %s - Version %s not found!\n", name, version)
			return
		}
	}

	// Updating dependency in entry file
	dependencyLine := fmt.Sprintf("\n\t_ \"%s\"", dependency)
	fileContent, _ := os.ReadFile(entry)
	importIndex := bytes.Index(fileContent, []byte("import (")) + 8
	var newContent strings.Builder
	newContent.Write(fileContent[:importIndex])
	newContent.WriteString(dependencyLine)
	newContent.Write(fileContent[importIndex:])
	os.WriteFile(entry, []byte(newContent.String()), 0644)

	// Running format and fetch the version, then tidy go modules
	cmd := exec.Command("go", "get", "-u", dependentVersion)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("grepo: [ERROR] Failed to get plugin version %s!\n", version)
	}
	cmd = exec.Command("go", "fmt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	log.Printf("grepo: [SUCCESS] Plugin %s update comlete!\n", name)
	log.Println("grepo: [INFO] It shall be loaded next time you run your gonebot!")
}

func init() {
	baseUrl = "https://raw.githubusercontent.com/gonebot-dev/gonebot-plugin-repo/main"
	entry = ""

	log.Println("grepo: [INFO] Running go fmt...")
	_, err := exec.LookPath("go")
	if err != nil {
		log.Println("grepo: [ERROR] Cannot find go! Didn't you installed go?")
		log.Println("grepo: [INFO] Grepo should skip plugin update.")
		return
	}
	err = exec.Command("go", "fmt").Run()
	if err != nil {
		log.Println("grepo: [ERROR] Failed to run go fmt!")
		log.Println("grepo: [INFO] Grepo should skip plugin update.")
		return
	}

	log.Println("grepo: [INFO] Looking for entry file ...")
	currentDir, _ := os.Getwd()
	filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Dir(path) != currentDir || info.IsDir() || entry != "" {
			return nil
		}
		fileContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		importIndex := bytes.Index(fileContent, []byte("func main()"))
		if importIndex == -1 {
			return nil
		}
		entry = path
		return nil
	})
	if entry != "" {
		log.Printf("grepo: [INFO] Entry is %s\n", entry)
	} else {
		log.Println("grepo: [ERROR] Cannot find entry file!")
		log.Println("grepo: [INFO] Grepo should skip plugin update.")
	}
}
