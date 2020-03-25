package controllers

type MainController struct {
	BaseController
}

func (c *MainController) Main() {
	c.Success(nil)
}