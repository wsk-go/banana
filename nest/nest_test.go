package nest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"testing"
)

type UserController struct {
}

func (th *UserController) HelloWorld() Mapping {
	return GetMapping{
		Path: "/hello",
		Handler: func(ctx *fiber.Ctx) error {
			return fmt.Errorf("xxxx")
			//return ctx.JSON(fiber.Map{"msg": "success"})
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
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(gin.H{
				"success": false,
				"code":    0,
				"message": err.Error(),
			})
		},
	})

	var nest = NewWithConfig(Config{
		app: app,
	})

	testBean := &TestBean{}

	err := nest.RegisterBean(&Bean{
		Value: testBean,
	}, &Bean{
		Value: &User{Name: "user"},
	}, &Bean{
		Value: &User{Name: "user2"},
		Name:  "user2",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = nest.RegisterController(&Bean{
		Value: &UserController{},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = nest.Run("0.0.0.0:9222")
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
