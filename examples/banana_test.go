package examples

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/wsk-go/banana"
	"github.com/wsk-go/banana/impl/fiberengine"
	"github.com/wsk-go/banana/logger"
	"reflect"
	"testing"
)

type User2 struct {
	Name string
}

type UserRegister struct {
	Logger *logger.Logger `inject:""`
}

func (u *UserRegister) UserService(Logger *logger.Logger) (*User2, string, error) {
	return &User2{}, "", nil
}

func (u *UserRegister) Configuration() banana.ConfigurationFunc {
	return func(application banana.Application) (*banana.Configuration, error) {
		return &banana.Configuration{
			Controllers: nil,
			Beans: []*banana.Bean{
				{
					Value: &User2{Name: "aaa"},
				},
			},
		}, nil
	}
}

type UserController struct {
}

func (th *UserController) HelloWorld() banana.Mapping {
	return banana.GetMapping{
		Path: "/hello",
		Handler: func(ctx banana.Context) error {

			return ctx.JSON(fiber.Map{"msg": "success"})
		},
	}
}

func (th *UserController) HelloWorld2() banana.Mapping {
	return banana.GetMapping{
		Path: "/hello2",
		Handler: func(ctx banana.Context) error {
			return ctx.JSON(fiber.Map{"msg": "success"})
		},
	}
}

type User struct {
	Name string
}

type TestBean struct {
	User  *User `json:"user" inject:""`
	User2 *User `inject:"user2"`
}

func (t *TestBean) Loaded() {
	fmt.Println("loaded")
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func TestRegister(t *testing.T) {
	engine := fiberengine.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(fiber.Map{
				"success": false,
				"code":    0,
				"message": err.Error(),
			})
		},
	})

	engine.App().Use(cors.New())
	engine.App().Use(_recover.New())
	engine.App().Use(func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				_ = ctx.JSON(fiber.Map{
					"success": false,
					"code":    0,
					"message": fmt.Sprintf("%v", r),
				})
			}
		}()

		return ctx.Next()
	})

	engine.App().Get("/swagger/*", swagger.HandlerDefault)
	engine.App().Get("/swagger/*", swagger.New(swagger.Config{ // custom
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "Jack",
			ClientId: "123456",
		},
	}))

	var application = banana.New(banana.Config{
		Engine: engine,
	})

	testBean := &TestBean{}

	err := application.Import()

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterBean(&banana.Bean{
		Value: testBean,
	}, &banana.Bean{
		Value: &User{Name: "user"},
	}, &banana.Bean{
		Value: &User{Name: "user2"},
		Name:  "user2",
	}, &banana.Bean{Value: &UserRegister{}})

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterController(&banana.Bean{
		Value: &UserController{},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = application.Run("0.0.0.0:9222")
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetBean(t *testing.T) {
	engine := fiberengine.New()
	var application = banana.New(banana.Config{
		Engine: engine,
	})

	err := application.Import()

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterController(&banana.Bean{
		Value: &UserController{},
	})

	if err != nil {
		t.Fatal(err)
	}

	c, _ := banana.GetBeanByType[*UserController](application)
	fmt.Println(c)

	cc := banana.MustGetBeanByName[*User](application, "aaa")
	fmt.Println(cc)
	err = application.Run("0.0.0.0:9222")
	if err != nil {
		t.Fatal(err)
	}
}

type MyMiddleware struct {
	UserController *UserController `inject:""`
}

func (th *MyMiddleware) Handle(ctx banana.Context, application *banana.Banana) error {
	fmt.Println(th.UserController)
	return ctx.Next()
}

func TestMiddleware(t *testing.T) {
	e := fiberengine.New()
	var application = banana.New(banana.Config{
		Engine: e,
	})
	var err error

	application.Use(func(ctx banana.Context, application *banana.Banana) error {
		c := banana.MustGetBeanByType[*UserController](application)
		fmt.Println(c)

		return ctx.Next()
	})

	err = application.RegisterMiddleware(&MyMiddleware{})
	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterController(&banana.Bean{
		Value: &UserController{},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = application.Run("0.0.0.0:9222")
	if err != nil {
		t.Fatal(err)
	}
}

func TestReflect(t *testing.T) {
	u := reflect.TypeOf(&UserController{})
	m := u.Method(0)

	tt := m.Type.Out(0)
	fmt.Println(reflect.TypeOf((*banana.Mapping)(nil)).Elem().AssignableTo(tt))
}

func TestReflect2(t *testing.T) {
	//u := reflect.ValueOf(&UserController{})
	//m := u.Method(0)
	//
	//fmt.Println(reflect.TypeOf((*server.Mapping)(nil)).Elem().AssignableTo(tt))
}

//type MyRegister struct {
//	Logger *logger.Logger `inject:"" method:"UserService,UserService2"`
//	Logger *logger.Logger `inject:""`
//	Logger *logger.Logger `inject:""`
//}
//
//func (u *MyRegister) UserService() (*User2, string, error) {
//	return &User2{}, "", nil
//}
//
//func (u *MyRegister) UserService2() (*User2, string, error) {
//	return &User2{}, "", nil
//}

//func TestMyRegister(test *testing.T) {
//	t := reflect.TypeOf(&MyRegister{})
//	for i := 0; i < t.NumMethod(); i++ {
//		method := t.Method(i)
//		for j := 0; j < method.Type.NumIn(); j++ {
//			argType := method.Type.In(j)
//			v := reflect.ValueOf(&MyRegister{})
//			fmt.Println(argType.String())
//		}
//	}
//}

type GenericStruct[T any] struct {
	Data T
}

func (g GenericStruct[T]) GetType() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func TestReflect3(test *testing.T) {
	t := reflect.ValueOf(GenericStruct[User]{})
	tt := reflect.ValueOf(GenericStruct[User]{})
	fmt.Println(t == tt)
	fmt.Println(GenericStruct[User]{}.GetType().String())
}

type aa int64

func TestPointer(test *testing.T) {
	var total *int64 = nil
	var hh = (*aa)(total)
	fmt.Println(hh)
	//total = utils.ToPtr[int64](1)
	//fmt.Println(*total)
}

type MappingTest interface {
	GetData() any
}

type MappingTestImpl[T any] struct {
}

func (th *MappingTestImpl[T]) GetData() any {
	var t T
	return t
}

func TestReflect4(test *testing.T) {
	aa := &MappingTestImpl[User]{}

	print(aa)
}

func print(test any) {

	tt := reflect.TypeOf(test)
	fmt.Println(tt)
	//reflect.TypeOf(MappingTestImpl)
	//tt.AssignableTo()
	if aa, ok := test.(*MappingTestImpl[any]); ok {
		data := aa.GetData()
		fmt.Println(data)
	}
}
