package repo

import (
	"context"
	"errors"
	"time"

	"github.com/deicod/dysv/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Ensure AddressRepo implements AddressRepository
var _ AddressRepository = (*AddressRepo)(nil)

type AddressRepo struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewAddressRepo(db *mongo.Database, timeout time.Duration) *AddressRepo {
	return &AddressRepo{
		collection: db.Collection("addresses"),
		timeout:    timeout,
	}
}

func (r *AddressRepo) Create(ctx context.Context, addr *model.Address) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if addr.ID == "" {
		addr.ID = bson.NewObjectID().Hex()
	}
	// CreatedAt/UpdatedAt should ideally be set in Service too, but Repo setting them is acceptable provided it's consistent.
	// However, stricter "thin repo" implies just dumping data.
	// But let's leave ID generation and timestamps here for now or move to Service?
	// CartRepo sets ID on insert result.
	// Let's keep ID gen here for simplicity if ID is empty.

	_, err := r.collection.InsertOne(ctx, addr)
	return err
}

func (r *AddressRepo) ListByUserID(ctx context.Context, userID string) ([]model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	var addresses []model.Address
	if err := cursor.All(ctx, &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepo) Get(ctx context.Context, id, userID string) (*model.Address, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var addr model.Address
	err := r.collection.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&addr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &addr, nil
}

func (r *AddressRepo) Update(ctx context.Context, addr *model.Address) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": addr.ID, "user_id": addr.UserID}
	update := bson.M{"$set": addr}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *AddressRepo) Delete(ctx context.Context, id, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	return err
}

func (r *AddressRepo) UnsetDefaults(ctx context.Context, userID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "is_default": true},
		bson.M{"$set": bson.M{"is_default": false}},
	)
	return err
}
