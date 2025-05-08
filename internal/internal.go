package internal

import (
	"context"
	"database/sql"
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

// Calculation handler
func CalculationExpression(userId int, expression string) string {
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
	expressionID++

	err = database.InsertData(db, "calculations", map[string]interface{}{
		"userId":     userId,
		"calculation": expression,
		"result":     finalResult,
		"id":         expressionID,
	})

	if err != nil {
		log.Fatalf("❌ Failed to save calculation: %v", err)
	}

	return fmt.Sprintf("Your expression was saved with ID %d", expressionID)
}

// New gRPC method to get calculation by user ID and expression ID
func (s *Server) GetExpressionByID(ctx context.Context, req *user.GetUserCalculationRequest) (*user.UserCalculationResponse, error) {
	// Open the database connection
	dbPath := config.GetDatabasePath() // Replace with your actual database path or method
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Query for the calculation expression by userId and customId (expression ID)
	var expression string
	query := `SELECT calculation FROM calculations WHERE userId = ? AND id = ?`
	err = db.QueryRow(query, req.UserId, req.CustomId).Scan(&expression)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no calculation found for UserId=%d and CustomId=%d", req.UserId, req.CustomId)
		}
		return nil, fmt.Errorf("failed to retrieve calculation: %v", err)
	}

	// Return the calculation expression
	return &user.UserCalculationResponse{
		Expression: expression,
	}, nil
}

// gRPC method that handles user data (expression + user ID)
func (s *Server) SendUserData(ctx context.Context, req *user.UserDataRequest) (*user.UserDataResponse, error) {
    userId := int(req.UserId)
    expressionID := int(req.CustomId)

    // Open the database connection
    dbPath := config.GetDatabasePath()
    db, err := database.OpenDatabase(dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }
    defer db.Close()

    // Query for the calculation expression by userId and customId (expression ID)
    var expression string
    query := `SELECT calculation FROM calculations WHERE userId = ? AND id = ?`
    err = db.QueryRow(query, userId, expressionID).Scan(&expression)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no calculation found for UserId=%d and ExpressionId=%d", userId, expressionID)
        }
        return nil, fmt.Errorf("failed to retrieve calculation: %v", err)
    }

    // Return the calculation expression
    return &user.UserDataResponse{
        Message: fmt.Sprintf("Retrieved expression: %s", expression),
    }, nil
}


// Start gRPC TCP listener
func StartTCPListener() {
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
