package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "google.golang.org/grpc"
    "github.com/go-redis/redis/v8"
    "github.com/Haule9-2/microservice/adapter/userclient/generatedclient"
    "google.golang.org/grpc/reflection"
    "google.golang.org/protobuf/proto" // Import the protobuf package
)

// server struct implements the UserServiceServer interface.
type server struct {
    generatedclient.UnimplementedUserServiceServer
    redisClient *redis.Client // Redis client for database interactions
}

// GetUser retrieves a user by user_id from Redis.
func (s *server) GetUser(ctx context.Context, req *generatedclient.UserRequest) (*generatedclient.UserResponse, error) {
    // Fetch user data from Redis
    data, err := s.redisClient.Get(ctx, "user:"+req.UserId).Bytes()
    if err != nil {
        if err == redis.Nil {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %v", err)
    }

    // Deserialize user data
    user := &generatedclient.UserResponse{}
    if err := proto.Unmarshal(data, user); err != nil {
        return nil, fmt.Errorf("failed to unmarshal user: %v", err)
    }

    return user, nil
}

// AddUser adds a new user to Redis.
func (s *server) AddUser(ctx context.Context, req *generatedclient.AddUserRequest) (*generatedclient.UserResponse, error) {
    // Create a new user
    newUser := &generatedclient.UserResponse{
        Name: req.Name,
        Age:  req.Age,
    }

    // Serialize user data
    data, err := proto.Marshal(newUser)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal user: %v", err)
    }

    // Store serialized user data in Redis
    if err := s.redisClient.Set(ctx, "user:"+req.Name, data, 0).Err(); err != nil {
        return nil, fmt.Errorf("failed to add user to Redis: %v", err)
    }

    return newUser, nil
}

// UpdateUser updates an existing user in Redis.
func (s *server) UpdateUser(ctx context.Context, req *generatedclient.UpdateUserRequest) (*generatedclient.UserResponse, error) {
    userID := "user:" + req.UserId

    // Check if user exists in Redis
    if err := s.redisClient.Exists(ctx, userID).Err(); err != nil {
        return nil, fmt.Errorf("failed to check user existence: %v", err)
    }

    // Create an updated user
    updatedUser := &generatedclient.UserResponse{
        Name: req.Name,
        Age:  req.Age,
    }

    // Serialize updated user data
    data, err := proto.Marshal(updatedUser)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal updated user: %v", err)
    }

    // Store updated user data in Redis
    if err := s.redisClient.Set(ctx, userID, data, 0).Err(); err != nil {
        return nil, fmt.Errorf("failed to update user in Redis: %v", err)
    }

    return updatedUser, nil
}

// RemoveUser removes a user by user_id from Redis.
func (s *server) RemoveUser(ctx context.Context, req *generatedclient.RemoveUserRequest) (*generatedclient.UserResponse, error) {
    userID := "user:" + req.UserId

    // Fetch user data before deletion
    user, err := s.GetUser(ctx, &generatedclient.UserRequest{UserId: req.UserId})
    if err != nil {
        return nil, err
    }

    // Delete the user from Redis
    err = s.redisClient.Del(ctx, userID).Err()
    if err != nil {
        return nil, fmt.Errorf("failed to remove user from Redis: %v", err)
    }

    return user, nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Redis client configuration
    redisClient := redis.NewClient(&redis.Options{
        Addr: "redis:6379", // Use the container name for Docker networking
    })
    defer redisClient.Close()

    // Test Redis connection
    err := redisClient.Set(ctx, "test_key", "test_value", 0).Err()
    if err != nil {
        log.Fatalf("Failed to set key: %v", err)
    }

    val, err := redisClient.Get(ctx, "test_key").Result()
    if err != nil {
        log.Fatalf("Failed to get key: %v", err)
    }
    log.Printf("Value for 'test_key': %s", val)

    // Initialize gRPC server
    s := &server{redisClient: redisClient}

    // Set up gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    generatedclient.RegisterUserServiceServer(grpcServer, s)
    reflection.Register(grpcServer)
    log.Println("gRPC server is running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
    