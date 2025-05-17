package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

var (
	Client *supabase.Client
)

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_KEY environment variables must be set")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %v", err)
	}

	Client = client
	log.Println("Connected to Supabase!")
	return nil
}

// Helper function to get the Supabase client
func GetClient() *supabase.Client {
	return Client
}

// Database operations
func Insert(table string, data interface{}) error {
	fmt.Print("===============D", data)
	_, _, err := Client.From(table).Insert(data, false, "", "", "").Execute()

	return err
}

func Select(table string, data interface{}, query ...interface{}) error {
	_, _, err := Client.From(table).Select("*", "", false).Execute()
	if err != nil {
		return err
	}
	return nil
}

func Update(table string, data map[string]interface{}, query ...interface{}) error {
	_, _, err := Client.From(table).Update(data, "", "").Execute()
	return err
}

func Delete(table string, query ...interface{}) error {
	_, _, err := Client.From(table).Delete("", "").Execute()
	return err
}

// Auth operations
func SignUp(ctx context.Context, email, password string, data map[string]interface{}) (*types.User, error) {
	req := types.SignupRequest{
		Email:    email,
		Password: password,
		Data:     data,
	}
	fmt.Print("===============D")
	resp, err := Client.Auth.Signup(req)
	if err != nil {
		fmt.Print("===============D", err)
		return nil, err
	}
	return &resp.User, nil
}

func SignIn(ctx context.Context, email, password string) (*types.Session, error) {
	session, err := Client.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func GetUser(ctx context.Context, token string) (*types.User, error) {
	resp, err := Client.Auth.GetUser()
	if err != nil {
		return nil, err
	}
	return &resp.User, nil
}
