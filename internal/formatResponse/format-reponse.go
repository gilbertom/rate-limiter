package formatresponse

// ResponseOutput represents the body to reponse http
type ResponseOutput struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

// GetResponse returns the response based on the success flag.
func (m *ResponseOutput) GetResponse(isBlock bool) ResponseOutput {
	
	if isBlock == false {
		return ResponseOutput{
			Code:    200,
			Message: "Success",
		}
	}

	return ResponseOutput{
		Code:    429,
		Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
	}
}
