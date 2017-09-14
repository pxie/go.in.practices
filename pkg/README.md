# golang包管理

1. go项目GOPATH查找你的项目依赖
```go
// 包的路径名要和文件系统的路径一致
import "go.in.practices/pkg/client"
```
2. 用`package <package_name>`语句定义你的包名，并且在项目中用包名来引用包中相应的方法
```go
// client.go
package lib     // 定义包名为lib

// main.go
import "go.in.practices/pkg/client"
func main() {
	c := lib.NewClient()    // 在main()函数中，通过lib来使用相应的方法
```
3. 小写字符开头的变量或者函数只能在对应的包内使用，所以在main()函数中初始化`lib.client`是非法的
```go
func main() {
    c := &lib.client{}    // 在main()函数中，直接初始化lib.client是非法的
                          // 只能使用lib.NewClient()方法初始化client变量
```
4. 如果返回一个包外可见的接口，则必须实现接口所有的方法。如果，有任何一个接口方法没有实现，则会编译出错。
```go
// client.go
func NewClient() Client {       // 因为，NewClient返回包外可见的接口Client，
}                               // 所以，必须实现Client接口的所有方法，
                                //  IsConnected()
                                //  Connect()
                                //  Disconnect()
                                //  GetID()

type Client interface {
	IsConnected() bool
	Connect() error
	Disconnect(quiesce uint)
	GetID() string
}
```
