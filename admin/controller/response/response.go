package response

// Response 用户响应数据
type RespData struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ShowErrormsg(msg string) RespData {
	return RespData{
		Code:    400,
		Message: msg}

	/*
		msg = lang.Get(msg)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
		})
	*/
}

func ShowErrorParams(msg string) RespData {
	return RespData{
		Code:    400,
		Message: msg}

	/*
		msg = msg + lang.Get("not_exists")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
		})
	*/
}

func ShowSuccess(msg string) RespData {
	return RespData{
		Code:    20000,
		Message: msg}
	/*
		msg = lang.Get(msg)
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  msg,
		})
	*/
}
func ShowData(data interface{}) RespData {
	return RespData{
		Code: 20000,
		Data: data}
	/*
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": data,
		})
	*/
}
