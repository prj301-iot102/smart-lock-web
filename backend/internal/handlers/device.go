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

type DeviceResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

func (dr *DeviceResource) ListDevices(c fuego.ContextNoBody) ([]database.Device, error) {
	ctx := context.Background()
	queries := database.New(dr.db)
	devices, err := queries.ListDevices(ctx)
	if err != nil {
		return []database.Device{}, fuego.InternalServerError{}
	}

	return devices, nil
}

func (dr *DeviceResource) GetDevice(c fuego.ContextNoBody) (database.Device, error) {
	device_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return database.Device{}, fuego.BadRequestError{
			Detail: "Invalid uuid",
		}
	}

	ctx := context.Background()
	queries := database.New(dr.db)
	device, err := queries.GetDeviceById(ctx, device_id)
	if err != nil {
		return database.Device{}, fuego.NotFoundError{
			Detail: "Device id not found",
		}
	}

	return device, nil
}

type DeviceFlagRequest struct {
	Mac string `json:"mac_address"`
}

func (dr *DeviceResource) GetDeviceFlag(c fuego.ContextWithBody[DeviceFlagRequest]) (bool, error) {
	req, err := c.Body()
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid body",
		}
	}

	ctx := context.Background()
	queries := database.New(dr.db)

	device, err := queries.GetDeviceByMac(ctx, req.Mac)
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Device not found",
		}
	}

	return device.CanCreate, nil
}

func (dr *DeviceResource) EnableCreate(c fuego.ContextNoBody) (bool, error) {
	device_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return false, fuego.BadRequestError{}
	}

	ctx := context.Background()
	queries := database.New(dr.db)
	_, err = queries.UpdateDeviceCanCreate(ctx, database.UpdateDeviceCanCreateParams{
		CanCreate: true,
		ID:        device_id,
	})
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invaid body",
		}
	}

	return true, nil
}

func DeviceRoute(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := DeviceResource{
		db:  db,
		jwt: jwt,
	}
	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/devices")
	fuego.Get(group, "/", rs.ListDevices,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Post(group, "/flag", rs.GetDeviceFlag)
	fuego.Get(group, "/{id}", rs.GetDevice,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Patch(group, "/{id}/enable", rs.EnableCreate,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
}
