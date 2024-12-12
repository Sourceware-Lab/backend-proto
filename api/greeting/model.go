package greeting

type Output struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type Input struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type PostBodyInput struct {
	Body struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}
}
