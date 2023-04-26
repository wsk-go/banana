package banana

import (
	"fmt"
	"github.com/JackWSK/banana/defines"
	"github.com/JackWSK/banana/zaplogger"
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
	Logger *zaplogger.Logger `inject:""`
}

func (th *UserController) HelloWorld() Mapping {
	return GetMapping{
		Path: "/hello",
		Handler: func(ctx *fiber.Ctx) error {
			th.Logger.Info("hello world", zap.Any("aaa", "bbb"))
			th.Logger.Error("hello world", zap.Any("aaa", "bbb"))
			return ctx.JSON(fiber.Map{"msg": "success"})
		},
	}
}

func (th *UserController) HelloWorld2() Mapping {
	return GetMapping{
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
	var application = New(Config{
		Engine: engine,
	})

	testBean := &TestBean{}

	err := application.Import(zaplogger.Configuration(zaplogger.LoggerConfig{
		Level:  zapcore.DebugLevel,
		Writer: zaplogger.NewFileWriter("logger.default"),
		LevelWriter: map[zapcore.Level]io.Writer{
			zapcore.InfoLevel: zaplogger.NewFileWriter("logger.info"),
		},
	}))

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterBean(&defines.Bean{
		Value: testBean,
	}, &defines.Bean{
		Value: &User{Name: "user"},
	}, &defines.Bean{
		Value: &User{Name: "user2"},
		Name:  "user2",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = application.RegisterController(&defines.Bean{
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
	fmt.Println(reflect.TypeOf((*Mapping)(nil)).Elem().AssignableTo(tt))
}

func TestReflect2(t *testing.T) {
	//u := reflect.ValueOf(&UserController{})
	//m := u.Method(0)
	//
	//fmt.Println(reflect.TypeOf((*server.Mapping)(nil)).Elem().AssignableTo(tt))
}
