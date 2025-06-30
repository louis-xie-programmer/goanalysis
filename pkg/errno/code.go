package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	PARAMERR            = &Errno{Code: 1001, Message: "Param Error"}
	InternalServerError = &Errno{Code: 500, Message: "Internal Server Error"}
	ModelError          = &Errno{Code: 500, Message: "Model Error"}
	ApiServerError      = &Errno{Code: 1002, Message: "Http Api Server Error"}
)
