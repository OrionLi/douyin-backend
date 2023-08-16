package error

var MsgFlags = map[int]string{}

//GstMag 获取状态码对应信息

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return "error"
	}
	return msg
}
