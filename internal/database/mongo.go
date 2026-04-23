package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("mongodb ping failed: %w", err)
	}

	return client, nil
}

func EnsureSchema(ctx context.Context, client *mongo.Client, dbName string) error {
	if client == nil {
		return fmt.Errorf("mongo client is nil")
	}
	if dbName == "" {
		return fmt.Errorf("mongo db name is empty")
	}

	db := client.Database(dbName)
	required := map[string][]mongo.IndexModel{
		"admin_users": {
			{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		"portfolio_content": { // legacy aggregate collection (kept for migration/backward compatibility)
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "status", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "version", Value: -1}},
				Options: options.Index(),
			},
		},
		"portfolio_projects": {
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "status", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "version", Value: -1}},
				Options: options.Index(),
			},
		},
		"portfolio_technical": {
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "status", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "version", Value: -1}},
				Options: options.Index(),
			},
		},
		"portfolio_info": {
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "status", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "version", Value: -1}},
				Options: options.Index(),
			},
		},
		"portfolio_content_history": {
			{
				Keys:    bson.D{{Key: "locale", Value: 1}, {Key: "version", Value: -1}},
				Options: options.Index(),
			},
		},
	}

	for collectionName, indexes := range required {
		if err := ensureCollection(ctx, db, collectionName); err != nil {
			return err
		}
		if len(indexes) > 0 {
			if _, err := db.Collection(collectionName).Indexes().CreateMany(ctx, indexes); err != nil {
				return fmt.Errorf("create indexes for %s: %w", collectionName, err)
			}
		}
	}

	return nil
}

func MigrateLegacyPortfolioContent(ctx context.Context, client *mongo.Client, dbName string) error {
	if client == nil {
		return fmt.Errorf("mongo client is nil")
	}
	if dbName == "" {
		return fmt.Errorf("mongo db name is empty")
	}

	db := client.Database(dbName)
	legacy := db.Collection("portfolio_content")
	projects := db.Collection("portfolio_projects")
	technical := db.Collection("portfolio_technical")
	info := db.Collection("portfolio_info")

	type legacyDoc struct {
		Locale    string    `bson:"locale"`
		Status    string    `bson:"status"`
		Version   int       `bson:"version"`
		Content   bson.M    `bson:"content_json"`
		UpdatedBy string    `bson:"updated_by"`
		UpdatedAt time.Time `bson:"updated_at"`
	}

	cur, err := legacy.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var doc legacyDoc
		if err := cur.Decode(&doc); err != nil {
			return err
		}
		filter := bson.M{"locale": doc.Locale, "status": doc.Status}
		projectsJSON := extractMapValue(doc.Content, "projects", []interface{}{})
		technicalJSON := extractMapValue(doc.Content, "technical", []interface{}{})
		infoJSON := extractMapValue(doc.Content, "portfolioInfo", bson.M{})

		base := bson.M{
			"locale":     doc.Locale,
			"status":     doc.Status,
			"version":    doc.Version,
			"updated_by": doc.UpdatedBy,
			"updated_at": doc.UpdatedAt,
		}

		if _, err := projects.UpdateOne(ctx, filter, bson.M{
			"$set": mergeMaps(base, bson.M{"projects_json": projectsJSON}),
		}, options.Update().SetUpsert(true)); err != nil {
			return err
		}

		if _, err := technical.UpdateOne(ctx, filter, bson.M{
			"$set": mergeMaps(base, bson.M{"technical_json": technicalJSON}),
		}, options.Update().SetUpsert(true)); err != nil {
			return err
		}

		if _, err := info.UpdateOne(ctx, filter, bson.M{
			"$set": mergeMaps(base, bson.M{"info_json": infoJSON}),
		}, options.Update().SetUpsert(true)); err != nil {
			return err
		}
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

func ensureCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	names, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: collectionName}})
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}
	if len(names) > 0 {
		return nil
	}
	if err := db.CreateCollection(ctx, collectionName); err != nil {
		return fmt.Errorf("create collection %s: %w", collectionName, err)
	}
	return nil
}

func extractMapValue(m bson.M, key string, fallback interface{}) interface{} {
	if m == nil {
		return fallback
	}
	v, ok := m[key]
	if !ok {
		return fallback
	}
	return v
}

func mergeMaps(a bson.M, b bson.M) bson.M {
	out := bson.M{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
