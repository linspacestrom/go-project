package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linspacestrom/go-project/internal/domain"
)

func (r *Repository) CreateCity(ctx context.Context, name string, isActive bool) (*domain.City, error) {
	query, args, err := psql.Insert("cities").
		Columns("name", "is_active").
		Values(name, isActive).
		Suffix("RETURNING id, name, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create city query: %w", err)
	}

	var city domain.City
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&city.ID, &city.Name, &city.IsActive, &city.CreatedAt, &city.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create city: %w", err)
	}

	return &city, nil
}

func (r *Repository) UpdateCity(ctx context.Context, cityID uuid.UUID, name *string, isActive *bool) (*domain.City, error) {
	builder := psql.Update("cities").Set("updated_at", squirrel.Expr("now()"))
	if name != nil {
		builder = builder.Set("name", *name)
	}
	if isActive != nil {
		builder = builder.Set("is_active", *isActive)
	}

	query, args, err := builder.Where(squirrel.Eq{"id": cityID}).Suffix("RETURNING id, name, is_active, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update city query: %w", err)
	}

	var city domain.City
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&city.ID, &city.Name, &city.IsActive, &city.CreatedAt, &city.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCityNotFound
		}

		return nil, fmt.Errorf("update city: %w", err)
	}

	return &city, nil
}

func (r *Repository) ListCities(ctx context.Context, limit, offset uint64) ([]domain.City, uint64, error) {
	q := psql.Select("id", "name", "is_active", "created_at", "updated_at").
		From("cities").
		OrderBy("name ASC").
		Limit(limit).Offset(offset)
	query, args, err := q.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list cities query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list cities: %w", err)
	}
	defer rows.Close()

	cities := make([]domain.City, 0)
	for rows.Next() {
		var c domain.City
		if scanErr := rows.Scan(&c.ID, &c.Name, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); scanErr != nil {
			return nil, 0, fmt.Errorf("scan city: %w", scanErr)
		}
		cities = append(cities, c)
	}

	countQ, countArgs, err := psql.Select("COUNT(1)").From("cities").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count cities query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cities: %w", err)
	}

	return cities, total, nil
}

func (r *Repository) CreateHub(ctx context.Context, cityID uuid.UUID, name, address string, isActive bool) (*domain.Hub, error) {
	query, args, err := psql.Insert("hubs").
		Columns("city_id", "name", "address", "is_active").
		Values(cityID, name, address, isActive).
		Suffix("RETURNING id, city_id, name, address, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create hub query: %w", err)
	}

	var hub domain.Hub
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&hub.ID, &hub.CityID, &hub.Name, &hub.Address, &hub.IsActive, &hub.CreatedAt, &hub.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create hub: %w", err)
	}

	return &hub, nil
}

func (r *Repository) UpdateHub(ctx context.Context, hubID uuid.UUID, name, address *string, isActive *bool) (*domain.Hub, error) {
	builder := psql.Update("hubs").Set("updated_at", squirrel.Expr("now()"))
	if name != nil {
		builder = builder.Set("name", *name)
	}
	if address != nil {
		builder = builder.Set("address", *address)
	}
	if isActive != nil {
		builder = builder.Set("is_active", *isActive)
	}

	query, args, err := builder.Where(squirrel.Eq{"id": hubID}).Suffix("RETURNING id, city_id, name, address, is_active, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update hub query: %w", err)
	}

	var hub domain.Hub
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&hub.ID, &hub.CityID, &hub.Name, &hub.Address, &hub.IsActive, &hub.CreatedAt, &hub.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrHubNotFound
		}
		return nil, fmt.Errorf("update hub: %w", err)
	}

	return &hub, nil
}

func (r *Repository) ListHubsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Hub, uint64, error) {
	query, args, err := psql.Select("id", "city_id", "name", "address", "is_active", "created_at", "updated_at").
		From("hubs").
		Where(squirrel.Eq{"city_id": cityID}).
		OrderBy("name ASC").
		Limit(limit).Offset(offset).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list hubs query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list hubs: %w", err)
	}
	defer rows.Close()

	hubs := make([]domain.Hub, 0)
	for rows.Next() {
		var h domain.Hub
		if scanErr := rows.Scan(&h.ID, &h.CityID, &h.Name, &h.Address, &h.IsActive, &h.CreatedAt, &h.UpdatedAt); scanErr != nil {
			return nil, 0, fmt.Errorf("scan hub: %w", scanErr)
		}
		hubs = append(hubs, h)
	}

	countQ, countArgs, err := psql.Select("COUNT(1)").From("hubs").Where(squirrel.Eq{"city_id": cityID}).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count hubs query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count hubs: %w", err)
	}

	return hubs, total, nil
}

