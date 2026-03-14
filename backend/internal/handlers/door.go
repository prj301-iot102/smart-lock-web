package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/middlewares"
	"github.com/prj301-iot102/smart-lock-web/backend/internal/token"
)

type DoorResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

type DoorPermissonRequest struct {
	RoleID string `json:"role_id"`
}

type GetDoorResponse struct {
	database.Door
	Roles []string `json:"roles"`
}

func (dr *DoorResource) ListDoors(c fuego.ContextNoBody) ([]database.Door, error) {
	ctx := context.Background()
	queries := database.New(dr.db)
	doors, err := queries.ListDoors(ctx)
	if err != nil {
		return []database.Door{}, fuego.InternalServerError{
			Detail: "Unable to get doors",
			Err:    err,
		}
	}

	return doors, nil
}

func (dr *DoorResource) GetDoor(c fuego.ContextNoBody) (GetDoorResponse, error) {
	doorID, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return GetDoorResponse{}, fuego.BadRequestError{
			Detail: "Invalid uuid",
		}
	}

	ctx := context.Background()
	queries := database.New(dr.db)
	door, err := queries.GetDoorById(ctx, doorID)
	if err != nil {
		return GetDoorResponse{}, fuego.NotFoundError{
			Detail: "Door id not found",
		}
	}

	roleNames, err := queries.GetDoorPermissonByDoorId(ctx, doorID)
	if err != nil {
		return GetDoorResponse{}, fuego.NotFoundError{
			Detail: "Door id not foud",
			Err:    err,
		}
	}

	res := GetDoorResponse{
		door,
		roleNames,
	}

	return res, nil
}

func (dr *DoorResource) AddDoorPermisson(c fuego.ContextWithBody[DoorPermissonRequest]) (string, error) {
	req, err := c.Body()
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Invalid body",
			Err:    err,
		}
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Invalid uuid",
			Err:    err,
		}
	}
	doorID, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Invalid uuid",
			Err:    err,
		}
	}

	ctx := context.Background()
	queries := database.New(dr.db)
	doorPermID, err := queries.AddDoorPermissionRole(ctx, database.AddDoorPermissionRoleParams{
		DoorID: doorID,
		RoleID: roleID,
	})
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Cannot add role to door",
			Err:    err,
		}
	}

	return doorPermID.String(), nil
}

func (dr *DoorResource) DeleteDoorPermisson(c fuego.ContextWithBody[DoorPermissonRequest]) (bool, error) {
	req, err := c.Body()
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid body",
			Err:    err,
		}
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid uuid " + req.RoleID,
			Err:    err,
		}
	}
	doorID, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid uuid",
			Err:    err,
		}
	}

	ctx := context.Background()
	queries := database.New(dr.db)
	err = queries.DeleteDoorPermissonRole(ctx, database.DeleteDoorPermissonRoleParams{
		DoorID: doorID,
		RoleID: roleID,
	})
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Cannot delete role for this door",
			Err:    err,
		}
	}

	return true, nil
}

func DoorRoute(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := DoorResource{
		db:  db,
		jwt: jwt,
	}

	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/door",
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()),
	)

	fuego.Get(group, "/", rs.ListDoors)
	fuego.Get(group, "/{id}", rs.GetDoor)
	fuego.Patch(group, "/{id}", rs.AddDoorPermisson)
	fuego.Put(group, "/{id}", rs.DeleteDoorPermisson)
}
