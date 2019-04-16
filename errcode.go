package base

// ErrCode 错误码
var ErrCode = map[string]string{
	"10000": "Token Error",
	"40400": "Not Found",
	"40401": "Product Not Found",

	"50000": "Server error",
	"50001": "Directory creation failed",
	"50002": "File copy failed",
	"50003": "This Email has submited",
	"50004": "Wrong account or password",

	"50005": "Create Order fail",
	"50006": "User Exists",
	"50007": "Email Exists",
	"50008": "Password mismatch",
	"50009": "Email Not Exists",
	"50010": "Did not pay",
	"50011": "The order has been paid",
	"50012": "Please place a new order",
	"50013": "The order has been expired",
	"50014": "Please sign up first",
	"50015": "Please fill in the number of participants",
	"50016": "Parameter error",
	"50017": "Verification code error",
}
