package postgres

import (
	"auth-service/generated/user"
	"fmt"
	"log"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into users(id, fullname, username, email, password, phone, image)
	values ($1, $2, $3, $4, $5, $6, $7)`,
		"3c9c52ef-8328-4801-b938-25f060f9337a", "test name", "test", "test@gmail.com",
		"test", "+987654321123", "test")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepositoryImpl{db: db}

	userRes, err := userRepo.GetUserByID("3c9c52ef-8328-4801-b938-25f060f9337a")
	if err != nil {
		log.Fatal(err)
	}
	if userRes.Id != "3c9c52ef-8328-4801-b938-25f060f9337a" {
		log.Fatal("false user")
	}

	_, err = db.Exec(`delete from users where id = $1`, userRes.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllUsers(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into users(id, fullname, username, email, password, phone, image)
	values ($1, $2, $3, $4, $5, $6, $7)`,
		"3c9c52ef-8328-4801-b938-25f060f9337b", "test name", "test", "test@gmail.com",
		"test", "+987654321123", "test")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepositoryImpl{db: db}

	filter := user.GetAllUsersRequest{}

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := userRepo.GetAllUsers(&filter)
		if err != nil {
			log.Fatal(err)
		}
		if users == nil {
			log.Fatal("false user")
		}
	})
}

func TestUpdateUser(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//_, err = db.Exec(`insert into users(id, fullname, username, email, password, phone, image)
	//	values ($1, $2, $3, $4, $5, $6, $7)`,
	//	"3c9c52ef-8328-4801-b938-25f060f9337b", "test name", "test", "test@gmail.com",
	//	"test", "+987654321123", "test")
	//if err != nil {
	//	log.Fatal(err)
	//}

	userRepo := userRepositoryImpl{db: db}
	updateUser := user.UpdateUserRequest{
		Id:          "3c9c52ef-8328-4801-b938-25f060f9337a",
		Username:    "new_username",
		Password:    "test",
		NewPasswrod: "new_password",
		Email:       "new_email",
		FullName:    "new_full_name",
		PhoneNumber: "+998999774754",
	}
	user, err := userRepo.UpdateUser(&updateUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}

func TestUpdateUserByID(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into users(id, fullname, username, email, password, phone, image)
			values ($1, $2, $3, $4, $5, $6, $7)`,
		"3c9c52ef-8328-4801-b938-25f060f12365", "test name", "test", "test@gmail.com",
		"test", "+987654321123", "test")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userRepositoryImpl{db: db}

	t.Run("Update user", func(t *testing.T) {
		updateUser := user.UpdateUserByIdRequest{
			Id:       "3c9c52ef-8328-4801-b938-25f060f12365",
			FullName: "old_full_name",
			Username: "old_username",
			Email:    "old_email",
		}
		us, err := userRepo.UpdateUserById(&updateUser)
		if err != nil {
			log.Fatal(err)
		}
		if us == nil {
			log.Fatal("nil user")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := userRepositoryImpl{db: db}
	id := user.DeleteUserRequest{
		Id:       "3c9c52ef-8328-4801-b938-25f060f65456",
		Password: "new_password",
	}
	res, err := userRepo.DeleteUser(&id)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatalln(res)
}
