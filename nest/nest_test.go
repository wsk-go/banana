package nest

import (
	"fmt"
	"github.com/JackWSK/go-nest/nestgin"
	"github.com/JackWSK/go-nest/server"
	"github.com/gin-gonic/gin"
	"testing"
)

type UserController struct {
}

func (th *UserController) HelloWorld() server.Mapping {
	return server.GetMapping{
		Path: "/hello",
		Handler: func(ctx server.Context) (body any, err error) {
			return gin.H{
				"msg": "success",
			}, err
		},
	}
}

func (th *UserController) HelloWorld2() server.Mapping {
	return server.GetMapping{
		Path: "/hello2",
		Handler: func(ctx server.Context) (body any, err error) {
			return gin.H{
				"msg": "success",
			}, fmt.Errorf("xxxx")
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
	var nest Nest

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

	nest.useServer(nestgin.New())
	err = nest.Run("0.0.0.0:9222")
	if err != nil {
		t.Fatal(err)
	}

}
