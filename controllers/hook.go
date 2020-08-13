package controllers

//检查老师的access token
//func GatewayAccessTeacher(ctx *context.Context) {
//	datas := make(map[string]interface{})
//	F_teacher_id := ctx.Input.Query("F_teacher_id")
//	F_accesstoken := ctx.Input.Query("F_accesstoken")
//	//检查参数
//	if len(F_teacher_id) <= 0 || len(F_accesstoken) <= 0 {
//		ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
//		datas["F_responseNo"] = models.RESP_PARAM_ERR
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//		ctx.Output.JSON(datas, false, false)
//		return
//	}
//	//验证accesstoken
//	respNo := models.RESP_OK
//	var mAuthCurlObj *models.MAuthCurl
//	if strings.Contains(ctx.Request.UserAgent(), "ebagwechatserver") {
//		//网页端登录
//		respNo = mAuthCurlObj.CheckAccessToken(F_teacher_id, F_accesstoken, models.ROLE_TEACHER, models.PLATFORM_WEBCHAT)
//	} else if strings.Contains(ctx.Request.UserAgent(), "Gecko") || strings.Contains(ctx.Request.UserAgent(), "Windows") || strings.Contains(ctx.Request.UserAgent(), "HTML") || ctx.Request.UserAgent() == "" {
//		//网页端登录
//		respNo = mAuthCurlObj.CheckAccessToken(F_teacher_id, F_accesstoken, models.ROLE_TEACHER, models.PLATFORM_WEB)
//	} else {
//		//android登录
//		respNo = mAuthCurlObj.CheckAccessToken(F_teacher_id, F_accesstoken, models.ROLE_TEACHER, models.PLATFORM_ANDROID)
//	}
//	if respNo != models.RESP_OK {
//		//		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
//		datas["F_responseNo"] = respNo
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//		ctx.Output.JSON(datas, false, false)
//		return
//	}
//	return
//}
