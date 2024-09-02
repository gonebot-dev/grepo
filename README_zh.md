<div align="center">
  <a href="https://github.com/gonebot-dev">
    <img width="160" src="/assets/gonebot-logo.png" />
  </a>
  <h1>grepo(<a href="/README.md">English</a>)</h1>
  <p>🔧 Gonebot 官方简易插件管理器</p>
</div>

## Grepo

`grepo` 是一款使用 Go 语言编写的 [gonebot](https://github.com/gonebot-dev/gonebot)插件管理器

它基于 `gonebot` 提供的运行时但不依赖于 gonebot，通过调用一些简单的方法，你可以轻松地从我们的[官方插件仓库](https://github.com/gonebot-dev/gonebot-plugin-repo)安装和更新 gonebot 插件！

## 用法

通过一些简单的步骤，你就可以使用 `grepo` 来加载插件了！

##### 1. 在你的项目中引用 `grepo`
```go
import (
  ...
  
  "github.com/gonebot-dev/grepo"
)
```

在运行如下命令之后，`grepo` 应该已经安装到你的项目中了。

```sh
go get -u github.com/gonebot-dev/grepo
```

#### 注意！

也许你已经注意到了 `grepo` 需要你使用 `go run` 来运行你的项目，所以如果你想要使用 `go build` 来构建你的项目，你需要在构建之前移除 `grepo`！

##### 2. 在你的主函数中使用 `Require` 方法
```go
func main() {
  grepo.Require("test", "v0.0.1")
  grepo.Require("echo", "latest")

  ...
}
```
当你使用 `go run` 命令运行你的项目时，`grepo` 会自动搜索你的入口文件，查找官方仓库中你需要的插件的版本，并尝试安装。安装完成后，你就可以在下次编译和运行项目时使用这些插件了！

## 文档和注意点

- ##### `Require` 方法
  
    `Require` 允许你从官方仓库安装插件，无需你自己查找相关的插件及其地址！它提供了一条非常方便的途径来安装和使用插件。
    
    不过你需要知道，安装过程非常`愚蠢`，它假设你的入口文件中一定存在多行导入，像这样：
    ```go
    import (
      "fmt"
      "os"
      ...
    )
    ```
    当你调用 `Require` 方法时，`grepo` 会自动在 `import` 语句中添加一行，像这样：
    ```go
    import (
      _ "github.com/Kingcxp/gonebot-plugin-test"
      "fmt"
      "os"
      ...
    )
    ```
    然后它会尝试运行 `go fmt` 和 `go get -u package` 来安装插件。

- ##### `SetProxy` 方法
    默认地，`grepo` 会从 `https://raw.githubusercontent.com/gonebot-dev/gonebot-plugin-repo/main` 这个官方仓库的 main 分支中寻找插件，但是如果你有代理或者自己创建了一个仓库，你可以使用 `SetProxy` 方法来改变这个 url，我们假设它为 `proxyUrl`，然后 `grepo` 会从 `proxyUrl/plugins/x/xxx/xxx.json` 中寻找插件来加载。

- ##### `SetEntry` 方法
    默认地，`grepo` 会从寻找你的入口文件（包含 `func main()` 字符串并处于 `os.Getwd()` 目录下的文件），如果 `grepo` 找不到入口文件，`grepo` 会跳过所有的 `Require` 方法。
    想要避免这种情况，你可以使用 `SetEntry` 方法来指定入口文件路径，给定的路径会被自动转换为绝对路径。

- ##### The `Disable` method
    `Require` 方法在你每次启动项目时都会运行并且非常缓慢（之后会提到）。如果你不想每次都更新插件，你可以在 `main` 方法中调用 `Disable` 方法，像这样：
    ```go
    func main() {
      grepo.Disable()
      ...
    }
    ```
    这样，`Require` 方法之后的所有方法都会被禁用，这允许你更加快速地运行你的项目。

### 为什么 `Require` 这么慢？

每次你使用 `go run` 运行你的项目，`Require` 都会重新安装你需要的插件，这会导致以下过程：

- 从远程仓库获取 json 数据。
- 解析 json 数据并找到 `latest` 版本（用于插件的模块名称）和你的所需版本（用于 `go get -u`）。
- 在你的入口文件中添加导入。
- 运行 `go fmt` 以删除重复的导入，然后运行 `go get -u` 来安装依赖。
- 运行 `go mod tidy` 以删除未使用的或重复的依赖。

Ass you can see♂，这个过程非常的缓慢，所以当你需要短时间内多次运行你的项目时，我们推荐你使用 `Disable` 方法来大大加快项目的启动速度。