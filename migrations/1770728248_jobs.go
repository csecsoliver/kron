package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("jobs")
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
		    return err
		}
		collection.Fields.Add(&core.RelationField{
			Name: "user_id",
			Required: true,
			CollectionId: usersCollection.Id,
			CascadeDelete: true,
		})
		collection.Fields.Add(&core.TextField{
			Name: "name",
		})
		collection.Fields.Add(&core.URLField{
			Name: "target",
			Required: true,
		})
		collection.Fields.Add(&core.JSONField{
			Name: "request",
		})
		collection.Fields.Add(&core.JSONField{
			Name: "expected_response",
		})
		collection.Fields.Add(&core.JSONField{
			Name: "schedule",
			
		})
		err = app.Save(collection)
		if err != nil {
		    return err
		}

		return nil
	}, func(app core.App) error {
		// add down queries...
		jobsCollection, err := app.FindCollectionByNameOrId("jobs")
		if err != nil {
		    return err
		}
		app.Delete(jobsCollection)
		
		
		return nil
	})
}
