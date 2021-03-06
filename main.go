package main

import (
	"fmt"
	"net/http"

	_ "gmsec/routers" // debug模式需要添加[mod]/routers 注册注解路由

	"github.com/gin-gonic/gin"
	"github.com/xxjwxc/ginrpc"
	"github.com/xxjwxc/ginrpc/api"
)

// ReqTest demo struct
type ReqTest struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"user_name" binding:"required"` // 带校验方式
	Password    string `json:"password"`
}

// Hello ...
type Hello struct {
	Index int
}

// Hello 带注解路由(参考beego形式)
// @Router /block [post,get]
func (s *Hello) Hello(c *api.Context, req *ReqTest) {
	fmt.Println(req)
	fmt.Println(s.Index)
	c.JSON(http.StatusOK, "ok")
}

// Hello2 不带注解路由(参数为2默认post)
func (s *Hello) Hello2(c *gin.Context, req ReqTest) {
	fmt.Println(req)
	fmt.Println(s.Index)
	c.JSON(http.StatusOK, "ok")
}

//Hello3 带自定义context跟已解析的req参数回调方式,err,resp 返回模式
func (s *Hello) Hello3(c *gin.Context, req ReqTest) (*ReqTest, error) {
	fmt.Println(req)
	//c.JSON(http.StatusOK, req)
	return &req, nil
}

//TestFun6 带自定义context跟已解析的req参数回调方式,err,resp 返回模式
func TestFun6(c *gin.Context, req ReqTest) (*ReqTest, error) {
	fmt.Println(req)
	//c.JSON(http.StatusOK, req)
	return &req, nil
}

func main() {

	// debug test
	// wp := workpool.New(1000)    // Set the maximum number of threads
	// for i := 0; i < 2000; i++ { // Open 20 requests
	// 	wp.Do(func() error {
	// 		orm := core.Dao.GetDBr()
	// 		var ut []model.UserInfoTbl
	// 		err := orm.Table("user_info_tbl").Find(&ut).Error
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 			return err
	// 		}
	// 		fmt.Println(tools.GetJSONStr(ut, true))
	// 		return nil
	// 	})
	// }

	// wp.Wait()
	// fmt.Println("down")

	// swagger
	// -----end --

	base := ginrpc.New(ginrpc.WithCtx(func(c *gin.Context) interface{} {
		return api.NewCtx(c)
	}), ginrpc.WithDebug(true), ginrpc.WithGroup("xxjwxc"))
	router := gin.Default()

	h := new(Hello)
	h.Index = 123
	base.Register(router, h)                                                     // 对象注册
	router.POST("/test6", base.HandlerFunc(TestFun6))                            // 函数注册
	base.RegisterHandlerFunc(router, []string{"post", "get"}, "/test", TestFun6) // 多种请求方式注册

	router.Run(":8080")
}
