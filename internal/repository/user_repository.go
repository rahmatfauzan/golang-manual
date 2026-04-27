package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/rahmatfauzan/golang-manual/internal/model"
)

type UserRepositoryContract interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUserByID(ctx context.Context, id string) (*model.User, error)
    GetUserByEmail(ctx context.Context, email string) (*model.User, error)
    GetUserByUsername(ctx context.Context, username string) (*model.User, error)
    ListUsers(ctx context.Context) ([]*model.User, error)
    UpdateUser(ctx context.Context, user *model.User) error
    DeleteUser(ctx context.Context, id string) error
    UpdateUserPassword(ctx context.Context, id, passwordHash string) error
    MarkEmailVerified(ctx context.Context, id string) error
}

type UserRepository struct {
    DB DBTX // Asumsi: DBTX ini interface yang punya method sqlx (GetContext, SelectContext, dll)
}

var _ UserRepositoryContract = (*UserRepository)(nil)

func NewUserRepository(db DBTX) UserRepositoryContract {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
    query := `INSERT INTO users (username, email, password_hash, full_name, avatar_url, bio)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`

    err := r.DB.QueryRowContext(ctx, query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.AvatarURL,
        user.Bio).Scan(&user.ID)

    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }

    return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    // Diperbaiki: Tambah koma setelah bio, dan tambah IS sebelum NULL
    query := `SELECT id, username, email, full_name, avatar_url, bio, is_active, email_verified
    FROM users
    WHERE id=$1 AND deleted_at IS NULL`

    var user model.User
    err := r.DB.GetContext(ctx, &user, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("get user by id: %w", err)
    }
    return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    // Diperbaiki: Ubah id=$1 jadi email=$1
    query := `SELECT id, username, email, full_name, avatar_url, bio, is_active, email_verified
    FROM users
    WHERE email=$1 AND deleted_at IS NULL`

    var user model.User
    if err := r.DB.GetContext(ctx, &user, query, email); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("get user by email: %w", err)
    }
    return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    // Diperbaiki: Ubah id=$1 jadi username=$1
    query := `SELECT id, username, email, full_name, avatar_url, bio, is_active, email_verified
    FROM users
    WHERE username=$1 AND deleted_at IS NULL`

    var user model.User
    if err := r.DB.GetContext(ctx, &user, query, username); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("get user by username: %w", err)
    }
    return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*model.User, error) {
    // Diperbaiki: Tambah koma setelah bio, dan tambah IS sebelum NULL
    query := `SELECT id, username, email, full_name, avatar_url, bio, is_active, email_verified
    FROM users
    WHERE deleted_at IS NULL`

    var users []*model.User
    err := r.DB.SelectContext(ctx, &users, query)
    if err != nil {
        return nil, fmt.Errorf("list users: %w", err)
    }

    return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
    // Diperbaiki: Tambah koma untuk setiap SET, ubah fullname jadi full_name, tambah IS NULL
    query := `
        UPDATE users
        SET username = $1,
            full_name = $2,
            bio = $3,
            updated_at = NOW()
        WHERE id = $4 AND deleted_at IS NULL
    `
    _, err := r.DB.ExecContext(ctx, query, user.Username, user.FullName, user.Bio, user.ID)

    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
    query := `
    UPDATE users
    SET deleted_at = NOW()
    WHERE id = $1 AND deleted_at IS NULL`
    
    _, err := r.DB.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    return nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, id, passwordHash string) error {
    query := `
    UPDATE users
    SET password_hash = $1,
        updated_at = NOW() 
    WHERE id = $2 AND deleted_at IS NULL`

    // Tambahan minor: Biasanya saat update password, updated_at juga ikut diubah
    _, err := r.DB.ExecContext(ctx, query, passwordHash, id)
    if err != nil {
        return fmt.Errorf("update user password: %w", err)
    }

    return nil
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, id string) error {
    query := `
        UPDATE users
        SET email_verified = TRUE,
            updated_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
    _, err := r.DB.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("mark email verified: %w", err)
    }

    return nil
}