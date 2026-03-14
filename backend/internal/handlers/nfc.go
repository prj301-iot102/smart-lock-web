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

type NfcResource struct {
	db  *pgxpool.Pool
	jwt *token.JwtAuth
}

func (nr *NfcResource) GetNfc(c fuego.ContextNoBody) (database.GetTagByIdRow, error) {
	nfc_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return database.GetTagByIdRow{}, fuego.BadRequestError{
			Detail: "Invalid uuid",
		}
	}

	ctx := context.Background()
	queries := database.New(nr.db)
	nfc, err := queries.GetTagById(ctx, nfc_id)
	if err != nil {
		return database.GetTagByIdRow{}, fuego.NotFoundError{
			Detail: "Nfc id not exists",
			Err:    err,
		}
	}

	return nfc, nil
}

func (nr *NfcResource) ActiveNfc(c fuego.ContextNoBody) (string, error) {
	nfc_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Invalid uuid",
		}
	}

	ctx := context.Background()
	queries := database.New(nr.db)
	nfc, err := queries.UpdateTagStatus(ctx, database.UpdateTagStatusParams{
		ID:       nfc_id,
		IsActive: true,
	})
	if err != nil {
		return "", fuego.NotFoundError{
			Detail: "Nfc id not exists",
			Err:    err,
		}
	}

	return nfc.String(), nil
}

type ValidateNfcRequest struct {
	UID       string `json:"uid"`
	DeviceMac string `json:"mac_device"`
}

func (nr *NfcResource) ValidateNfc(c fuego.ContextWithBody[ValidateNfcRequest]) (bool, error) {
	req, err := c.Body()
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid body",
		}
	}

	ctx := context.Background()
	queries := database.New(nr.db)

	device, err := queries.GetDeviceByMac(ctx, req.DeviceMac)
	if err != nil {
		// return false, fuego.NotFoundError{
		// 	Detail: "Device not found",
		// }
		return false, nil
	}

	tag, err := queries.GetTagByUid(ctx, req.UID)
	if err != nil {
		queries.CreateAccessLog(ctx, database.CreateAccessLogParams{
			EmployeeID: uuid.Nil,
			Status:     database.StatusDenied,
		})
		// return false, fuego.NotFoundError{
		// 	Detail: "This tag does not exist " + req.UID,
		// }

		return false, nil
	}

	employee, err := queries.GetEmployeeById(ctx, tag.EmployeeID)
	if err != nil {
		return false, fuego.InternalServerError{}
	}

	door, err := queries.GetDoorByDeviceId(ctx, device.ID)
	if err != nil {
		// return false, fuego.BadRequestError{
		// 	Detail: "This device is not on this door or door do not exists",
		// }

		return false, nil
	}

	_, err = queries.CheckDoorPermissons(ctx, database.CheckDoorPermissonsParams{
		DoorID: door.ID,
		RoleID: employee.RoleID,
	})
	if err != nil {
		queries.CreateAccessLog(ctx, database.CreateAccessLogParams{
			EmployeeID: employee.ID,
			NfcTagID:   tag.ID,
			DoorID:     door.ID,
			Status:     database.StatusDenied,
		})
		return false, fuego.BadRequestError{}
	}

	// if employee.RoleName != door_permisson.RoleName {
	// 	queries.CreateAccessLog(ctx, database.CreateAccessLogParams{
	// 		EmployeeID: employee.ID,
	// 		NfcTagID:   tag.ID,
	// 		DoorID:     door.ID,
	// 		Status:     database.StatusDenied,
	// 	})
	// 	return false, fuego.BadRequestError{}
	// }

	queries.CreateAccessLog(ctx, database.CreateAccessLogParams{
		EmployeeID: employee.ID,
		NfcTagID:   tag.ID,
		DoorID:     door.ID,
		Status:     database.StatusGranted,
	})

	return true, nil
}

func (nr *NfcResource) RevokeNfc(c fuego.ContextNoBody) (bool, error) {
	nfc_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Invalid uuid",
		}
	}

	ctx := context.Background()
	queries := database.New(nr.db)
	_, err = queries.UpdateTagStatus(ctx, database.UpdateTagStatusParams{
		IsActive: false,
		ID:       nfc_id,
	})
	if err != nil {
		return false, fuego.BadRequestError{
			Detail: "Nfc is inactive",
		}
	}

	return true, nil
}

func (nr *NfcResource) EnableCreate(c fuego.ContextNoBody) (bool, error) {
	device_id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return false, fuego.BadRequestError{}
	}
	ctx := context.Background()
	queries := database.New(nr.db)
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

func (nr *NfcResource) ListNfcTags(c fuego.ContextNoBody) ([]database.ListNfcTagsRow, error) {
	ctx := context.Background()
	queries := database.New(nr.db)
	tags, err := queries.ListNfcTags(ctx)
	if err != nil {
		return []database.ListNfcTagsRow{}, fuego.InternalServerError{}
	}

	return tags, nil
}

type CreateNfcRequest struct {
	Uid string `json:"uid"`
	Mac string `json:"mac_address"`
}

func (nc *NfcResource) CreateNfc(c fuego.ContextWithBody[CreateNfcRequest]) (string, error) {
	req, err := c.Body()
	if err != nil {
		return "", fuego.BadRequestError{
			Err:    err,
			Detail: "Invalid body",
		}
	}

	queries := database.New(nc.db)
	ctx := context.Background()

	device, err := queries.GetDeviceByMac(ctx, req.Mac)
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "Device Mac not found",
		}
	}

	if !device.CanCreate {
		return "", fuego.BadRequestError{
			Detail: "This device cannot create",
		}
	}
	existNfc, _ := queries.CheckUidExist(ctx, req.Uid)
	if existNfc.Uid != "" {
		return "", fuego.BadRequestError{
			Detail: "NFC already exists",
		}
	}

	new_nfc, err := queries.CreateNfcTag(ctx, req.Uid)
	if err != nil {
		return "", fuego.BadRequestError{
			Detail: "",
			Err:    err,
		}
	}

	queries.UpdateDeviceCanCreate(ctx, database.UpdateDeviceCanCreateParams{
		CanCreate: false,
		ID:        device.ID,
	})

	return new_nfc.String(), nil
}

func NfcRoute(s *fuego.Server, db *pgxpool.Pool, jwt *token.JwtAuth) {
	rs := NfcResource{
		db:  db,
		jwt: jwt,
	}

	authMiddleware := middlewares.NewAuthMiddleware(jwt)

	group := fuego.Group(s, "/api/nfc")

	fuego.Get(group, "/", rs.ListNfcTags,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Get(group, "/{id}", rs.GetNfc,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Patch(group, "/{id}", rs.ActiveNfc,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))
	fuego.Post(group, "/validate", rs.ValidateNfc)
	fuego.Patch(group, "/{id}/revoke", rs.RevokeNfc,
		option.Middleware(authMiddleware.RequireAuthentication),
		option.Header("Authorization", "Bearer token", param.Required()))

	fuego.Post(group, "/create", rs.CreateNfc)
}
