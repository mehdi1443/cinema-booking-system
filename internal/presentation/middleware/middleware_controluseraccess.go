package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

// FIXME : Response Format
func ControlUserAccess(
// !permision list
) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exist := c.Get(domainjwt.CtxUserIDKey)
		if !exist {
			panic(fmt.Errorf("WTF!!!!"))
		}
		userID, ok := v.(int64)
		if !ok {
			panic(fmt.Errorf("WTF!!!!"))
		}
		fmt.Println("middleware_controluseraccess:ControlUserAccess:user_id:", userID)
		//FIXME
		// get roles and its permisions
		// check permison list
		/* چیزی که توی ذهن خودمه : توی ردیس ی کلید برای یوزرادی و برای دسترسی هاش بذاریم و توش دسترسی هاشو بذاریم به عنوان کش
		یعنی اگر از توی ردیس هیت نشد بره از دیتابیس بیاره و بعد توی ردیس کشش کنه به هر حال اخرش لیست دسترسی های ی کاربر در میاد
		بعد استافده میشه
		حالا اگر روولش کم و زیاد شد کلا این رکورد توی ردیس پاک میشه
		اگرم دسترسی های ی روول کم زیاد شد باز تمام یوزرادی هایی که این روول رو هم داشتن باید این رکوردشون از ردیس پاک شه
		*/

		c.Next()
	}
}
