package handlers

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
)

type NfcResource struct {
	db *pgxpool.Pool
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

func NfcRoute(s *fuego.Server, db *pgxpool.Pool) {
	rs := NfcResource{
		db: db,
	}

	group := fuego.Group(s, "/api/nfc")

	fuego.Get(group, "/{id}", rs.GetNfc)
	fuego.Patch(group, "/{id}/revoke", rs.RevokeNfc)
}
