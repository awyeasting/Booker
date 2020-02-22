package models

type JoinRequest struct {
	UserId string `json:"userId"`
	BookIds []string `json:"bookIds"`
}

type ResponseResult struct {
	Error		string `json:"error"`
	Result 		string `json:"result"`
}

type BookQueue struct {
	BookId 		string `json:"bookId"`
	UserIds		[]string `json:"userIds"`
}