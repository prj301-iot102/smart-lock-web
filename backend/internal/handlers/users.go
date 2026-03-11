package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/middlewares"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/utils"
)

type UsersResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

func (ur *UsersResource) GetUser(c fuego.ContextNoBody) (database.GetAccountByIdRow, error) {
	user_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return database.GetAccountByIdRow{}, fuego.BadRequestError{
			Detail: "Wrong uuid format",
			Err:    err,
		}
	}

	queries := database.New(ur.db)
	ctx := context.Background()

	user, err := queries.GetAccountById(ctx, user_id)
	if err != nil {
		return database.GetAccountByIdRow{}, fuego.NotFoundError{
			Detail: "User id not exist",
			Err:    err,
		}
	}

	return user, nil
}

type CreateUserRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Department string `json:"department"`
}

func (ur *UsersResource) CreateUser(c fuego.ContextWithBody[CreateUserRequest]) (string, error) {
	req, err := c.Body()
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Invalid body",
		}
	}
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", fuego.InternalServerError{
			Detail: "Unable to hash password",
			Err:    err,
		}
	}

	ctx := context.Background()
	queries := database.New(ur.db)
	employeeID, err := queries.CreateEmployee(ctx, database.CreateEmployeeParams{
		FullName:   req.FullName,
		Department: req.Department,
	})
	if err != nil {
		return "", fuego.InternalServerError{
			Detail: "Unable to create new employee",
			Err:    err,
		}
	}

	userID, err := queries.CreateUser(ctx, database.CreateUserParams{
		Username:   req.Username,
		Password:   hashPassword,
		EmployeeID: employeeID,
	})
	if err != nil {
		return "", fuego.InternalServerError{
			Detail: "Unable to create new user",
			Err:    err,
		}
	}

	return userID.String(), nil
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

func (ur *UsersResource) UpdateUserPassword(c fuego.ContextWithBody[UpdatePasswordRequest]) (any, error) {
	req, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{
			Detail: "Invalid update body",
		}
	}

	if req.Password == "" {
		return nil, fuego.BadRequestError{
			Detail: "Empty password",
		}
	}

	user_id := c.Value(middlewares.AuthorizationTokenKey).(uuid.UUID)
	queries := database.New(ur.db)
	ctx := context.Background()
	hash_password, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fuego.InternalServerError{}
	}

	if err = queries.UpdatePassword(ctx, database.UpdatePasswordParams{
		Password: hash_password,
		ID:       user_id,
	}); err != nil {
		return nil, fuego.InternalServerError{}
	}

	return nil, nil
}

func UsersRoutes(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := UsersResource{
		db:  db,
		jwt: jwt,
	}
	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/users")
	fuego.Use(group, authMiddleware.RequireAuthentication)

	fuego.Get(group, "/{id}", rs.GetUser)
	fuego.Post(group, "/create", rs.CreateUser)
	fuego.Patch(group, "/update", rs.UpdateUserPassword)
}
