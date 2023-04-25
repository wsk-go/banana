package go_nest

import (
	"fmt"
	"github.com/JackWSK/go-nest/nest"
	"github.com/JackWSK/go-nest/nestzap"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_recover "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"reflect"
	"testing"
)

type UserController struct {
	Logger *nestzap.Logger `inject:""`
}

func (th *UserController) HelloWorld() nest.Mapping {
	return nest.GetMapping{
		Path: "/hello",
		Handler: func(ctx *fiber.Ctx) error {
			th.Logger.Info("hello world", zap.Any("aaa", "bbb"))
			th.Logger.Error("hello world", zap.Any("aaa", "bbb"))
			return ctx.JSON(fiber.Map{"msg": "success"})
		},
	}
}

func (th *UserController) HelloWorld2() nest.Mapping {
	return nest.GetMapping{
		Path: "/hello2",
		Handler: func(ctx *fiber.Ctx) error {
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

func TestRegister(t *testing.T) {
	engine := fiber.New(fiber.Config{

		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(gin.H{
				"success": false,
				"code":    0,
				"message": err.Error(),
			})
		},
	})

	engine.Use(cors.New())
	engine.Use(_recover.New())
	engine.Use(func(ctx *fiber.Ctx) (err error) {
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
	var application = nest.New(nest.Config{
		Engine: engine,
	})

	testBean := &TestBean{}

	err := application.Import(nestzap.Module(nestzap.LoggerConfig{
		Level:  zapcore.DebugLevel,
		Writer: nestzap.NewFileWriter("logger.default"),
		LevelWriter: map[zapcore.Level]io.Writer{
			zapcore.InfoLevel: nestzap.NewFileWriter("logger.info"),
		},
	}))

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterBean(&nest.Bean{
		Value: testBean,
	}, &nest.Bean{
		Value: &User{Name: "user"},
	}, &nest.Bean{
		Value: &User{Name: "user2"},
		Name:  "user2",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterController(&nest.Bean{
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
	fmt.Println(reflect.TypeOf((*nest.Mapping)(nil)).Elem().AssignableTo(tt))
}

func TestReflect2(t *testing.T) {
	//u := reflect.ValueOf(&UserController{})
	//m := u.Method(0)
	//
	//fmt.Println(reflect.TypeOf((*server.Mapping)(nil)).Elem().AssignableTo(tt))
}
