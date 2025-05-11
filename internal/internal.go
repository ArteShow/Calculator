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

func CalculationExpression(userId int, expression string) string {
	log.Printf("User %d requested: %s", userId, expression)

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
				log.Println("Timeout: Calculation aborted")
				return
			default:
				result, err, code := calculate.Calc(expr)
				mu.Lock()
				if err != nil {
					log.Println("Error in calculation:", err)
				} else {
					log.Printf("Calculation: %s = %f, StatusCode: %d\n", expr, result, code)
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
		log.Fatalf("Failed to open DB: %v", err)
	}
	expressionID, err := database.GetMaxExpressionIdByUserId(db, userId)
	if err != nil {
		log.Fatalf("Failed to get max expression ID: %v", err)
	}
	expressionID++

	err = database.InsertData(db, "calculations", map[string]interface{}{
		"userId":     userId,
		"calculation": expression,
		"result":     finalResult,
		"id":         expressionID,
	})

	if err != nil {
		log.Fatalf("Failed to save calculation: %v", err)
	}

	return fmt.Sprintf("Your expression was saved with ID %d", expressionID)
}

func (s *Server) SendUserData(ctx context.Context, req *user.UserDataRequest) (*user.UserDataResponse, error) {
	userId := int(req.UserId)
	expressionID := int(req.CustomId)

	var expressionInput string
	if req.Calculation != nil {
		expressionInput = req.Calculation.Expression
	}

	// If both are empty/zero, return error
	if expressionID == 0 && expressionInput == "" {
		return &user.UserDataResponse{
			Message: "❌ No expression or ID provided",
		}, nil
	}

	// Case: Calculation input present
	if expressionInput != "" {
		message := CalculationExpression(userId, expressionInput)
		return &user.UserDataResponse{
			Message: message,
		}, nil
	}

	// Case: ID present, fetch from DB
	dbPath := config.GetDatabasePath()
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var expression string
	query := `SELECT calculation FROM calculations WHERE userId = ? AND id = ?`
	err = db.QueryRow(query, userId, expressionID).Scan(&expression)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("❌ No calculation found for UserId=%d and ExpressionId=%d", userId, expressionID)
		}
		return nil, fmt.Errorf("❌ Failed to retrieve calculation: %v", err)
	}

	return &user.UserDataResponse{
		Message: fmt.Sprintf("✅ Retrieved expression: %s", expression),
	}, nil
}


func (s *Server) GetUserCalculations(ctx context.Context, req *user.UserIdRequest) (*user.UserCalculationsResponse, error) {
	userId := int(req.UserId)

	dbPath := config.GetDatabasePath()
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT calculation, result FROM calculations WHERE userId = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query calculations: %v", err)
	}
	defer rows.Close()

	var calculations []*user.Calculation
	for rows.Next() {
		var expression string
		var result float64
		err := rows.Scan(&expression, &result)
		if err != nil {
			continue
		}
		calculations = append(calculations, &user.Calculation{
			Expression: expression,
			Result:     float32(result),
		})
	}

	return &user.UserCalculationsResponse{
		Calculations: calculations,
	}, nil
}

func StartTCPListener() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &Server{})

	log.Println("Server is listening on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
