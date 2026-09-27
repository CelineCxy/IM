# Day1

## 一、项目代码书写顺序
    1. model
    2. repository
    3. service
    4. handler

## 二、每个文件书写方式
    1. 先定义结构体
        type XXX struct{ }
    2. 再初始化结构体（参数用指针进行传递）
        func NewStruct() *struct

## 三、不同HTTP请求的参数获取方式
| 参数类型 | 示例 | Gin 获取方式 |
|---|---|---|
| 路径参数 | `/users/1` | `c.Param("id")` |
| Query 参数 | `/users?id=1` | `c.Query("id")` |
| JSON 请求体 | `{"name":"Tom"}` | `c.ShouldBindJSON(&req)` |
| Form 表单 | `name=Tom` | `c.PostForm("name")` |
| Header | `Authorization: Bearer xxx` | `c.GetHeader("Authorization")` |
| Cookie | `token=xxx` | `c.Cookie("token")` |

## 四、函数和方法的区别
### 函数
    函数没有接收者（receiver）：
    func Add(a int, b int) int {
      return a + b
    }
    调用时直接使用函数名：
    result := Add(1, 2)
### 方法
    方法会在 func 和方法名之间声明接收者：
    type User struct {
        Name string
    }
    func (u User) SayHello() {
        fmt.Println("你好，我是", u.Name)
    }
    调用时通过对象调用：
    user := User{Name: "Tom"}
    user.SayHello()
    这里的 u User 就是接收者，表示 SayHello 属于 User 类型。
### 主要区别

| 对比 | 函数 | 方法 |
|---|---|---|
| 是否绑定类型 | 否 | 是 |
| 是否有接收者 | 没有 | 有 |
| 调用方式 | `Add(1, 2)` | `user.SayHello()` |
| 适合场景 | 通用功能 | 某个类型的行为 |

## 五、提交项目Git流程

### 1. 初始化当前项目

```powershell
git init
git branch -M main
```

如果项目目录中已经存在 `.git`，就不需要再次执行 `git init`。

### 2. 配置用户名和邮箱

```powershell
git config --global user.name "cxy"
git config --global user.email "cai07181123@gmail.com"
```

检查配置：

```powershell
git config --global --list
```

### 3. 创建 `.gitignore`

提交前应排除依赖、构建文件和敏感信息，例如：

```gitignore
.env
.env.*
node_modules/
__pycache__/
*.pyc
dist/
build/
.vscode/
.idea/
```

不要把密码、Token、私钥等敏感信息提交到 GitHub。

### 4. 创建第一次提交

```powershell
git add .
git status
git commit -m "Initial commit"
```

其中：

- `git add .`：将修改加入暂存区
- `git commit`：创建提交
- `-m`：填写提交说明

### 5. 在 GitHub 创建仓库

登录 GitHub，点击：

`New repository`

建议创建一个空仓库，不要勾选自动生成 README、`.gitignore` 或 License。

创建后复制仓库地址，例如：

```text
https://github.com/你的用户名/你的仓库名.git
```

### 6. 关联 GitHub 并推送

```powershell
git remote add origin https://github.com/你的用户名/你的仓库名.git
git push -u origin main
```

之后再推送时通常只需要：

```powershell
git push
```

如果提示远程仓库已存在，可以查看：

```powershell
git remote -v
```

修改远程地址：

```powershell
git remote set-url origin https://github.com/你的用户名/你的仓库名.git
```

使用 HTTPS 推送时，GitHub 通常需要登录或 Personal Access Token；也可以改用 SSH 地址：

```powershell
git remote set-url origin git@github.com:你的用户名/你的仓库名.git
```

### 常用命令

查看状态：

```powershell
git status
```

查看提交记录：

```powershell
git log --oneline --graph --all
```

提交修改：

```powershell
git add .
git commit -m "描述本次修改"
git push
```

只添加某个文件：

```powershell
git add src/app.js
```

查看文件差异：

```powershell
git diff
```

查看某次提交内容：

```powershell
git show <提交ID>
```

拉取远程最新代码：

```powershell
git pull
```

### 回退提交

#### 已经推送到 GitHub：推荐使用 `revert`

这会创建一个新的提交，用来抵消之前的提交，比较安全：

```powershell
git log --oneline
git revert <提交ID>
git push
```

例如：

```powershell
git revert a1b2c3d
git push
```

#### 只想回退本地最近一次提交，但保留修改

```powershell
git reset --soft HEAD~1
```

#### 回退最近一次提交，并取消暂存，但保留文件修改

```powershell
git reset HEAD~1
```

#### 彻底删除最近一次提交及其文件修改

```powershell
git reset --hard HEAD~1
```

`--hard` 会丢失未保存的修改，请谨慎使用。

如果已经把 reset 后的历史推送到 GitHub，可能需要：

```powershell
git push --force-with-lease
```

这会改写远程历史，使用前确认不会影响其他协作者。

一个常见的完整流程是：

```powershell
cd D:\IM
git init
git branch -M main
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/你的用户名/你的仓库名.git
git push -u origin main
```