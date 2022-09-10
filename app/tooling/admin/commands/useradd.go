package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/dudakovict/social-network/business/core/user"
	em "github.com/dudakovict/social-network/business/data/email"
	"github.com/dudakovict/social-network/business/sys/auth"
	"github.com/dudakovict/social-network/business/sys/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// UserAdd adds new users into the database.
func UserAdd(log *zap.SugaredLogger, cfg database.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial("email-service:50084", grpc.WithInsecure())

	if err != nil {
		return fmt.Errorf("connecting to gRPC: %w", err)
	}

	defer conn.Close()

	ec := em.NewEmailClient(conn)

	core := user.NewCore(log, db, ec)

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	usr, err := core.Create(ctx, nu, time.Now())
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
