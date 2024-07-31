package postgres

import (
	pb "auth-service/generated/user"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUserByID(id string) (*pb.GetUserResponse, error)
	GetAllUsers(fUser *pb.GetAllUsersRequest) ([]*pb.GetUserResponse, error)
	UpdateUser(user *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	UpdateUserById(user *pb.UpdateUserByIdRequest) (*pb.UpdateUserByIdResponse, error)
	DeleteUser(req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	DeleteUserByID(req *pb.DeleteByIdRequest) (*pb.DeleteUserResponse, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) CreateUser(user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var id string
	err = tx.QueryRow(`
        INSERT INTO users (
            username,
            fullname,
            email,
            phone,
			image,
			role,
			password
        ) VALUES ($1, $2, $3, $4, 'Unknown Image', $5, $6) RETURNING id`,
		user.Username,
		user.FullName,
		user.Email,
		user.Phone,
		user.Role,
		user.Password,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO user_locations (
			user_id,
			address,
			city,
			state,
			country,
			postal_code
		)
		VALUES (
			$1,
			'Unknown address',
			'Unknown city',
			'Unknown state',
            'Unknown country',
            'Unknown postal_code'
		)
	`, id)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (r *userRepositoryImpl) GetUserByID(id string) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse
	err := r.db.Get(&user, `
		SELECT 
			id,
			username,
            fullname,
            email,
			phone,
			image
		FROM 
			users 
		WHERE 
			id=$1 AND deleted_at = 0`, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetAllUsers(fUser *pb.GetAllUsersRequest) ([]*pb.GetUserResponse, error) {
	var (
		args   []interface{}
		filter string
	)

	query := `
        SELECT 
            id,
            username,
            fullname,
            email,
            phone,
            image,
			role
        FROM	
            users u
        JOIN
            user_locations ul ON u.id = ul.user_id
        WHERE deleted_at = 0 `

	if fUser.FullName != "" {
		filter += fmt.Sprintf(" AND u.full_name LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.FullName))
	}
	if fUser.City != "" {
		filter += fmt.Sprintf(" AND ul.city LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.City))
	}
	if fUser.State != "" {
		filter += fmt.Sprintf(" AND ul.state LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.State))
	}
	if fUser.Country != "" {
		filter += fmt.Sprintf(" AND ul.country LIKE $%d", len(args)+1)
		args = append(args, fmt.Sprintf("%%%s%%", fUser.Country))
	}

	// `LIMIT` va `OFFSET` qiymatlarini qo'shish
	filter += fmt.Sprintf(" LIMIT %d OFFSET %d", fUser.Limit, (fUser.Offset-1)*fUser.Limit)

	query += filter

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var users []*pb.GetUserResponse
	for rows.Next() {
		var user pb.GetUserResponse
		err = rows.Scan(
				&user.Id, 
				&user.Username, 
				&user.FullName, 
				&user.Email, 
				&user.Phone, 
				&user.Image,
				&user.Role,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, err
}

func (r *userRepositoryImpl) UpdateUser(user *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	_, err = tx.Exec(`
        UPDATE 
            users 
        SET 
            username=$1,
            fullname=$2,
            email=$3,
            phone=$4,
            image=$5,
			password=$6
        WHERE 
            id=$7 and password=$8 and deleted_at = 0`,
		user.Username, user.FullName, user.Email,
		user.PhoneNumber, user.Image, user.NewPasswrod, user.Id, user.Password)

	if err != nil {
		return &pb.UpdateUserResponse{
			Success: false,
			Message: "Error updating user",
		}, err
	}

	_, err = tx.Exec(`
		UPDATE 
			user_locations
		SET
			address = $1,
			city = $2,
			state = $3,
			country = $4,
			postal_code = $5
		WHERE
			user_id = $6
	`, user.Address, user.City, user.State, user.Country, user.PostalCode, user.Id)

	if err != nil {
		return &pb.UpdateUserResponse{
			Success: false,
			Message: "Error updating user",
		}, err
	}

	return &pb.UpdateUserResponse{
		Success: true,
		Message: "User updated successfully",
	}, err
}

func (r *userRepositoryImpl) UpdateUserById(user *pb.UpdateUserByIdRequest) (*pb.UpdateUserByIdResponse, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err := tx.Exec(`
        UPDATE 
            users 
        SET 
            username=$1,
            fullname=$2,
            email=$3,
            phone=$4,
            image=$5
        WHERE 
            id=$6 AND deleted_at = 0`,
		user.Username, user.FullName, user.Email,
		user.Phone, user.Image, user.Id)

	if err != nil {
		return &pb.UpdateUserByIdResponse{
			Success: false,
			Message: "Error updating user",
		}, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("user not found")
	}
	_, err = tx.Exec(`
		UPDATE 
			user_locations
		SET
			address = $1,
			city = $2,
			state = $3,
			country = $4,
			postal_code = $5
		WHERE
			user_id = $6
	`, user.Address, user.City, user.State, user.Country, user.PostalCode, user.Id)

	if err != nil {
		return &pb.UpdateUserByIdResponse{
			Success: false,
			Message: "Error updating user",
		}, err
	}
	return &pb.UpdateUserByIdResponse{
		Success: true,
		Message: "User updated successfully",
	}, nil
}

func (r *userRepositoryImpl) DeleteUser(req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := r.db.Exec(`
		UPDATE 
			users 
        SET 
            deleted_at=DATE_PART('epoch', CURRENT_TIMESTAMP)::INT
        WHERE 
            id=$1 AND  password = $2 and deleted_at = 0
	`, req.Id, req.Password)

	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: "Error deleting user",
		}, err
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

func (r *userRepositoryImpl) DeleteUserByID(req *pb.DeleteByIdRequest) (*pb.DeleteUserResponse, error) {
	_, err := r.db.Exec(`
		UPDATE 
			users 
        SET 
            deleted_at=DATE_PART('epoch', CURRENT_TIMESTAMP)::INT
        WHERE 
            id=$1 AND deleted_at = 0
	`, req.UserId)

	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: "Error deleting user",
		}, err
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}
