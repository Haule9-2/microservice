package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/Haule9-2/microservice/adapter/userclient/generatedclient"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    client := generatedclient.NewUserServiceClient(conn)

    for {
        fmt.Println("\nChoose an option:")
        fmt.Println("1. Add User")
        fmt.Println("2. Get User")
        fmt.Println("3. Update User")
        fmt.Println("4. Remove User")
        fmt.Println("5. Exit")

        var choice int
        _, err := fmt.Scanf("%d", &choice)
        if err != nil {
            fmt.Println("Invalid input, please enter a number between 1 and 5.")
            continue
        }

        switch choice {
        case 1:
            var name string
            var age int32
            fmt.Print("Enter user name: ")
            fmt.Scanf("%s", &name)
            fmt.Print("Enter user age: ")
            fmt.Scanf("%d", &age)

            addReq := &generatedclient.AddUserRequest{Name: name, Age: age}
            res, err := client.AddUser(context.Background(), addReq)
            if err != nil {
                log.Fatalf("Error adding user: %v", err)
            }
            fmt.Printf("Added User: %s, Age: %d\n", res.Name, res.Age)

        case 2:
            var userID string
            fmt.Print("Enter user ID: ")
            fmt.Scanf("%s", &userID)

            getReq := &generatedclient.UserRequest{UserId: userID}
            res, err := client.GetUser(context.Background(), getReq)
            if err != nil {
                log.Fatalf("Error getting user: %v", err)
            }
            fmt.Printf("Retrieved User: %s, Age: %d\n", res.Name, res.Age)

        case 3:
            var userID string
            var newName string
            var newAge int32
            fmt.Print("Enter user ID to update: ")
            fmt.Scanf("%s", &userID)
            fmt.Print("Enter new name: ")
            fmt.Scanf("%s", &newName)
            fmt.Print("Enter new age: ")
            fmt.Scanf("%d", &newAge)

            updateReq := &generatedclient.UpdateUserRequest{UserId: userID, Name: newName, Age: newAge}
            res, err := client.UpdateUser(context.Background(), updateReq)
            if err != nil {
                log.Fatalf("Error updating user: %v", err)
            }
            fmt.Printf("Updated User: %s, Age: %d\n", res.Name, res.Age)

        case 4:
            var userID string
            fmt.Print("Enter user ID to remove: ")
            fmt.Scanf("%s", &userID)

            removeReq := &generatedclient.RemoveUserRequest{UserId: userID}
            res, err := client.RemoveUser(context.Background(), removeReq)
            if err != nil {
                log.Fatalf("Error removing user: %v", err)
            }
            fmt.Printf("Removed User: %s, Age: %d\n", res.Name, res.Age)

        case 5:
            fmt.Println("Exiting...")
            os.Exit(0)

        default:
            fmt.Println("Invalid option. Please choose a valid option between 1 and 5.")
        }
    }
}
