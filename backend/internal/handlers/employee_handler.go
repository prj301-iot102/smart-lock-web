package handlers

import (
	"context"
	"math"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prj301-iot102/smart-lock-web/backend/internal/database"
)

type EmployeeResource struct {
	db *pgxpool.Pool
}

type CreateEmployeeRequest struct {
	FullName   string `json:"full_name"  validate:"required"`
	Department string `json:"department" validate:"required"`
	Birth      string `json:"birth"`
}

type UpdateEmployeeRequest struct {
	FullName   string `json:"full_name"`
	Department string `json:"department"`
	Birth      string `json:"birth"`
}

type EmployeeResponse struct {
	ID         uuid.UUID          `json:"id"`
	FullName   string             `json:"full_name"`
	Birth      pgtype.Date        `json:"birth"`
	Department string             `json:"department"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt  pgtype.Timestamptz `json:"updated_at"`
}

type EmployeeListResponse struct {
	Data       []EmployeeResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

type DeleteEmployeeResponse struct {
	Message string `json:"message"`
}

func parseDate(s string) (pgtype.Date, error) {
	if s == "" {
		return pgtype.Date{Valid: false}, nil
	}
	var d pgtype.Date
	if err := d.Scan(s); err != nil {
		return pgtype.Date{}, fuego.BadRequestError{
			Detail: "invalid date format, expected YYYY-MM-DD: " + s,
		}
	}
	return d, nil
}

func textOrNull(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func clampLimit(n int) int {
	if n < 5 {
		return 5
	}
	if n > 10 {
		return 10
	}
	return n
}

func getPagination(c fuego.ContextNoBody) (page, limit, offset int) {
	page = c.QueryParamInt("page")
	if page < 1 {
		page = 1
	}
	limit = clampLimit(c.QueryParamInt("limit"))
	offset = (page - 1) * limit
	return
}

func toResponse(r database.ListEmployeesRow) EmployeeResponse {
	return EmployeeResponse{
		ID:         r.ID,
		FullName:   r.FullName,
		Birth:      r.Birth,
		Department: r.Department,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func buildPage(rows []database.ListEmployeesRow, page, limit int) EmployeeListResponse {
	data := make([]EmployeeResponse, 0, len(rows))
	for _, r := range rows {
		data = append(data, toResponse(r))
	}
	total := int64(len(data))
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	return EmployeeListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

func (er *EmployeeResource) ListEmployees(c fuego.ContextNoBody) (EmployeeListResponse, error) {
	page, limit, offset := getPagination(c)

	rows, err := database.New(er.db).ListEmployees(context.Background(), database.ListEmployeesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return EmployeeListResponse{}, fuego.InternalServerError{Detail: "failed to list employees: " + err.Error()}
	}

	return buildPage(rows, page, limit), nil
}

func (er *EmployeeResource) SearchEmployees(c fuego.ContextNoBody) (EmployeeListResponse, error) {
	name := c.QueryParam("name")
	if name == "" {
		return EmployeeListResponse{}, fuego.BadRequestError{Detail: "query param 'name' is required"}
	}

	page, limit, offset := getPagination(c)

	rows, err := database.New(er.db).ListEmployees(context.Background(), database.ListEmployeesParams{
		Search: textOrNull(name),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return EmployeeListResponse{}, fuego.InternalServerError{Detail: "failed to search employees: " + err.Error()}
	}

	return buildPage(rows, page, limit), nil
}

func (er *EmployeeResource) FilterByBirth(c fuego.ContextNoBody) (EmployeeListResponse, error) {
	page, limit, offset := getPagination(c)

	if c.QueryParam("from") == "" && c.QueryParam("to") == "" {
		return EmployeeListResponse{}, fuego.BadRequestError{Detail: "at least one of 'from' or 'to' is required"}
	}

	birthFrom, err := parseDate(c.QueryParam("from"))
	if err != nil {
		return EmployeeListResponse{}, err
	}
	birthTo, err := parseDate(c.QueryParam("to"))
	if err != nil {
		return EmployeeListResponse{}, err
	}

	rows, err := database.New(er.db).ListEmployees(context.Background(), database.ListEmployeesParams{
		BirthFrom: birthFrom,
		BirthTo:   birthTo,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return EmployeeListResponse{}, fuego.InternalServerError{Detail: "failed to filter by birth: " + err.Error()}
	}

	return buildPage(rows, page, limit), nil
}

func (er *EmployeeResource) FilterByDepartment(c fuego.ContextNoBody) (EmployeeListResponse, error) {
	page, limit, offset := getPagination(c)

	dept := c.QueryParam("department")
	if dept == "" {
		return EmployeeListResponse{}, fuego.BadRequestError{Detail: "query param 'department' is required"}
	}

	rows, err := database.New(er.db).ListEmployees(context.Background(), database.ListEmployeesParams{
		Department: textOrNull(dept),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return EmployeeListResponse{}, fuego.InternalServerError{Detail: "failed to filter by department: " + err.Error()}
	}

	return buildPage(rows, page, limit), nil
}

func (er *EmployeeResource) FilterByDateAdded(c fuego.ContextNoBody) (EmployeeListResponse, error) {
	page, limit, offset := getPagination(c)

	if c.QueryParam("from") == "" && c.QueryParam("to") == "" {
		return EmployeeListResponse{}, fuego.BadRequestError{Detail: "at least one of 'from' or 'to' is required"}
	}

	createdFrom, err := parseDate(c.QueryParam("from"))
	if err != nil {
		return EmployeeListResponse{}, err
	}
	createdTo, err := parseDate(c.QueryParam("to"))
	if err != nil {
		return EmployeeListResponse{}, err
	}

	rows, err := database.New(er.db).ListEmployees(context.Background(), database.ListEmployeesParams{
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return EmployeeListResponse{}, fuego.InternalServerError{Detail: "failed to filter by date added: " + err.Error()}
	}

	return buildPage(rows, page, limit), nil
}

func (er *EmployeeResource) GetEmployee(c fuego.ContextNoBody) (EmployeeResponse, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return EmployeeResponse{}, fuego.BadRequestError{Detail: "invalid employee id"}
	}

	row, err := database.New(er.db).GetEmployeeById(context.Background(), id)
	if err != nil {
		return EmployeeResponse{}, fuego.NotFoundError{Detail: "employee not found"}
	}

	return EmployeeResponse{
		ID:         row.ID,
		FullName:   row.FullName,
		Birth:      row.Birth,
		Department: row.Department,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func (er *EmployeeResource) CreateEmployee(c fuego.ContextWithBody[CreateEmployeeRequest]) (EmployeeResponse, error) {
	req, err := c.Body()
	if err != nil {
		return EmployeeResponse{}, fuego.BadRequestError{Detail: "invalid request body"}
	}

	ctx := context.Background()
	q := database.New(er.db)

	created, err := q.CreateEmployee(ctx, database.CreateEmployeeParams{
		FullName:   req.FullName,
		Department: req.Department,
	})
	if err != nil {
		return EmployeeResponse{}, fuego.InternalServerError{Detail: "failed to create employee: " + err.Error()}
	}

	if req.Birth != "" {
		birth, err := parseDate(req.Birth)
		if err != nil {
			return EmployeeResponse{}, err
		}
		if err := q.UpdateEmployee(ctx, database.UpdateEmployeeParams{
			ID:    created.ID,
			Birth: birth,
		}); err != nil {
			return EmployeeResponse{}, fuego.InternalServerError{Detail: "failed to set birth date: " + err.Error()}
		}
		row, err := q.GetEmployeeById(ctx, created.ID)
		if err != nil {
			return EmployeeResponse{}, fuego.InternalServerError{Detail: "failed to retrieve employee: " + err.Error()}
		}
		return EmployeeResponse{
			ID:         row.ID,
			FullName:   row.FullName,
			Birth:      row.Birth,
			Department: row.Department,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		}, nil
	}

	return EmployeeResponse{
		ID:         created.ID,
		FullName:   created.FullName,
		Birth:      created.Birth,
		Department: created.Department,
		CreatedAt:  created.CreatedAt,
		UpdatedAt:  created.UpdatedAt,
	}, nil
}

func (er *EmployeeResource) UpdateEmployee(c fuego.ContextWithBody[UpdateEmployeeRequest]) (EmployeeResponse, error) {
	req, err := c.Body()
	if err != nil {
		return EmployeeResponse{}, fuego.BadRequestError{Detail: "invalid request body"}
	}

	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return EmployeeResponse{}, fuego.BadRequestError{Detail: "invalid employee id"}
	}

	birth, err := parseDate(req.Birth)
	if err != nil {
		return EmployeeResponse{}, err
	}

	ctx := context.Background()
	q := database.New(er.db)

	if err := q.UpdateEmployee(ctx, database.UpdateEmployeeParams{
		FullName:   textOrNull(req.FullName),
		Birth:      birth,
		Department: textOrNull(req.Department),
		ID:         id,
	}); err != nil {
		return EmployeeResponse{}, fuego.InternalServerError{Detail: "failed to update employee: " + err.Error()}
	}

	row, err := q.GetEmployeeById(ctx, id)
	if err != nil {
		return EmployeeResponse{}, fuego.NotFoundError{Detail: "employee not found"}
	}

	return EmployeeResponse{
		ID:         row.ID,
		FullName:   row.FullName,
		Birth:      row.Birth,
		Department: row.Department,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func (er *EmployeeResource) DeleteEmployee(c fuego.ContextNoBody) (DeleteEmployeeResponse, error) {
	id, err := uuid.Parse(c.PathParam("id"))
	if err != nil {
		return DeleteEmployeeResponse{}, fuego.BadRequestError{Detail: "invalid employee id"}
	}

	if err := database.New(er.db).DeleteEmployee(context.Background(), id); err != nil {
		return DeleteEmployeeResponse{}, fuego.InternalServerError{Detail: "failed to delete employee: " + err.Error()}
	}

	return DeleteEmployeeResponse{Message: "employee deleted successfully"}, nil
}

func EmployeeRoutes(s *fuego.Server, db *pgxpool.Pool) {
	rs := EmployeeResource{db: db}
	g := fuego.Group(s, "/employees")

	fuego.Get(g, "/", rs.ListEmployees,
		fuego.OptionQuery("page", "Page number (default: 1)"),
		fuego.OptionQuery("limit", "Page size 5-10 (default: 10)"),
	)
	fuego.Get(g, "/search", rs.SearchEmployees,
		fuego.OptionQuery("name", "Partial name to search (required)"),
		fuego.OptionQuery("page", "Page number (default: 1)"),
		fuego.OptionQuery("limit", "Page size 5-10 (default: 10)"),
	)
	fuego.Get(g, "/filter/birth", rs.FilterByBirth,
		fuego.OptionQuery("from", "Birth date range start (YYYY-MM-DD)"),
		fuego.OptionQuery("to", "Birth date range end (YYYY-MM-DD)"),
		fuego.OptionQuery("page", "Page number (default: 1)"),
		fuego.OptionQuery("limit", "Page size 5-10 (default: 10)"),
	)
	fuego.Get(g, "/filter/department", rs.FilterByDepartment,
		fuego.OptionQuery("department", "Department name (required)"),
		fuego.OptionQuery("page", "Page number (default: 1)"),
		fuego.OptionQuery("limit", "Page size 5-10 (default: 10)"),
	)
	fuego.Get(g, "/filter/date-added", rs.FilterByDateAdded,
		fuego.OptionQuery("from", "Date added range start (YYYY-MM-DD)"),
		fuego.OptionQuery("to", "Date added range end (YYYY-MM-DD)"),
		fuego.OptionQuery("page", "Page number (default: 1)"),
		fuego.OptionQuery("limit", "Page size 5-10 (default: 10)"),
	)
	fuego.Post(g, "/", rs.CreateEmployee)
	fuego.Get(g, "/{id}", rs.GetEmployee)
	fuego.Put(g, "/{id}", rs.UpdateEmployee)
	fuego.Delete(g, "/{id}", rs.DeleteEmployee)
}
