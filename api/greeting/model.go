package greeting

type OutputGreeting struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type InputGreeting struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type PostInputGreeting struct {
	Body struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}
}
