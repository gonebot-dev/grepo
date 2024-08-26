<div align="center">
  <a href="https://github.com/gonebot-dev">
    <img width="160" src="/assets/gonebot-logo.png" />
  </a>
  <h1>grepo(<a href="/README_zh.md">中文</a>)</h1>
  <p>🔧 An official simple gonebot plugin manager</p>
</div>

## Grepo

`grepo` is a [gonebot](https://github.com/gonebot-dev/gonebot) plugin manager written in Go.

It is based on the runtime provided by gonebot. By calling some simple methods in this module you can easily install and update gonebot plugins from our [official plugin repository](https://github.com/gonebot-dev/gonebot-plugin-repo)!

## Usage

With a few steps, then you are ready to load plugins via `grepo`!

##### 1. Import `grepo` in your gonebot project
```go
import (
  ...
  
  "github.com/gonebot-dev/grepo"
)
```

After running the command below, `grepo` should be installed to your project.
```sh
go get -u github.com/gonebot-dev/grepo
```

##### 2. `Require` plugins in your main method
```go
func main() {
  grepo.Require("test", "v0.0.1")
  grepo.Require("echo", "latest")

  ...
}
```
When you run your project by `go run` command, `Grepo` will search for your entry file, find versions of the required plugin in official repository and try to install it. And the installed plugins should be ready for you next time you compile and run the project!

## Documentation and Cautions

- ##### The `Require` method
  
    `Require` allows you to install plugins from official repository without looking for it yourself. It provides a convenient approach to install plugins.
    
    But you should be aware that the install process is dumb and it assumes that your entry file has `grouped import`, which looks like this:
    ```go
    import (
      "fmt"
      "os"
      ...
    )
    ```
    Then when you require plugins, `grepo` will add imports automatically in this group, like this:
    ```go
    import (
      _ "github.com/Kingcxp/gonebot-plugin-test"
      "fmt"
      "os"
      ...
    )
    ```
    Then `grepo` will try to run command `go fmt` and `go get -u package` to install the plugin in your project.

- ##### The `SetProxy` method
    By default, `grepo` will find plugins from `https://raw.githubusercontent.com/gonebot-dev/gonebot-plugin-repo/main`, which is the main branch of our official plugin repository. But in case you host a proxy to this repository or created your own repository, you can use `SetProxy` to change this url, let we call it `proxyUrl`, then `grepo` shall find plugins from `proxyUrl/plugins/x/xxx/xxx.json` to load.

- ##### The `SetEntry` method
    By default, `grepo` will search for your entry file (which has a `func main()` segment in path `os.Getwd()`)
    If `grepo` cannot find it, `grepo` will skip all the `Require` method.
    Just in case `grepo` cannot find entry file, you can use `SetEntry` to specify the entry file path, the path given will be automatically converted to absolute path.

- ##### The `Disable` method
    The `Require` method will run every time you run the project and is quite slow (we shall explain why later). So in case you don't want to update plugins every time, you can just call `Disable` method at the beginning of `main` method, like this:
    ```go
    func main() {
      grepo.Disable()
      ...
    }
    ```
    Then all the `Require` methods afterwards should be disabled, which allows you to debug your project extremely faster.

### Why `Require` is so slow?

Every time you run your project, `Require` will reinstall the plugins you required, which leads to the process below:

- Fetch json data from remote repository.
- Parse json data and find `latest` version (for the module name of the plugin) and your required version (to use in `go get -u`).
- Add import in your entry file
- Run `go fmt` to remove repeated imports and then `go get -u` to install the depedency.
- Run `go mod tidy` to remove unused or repeated dependencies.
  
As you can see, this process is extremely slow, so we recommend you to use `Disable` method to disable `Require` method when you need to run your project multiple times in a short period.