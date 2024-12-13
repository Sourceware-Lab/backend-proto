package dbexample

type PostOutputDbExample struct {
	Body struct {
		ID string `json:"id" example:"999" doc:"Id for new user"`
	}
}

type PostBodyInputDbExample struct {
	Body struct {
		Name string `path:"name" maxLength:"100" example:"Jo" doc:"Name for new user"`
		Age  uint8  `path:"age" example:"25" doc:"Age for new user"`

		// Optional
		Email        string  `path:"email" maxLength:"100" example:"jo@example.com" doc:"Email for new user" required:"false"`
		Birthday     *string `path:"birthday" example:"2006-01-02" doc:"Birthday for new user" required:"false"`
		MemberNumber *string `path:"member_number" example:"123456" doc:"Member number for new user" required:"false"`
	}
}
