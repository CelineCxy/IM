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
