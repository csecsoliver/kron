package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		
		collection := core.NewBaseCollection("status_logs")
		jobsCollection, err := app.FindCollectionByNameOrId("jobs")
		if err != nil {
			return err
		}
		collection.Fields.Add(&core.RelationField{
			Name: "job_id",
			Required: true,
			CollectionId: jobsCollection.Id,
			CascadeDelete: true,
		})
		collection.Fields.Add(&core.BoolField{
			Name: "successful",
		})
		collection.Fields.Add(&core.TextField{
			Name: "error",
		})
		collection.Fields.Add(&core.TextField{
			Name: "status_code",
		})
		collection.Fields.Add(&core.TextField{
			Name: "body",
			Max: 10000000,
		})
		collection.Fields.Add(&core.JSONField{
			Name: "headers",
		})
		collection.Fields.Add(&core.AutodateField{
			Name: "created_at",
			OnCreate: true,
		})
		err = app.Save(collection)
		if err != nil {
			return err
		}
		return nil
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("status_logs")
		if err != nil {
			return err
		}
		app.Delete(collection)
		return nil
	})
}
