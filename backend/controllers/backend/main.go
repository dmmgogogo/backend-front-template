package backend

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.Success("Hello World")
}
