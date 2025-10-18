// package repositories

// import (
// 	"context"
// 	"time"

// 	"consistent_1/Domain"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // UserRepository defines methods for interacting with user data.
// type UserRepository interface {
// 	CreateUser(ctx context.Context, user *domain.User) error
// 	GetUserByID(ctx context.Context, id string) (*domain.User, error)
// 	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
// 	UpdateUser(ctx context.Context, user *domain.User) error
// 	GetAllUsers(ctx context.Context) ([]domain.User, error)
// }

// type userRepository struct {
// 	collection *mongo.Collection
// }

// // NewUserRepository creates a new UserRepository.
// func NewUserRepository(db *mongo.Database) UserRepository {
// 	return &userRepository{
// 		collection: db.Collection("users"),
// 	}
// }

// // CreateUser inserts a new user into the database.
// func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
// 	user.ID = primitive.NewObjectID()
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	_, err := r.collection.InsertOne(ctx, user)
// 	return err
// }

// // GetUserByID retrieves a user by their ID.
// func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, domain.ErrUserNotFound // Or a more specific error for invalid ID format
// 	}

// 	var user domain.User
// 	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
// 	if err == mongo.ErrNoDocuments {
// 		return nil, domain.ErrUserNotFound
// 	}
// 	return &user, err
// }

// // GetUserByEmail retrieves a user by their email address.
// func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
// 	var user domain.User
// 	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
// 	if err == mongo.ErrNoDocuments {
// 		return nil, domain.ErrUserNotFound
// 	}
// 	return &user, err
// }

// // UpdateUser updates an existing user's information.
// func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
// 	user.UpdatedAt = time.Now()
// 	_, err := r.collection.UpdateOne(
// 		ctx,
// 		bson.M{"_id": user.ID},
// 		bson.M{"$set": user}, // Overwrite all fields with the updated user struct
// 	)
// 	return err
// }

// // GetAllUsers retrieves all users from the database.
// // Use with caution in production for very large datasets; typically, you'd paginate this.
// func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
// 	var users []domain.User
// 	cursor, err := r.collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	if err = cursor.All(ctx, &users); err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }




package repositories

import (
	"context"
	"time"

	"consistent_1/Domain" // Ensure this path is correct

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	
	UpdateUserLeetCodeStats(ctx context.Context, userID primitive.ObjectID, totalSolved int, lastCheckDate time.Time) error
	
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}


func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	
	user.LeetCodeLastTotalSolved = 0
	user.LeetCodeLastCheckDate = time.Time{} 

	_, err := r.collection.InsertOne(ctx, user)
	return err
}


func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrUserNotFound 
	}

	var user domain.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}


func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}


func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user}, 
	)
	return err
}


func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}


func (r *userRepository) UpdateUserLeetCodeStats(ctx context.Context, userID primitive.ObjectID, totalSolved int, lastCheckDate time.Time) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"leetcodeLastTotalSolved": totalSolved,
			"leetcodeLastCheckDate":   lastCheckDate,
			"updatedAt":               time.Now(), 
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}