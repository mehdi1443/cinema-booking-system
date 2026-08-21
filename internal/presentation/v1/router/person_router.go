package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpPersonRouter(r *gin.RouterGroup, personHandler *handler.PersonHandler) {
	persons := r.Group("/persons")
	{
		persons.POST("", personHandler.CreatePerson)
		persons.GET("", personHandler.GetPersons)
		persons.GET("/:id", personHandler.GetPersonByID)
		persons.PUT("/:id", personHandler.UpdatePerson)
		persons.DELETE("/:id", personHandler.DeletePerson)

		// Image management routes
		persons.POST("/:id/image", personHandler.UploadImage)
		persons.DELETE("/:id/image", personHandler.DeleteImage)
	}
}
