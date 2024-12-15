package greeting

type OutputGreeting struct {
	Body struct {
		Message string `doc:"Greeting message" example:"Hello, world!" json:"message"`
	}
}

type InputGreeting struct {
	Name string `doc:"Name to greet" example:"world" maxLength:"30" path:"name"`
}

type PostInputGreeting struct {
	Body struct {
		Name string `doc:"Name to greet" example:"world" json:"name" maxLength:"30" path:"name"`
	}
}
