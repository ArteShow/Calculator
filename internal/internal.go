package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	calculate "github.com/ArteShow/Calculator/pkg/Calculation"
	config "github.com/ArteShow/Calculator/pkg/Config"
	database "github.com/ArteShow/Calculator/pkg/Database"

	user "github.com/ArteShow/Calculator/proto"
	"google.golang.org/grpc"
)

type Server struct {
	user.UnimplementedUserServiceServer
}

// This is the calculation handler
func CalculationExpression(userId int, expression string) string {
	// You can replace this with real parsing logic
	logs := fmt.Sprintf("User %d requested: %s", userId, expression)
	fmt.Println(logs)
	

	var wg sync.WaitGroup
	var mu sync.Mutex
	var finalResult float64
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parts := strings.Fields(expression)
	wg.Add(len(parts))
	for _, expr := range parts {
		go func(expr string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				log.Println("🛑 Timeout: Calculation aborted")
				return
			default:
				result, err, code := calculate.Calc(expr)
				mu.Lock()
				if err != nil {
					log.Println("❌ Error in calculation:", err)
				} else {
					log.Printf("✅ Calculation: %s = %f, StatusCode: %d\n", expr, result, code)
					finalResult += result
				}
				mu.Unlock()
			}
		}(expr)
	}

	wg.Wait()
	dbPath := config.GetDatabasePath()
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	expressionID, err := database.GetMaxExpressionIdByUserId(db, userId)
	if err != nil {
		log.Fatalf("❌ Failed to get max expression ID: %v", err)
	}
	expressionID++ // Increment the max ID to get the new expression ID

	err = database.InsertData(db, "calculations", map[string]interface{}{
		"userId":     userId,
		"calculation": expression,
		"result":     finalResult,
		"id":         expressionID,
	})

	if err != nil {
		log.Fatalf("❌ Failed to save calculation: %v", err)
	}

	return fmt.Sprintf("Final result: %f", finalResult)
}

// gRPC method that handles incoming requests
func (s *Server) SendUserData(ctx context.Context, req *user.UserDataRequest) (*user.UserDataResponse, error) {
	userId := int(req.UserId)
	expr := req.Calculation.Expression
	var result string
	if expr != ""{
		result = CalculationExpression(userId, expr)
	}
	// Call the calc handler
	

	// Respond to client
	return &user.UserDataResponse{
		Message: result,
	}, nil
}

// Start the TCP gRPC server
func StartTPCListener() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &Server{})

	fmt.Println("Server is listening on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
