package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

// FIXME : Response Format
func ControlUserAccess(
	rbacService service_contract.RBACServiceInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userService service_contract.UserServiceInterface,
	requiredPermission string, // !permission list
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
		/* چیزی که توی ذهن خودمه : توی ردیس ی کلید برای یوزرادی و برای دسترسی هاش بذاریم و توش دسترسی هاشو بذاریم به عنوان کش
		یعنی اگر از توی ردیس هیت نشد بره از دیتابیس بیاره و بعد توی ردیس کشش کنه به هر حال اخرش لیست دسترسی های ی کاربر در میاد
		بعد استافده میشه
		حالا اگر روولش کم و زیاد شد کلا این رکورد توی ردیس پاک میشه
		اگرم دسترسی های ی روول کم زیاد شد باز تمام یوزرادی هایی که این روول رو هم داشتن باید این رکوردشون از ردیس پاک شه
		*/
		// add roles to ctx ??!?!?!?

		tokenData, err := contexTokenService.GetTokenFromGinContext(c)
		if err != nil {
			log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
			response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
			return
		}

		myRealPermissions, err := rbacService.GetUserPermissions(c.Request.Context(), &tokenData, tokenData.UserID)
		if len(myRealPermissions) <= 0 { //FIXME!!!!!!
			response.Error(c, http.StatusUnauthorized, "You have not enough access permission")
		}

		user, err := userService.GetUserByID(c.Request.Context(), &tokenData, tokenData.UserID, &dto.HideFieldUserOption{WalletBalance: true})
		if !user.IsPhoneVerified {
			response.Error(c, http.StatusUnauthorized, "You are not had been Verified your phone number")
		}

		c.Next()
	}
}
