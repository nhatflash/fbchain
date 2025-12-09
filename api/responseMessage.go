package api

type ApiResponse struct {
	Code		int
	Message 	string
	Data 		any
}

type ErrorResponse struct {
	Code		int
	Error		string
}