func (r *Repository) CreateRoom(ctx context.Context, hubID uuid.UUID, name string, description, roomType *string, capacity int, isActive bool) (*domain.Room, error) {
	query, args, err := psql.Insert("rooms").
		Columns("hub_id", "name", "description", "room_type", "capacity", "is_active").
		Values(hubID, name, description, roomType, capacity, isActive).
		Suffix("RETURNING id, hub_id, name, description, room_type, capacity, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create room query: %w", err)
	}

	var room domain.Room
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&room.ID, &room.HubID, &room.Name, &room.Description, &room.RoomType, &room.Capacity, &room.IsActive, &room.CreatedAt, &room.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	return &room, nil
}

func (r *Repository) UpdateRoom(ctx context.Context, roomID uuid.UUID, name, description, roomType *string, capacity *int, isActive *bool) (*domain.Room, error) {
	builder := psql.Update("rooms").Set("updated_at", squirrel.Expr("now()"))
	if name != nil {
		builder = builder.Set("name", *name)
	}
	if description != nil {
		builder = builder.Set("description", *description)
	}
	if roomType != nil {
		builder = builder.Set("room_type", *roomType)
	}
	if capacity != nil {
		builder = builder.Set("capacity", *capacity)
	}
	if isActive != nil {
		builder = builder.Set("is_active", *isActive)
	}

	query, args, err := builder.Where(squirrel.Eq{"id": roomID}).Suffix("RETURNING id, hub_id, name, description, room_type, capacity, is_active, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update room query: %w", err)
	}

	var room domain.Room
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&room.ID, &room.HubID, &room.Name, &room.Description, &room.RoomType, &room.Capacity, &room.IsActive, &room.CreatedAt, &room.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}

		return nil, fmt.Errorf("update room: %w", err)
	}

	return &room, nil
}

func (r *Repository) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	query, args, err := psql.Update("rooms").
		Set("is_active", false).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": roomID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete room query: %w", err)
	}

	tag, err := r.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete room: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRoomNotFound
	}

	return nil
}

func (r *Repository) ListRoomsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Room, uint64, error) {
	query, args, err := psql.Select("r.id", "r.hub_id", "r.name", "r.description", "r.room_type", "r.capacity", "r.is_active", "r.created_at", "r.updated_at").
		From("rooms r").
		Join("hubs h ON h.id = r.hub_id").
		Where(squirrel.Eq{"h.city_id": cityID}).
		OrderBy("r.name ASC").
		Limit(limit).Offset(offset).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list rooms query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list rooms: %w", err)
	}
	defer rows.Close()

	rooms := make([]domain.Room, 0)
	for rows.Next() {
		var room domain.Room
		if scanErr := rows.Scan(&room.ID, &room.HubID, &room.Name, &room.Description, &room.RoomType, &room.Capacity, &room.IsActive, &room.CreatedAt, &room.UpdatedAt); scanErr != nil {
			return nil, 0, fmt.Errorf("scan room: %w", scanErr)
		}
		rooms = append(rooms, room)
	}

	countQ, countArgs, err := psql.Select("COUNT(1)").
		From("rooms r").
		Join("hubs h ON h.id = r.hub_id").
		Where(squirrel.Eq{"h.city_id": cityID}).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count rooms query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count rooms: %w", err)
	}

	return rooms, total, nil
}

func (r *Repository) LockRoomByID(ctx context.Context, roomID uuid.UUID) (*domain.RoomView, error) {
	query, args, err := psql.Select(
		"r.id", "r.hub_id", "r.name", "r.description", "r.room_type", "r.capacity", "r.is_active", "r.created_at", "r.updated_at", "h.city_id", "h.name",
	).From("rooms r").
		Join("hubs h ON h.id = r.hub_id").
		Where(squirrel.Eq{"r.id": roomID}).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build lock room query: %w", err)
	}

	var room domain.RoomView
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&room.ID,
		&room.HubID,
		&room.Name,
		&room.Description,
		&room.RoomType,
		&room.Capacity,
		&room.IsActive,
		&room.CreatedAt,
		&room.UpdatedAt,
		&room.CityID,
		&room.HubName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, fmt.Errorf("lock room: %w", err)
	}

	room.HubIDView = room.HubID
	return &room, nil
}
