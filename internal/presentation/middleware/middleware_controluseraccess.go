package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func ControlUserAccess(
	rbacRepo DB.RBACInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	requiredPermission string, // !permission list
) gin.HandlerFunc {
	return func(c *gin.Context) {
		//FIXME
		/* چیزی که توی ذهن خودمه : توی ردیس ی کلید برای یوزرادی و برای دسترسی هاش بذاریم و توش دسترسی هاشو بذاریم به عنوان کش
		یعنی اگر از توی ردیس هیت نشد بره از دیتابیس بیاره و بعد توی ردیس کشش کنه به هر حال اخرش لیست دسترسی های ی کاربر در میاد
		بعد استافده میشه
		حالا اگر روولش کم و زیاد شد کلا این رکورد توی ردیس پاک میشه
		اگرم دسترسی های ی روول کم زیاد شد باز تمام یوزرادی هایی که این روول رو هم داشتن باید این رکوردشون از ردیس پاک شه
		*/
		// add roles to ctx ??!?!?!?

		user_id, err := contexTokenService.GetUserIDFromGinContext(c)
		if err != nil {
			log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
			response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
			c.Abort()

			return
		}

		myRealPermissions, err := rbacRepo.GetUserPermissions(c.Request.Context(), user_id)
		if len(myRealPermissions) < 0 || err != nil { //FIXME!!!!!!
			log.Println(err)
			response.Error(c, http.StatusUnauthorized, "You have not enough access permission")
			c.Abort()
			return
		}

		user, err := userRepo.FindUserByID(c.Request.Context(), user_id)
		if !user.IsPhoneVerified || err != nil {
			log.Println(err)
			response.Error(c, http.StatusUnauthorized, "You are not had been Verified your phone number")
			c.Abort()
			return
		}

		myRealPermissionsNames := make([]string, len(myRealPermissions))
		for i := 0; i < len(myRealPermissions); i++ {
			myRealPermissionsNames = append(myRealPermissionsNames, myRealPermissions[i].Name)
		}
		c.Set(domainjwt.CtxPermissionsKey, myRealPermissionsNames)
		log.Println(myRealPermissionsNames)
		c.Next()
	}
}
