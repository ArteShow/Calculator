package internal

import (
	"context"
	"database/sql"
	"net"
	"os"
	"testing"
	"time"

	Database "github.com/ArteShow/Calculator/pkg/Database"
	proto "github.com/ArteShow/Calculator/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var testDBPath = "./test.db"

func setupTestDB(t *testing.T) *sql.DB {
	os.Remove(testDBPath)
	db, err := Database.OpenDatabase(testDBPath)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE calculations (
			userId INTEGER,
			calculation TEXT,
			result REAL,
			id INTEGER
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestCalculationExpression(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer os.Remove(testDBPath)

	oldPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", testDBPath)
	defer os.Setenv("DB_PATH", oldPath)

	result := CalculationExpression(1, "2+2 3*3")
	assert.Contains(t, result, "saved with ID")

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM calculations`).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func startGRPCServer(t *testing.T, srv proto.UserServiceServer) net.Listener {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Second) // wait for server
	return lis
}

func TestGetExpressionByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer os.Remove(testDBPath)

	_, err := db.Exec(`INSERT INTO calculations (userId, calculation, result, id) VALUES (?, ?, ?, ?)`,
		42, "5+5", 10, 1)
	assert.NoError(t, err)

	os.Setenv("DB_PATH", testDBPath)

}

func TestSendUserData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer os.Remove(testDBPath)

	_, err := db.Exec(`INSERT INTO calculations (userId, calculation, result, id) VALUES (?, ?, ?, ?)`,
		7, "1+1", 2, 9)
	assert.NoError(t, err)

	os.Setenv("DB_PATH", testDBPath)

	server := &Server{}
	req := &proto.UserDataRequest{UserId: 7, CustomId: 9}
	resp, err := server.SendUserData(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "1+1")
}
