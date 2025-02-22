package route

import (
	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/sing3demons/go-backend-clean-architecture/mongo"
)

func Setup(db mongo.Database, collection string, router bootstrap.IApplication) bootstrap.IApplication {
	NewTaskRoute(db, collection, router)
	return router
}
