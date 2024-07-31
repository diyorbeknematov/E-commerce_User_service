package models

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
	ExpiresAt    string `json:"expires_at"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Errors struct {
	Error string `json:"error"`
}

type GettAllUserResponse struct {
	Users      []*GetUserResponse `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	TotalCount int32              `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"` 
}

type GetUserResponse struct {
	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" db:"id,omitempty"`
	FullName string `protobuf:"bytes,2,opt,name=fullName,proto3" json:"fullName,omitempty" db:"fullname,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty" db:"username,omitempty"`
	Phone    string `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty" db:"phone,omitempty"`
	Email    string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty" db:"email,omitempty"`
	Image    string `protobuf:"bytes,6,opt,name=image,proto3" json:"image,omitempty" db:"image,omitempty"`
	Role     string `protobuf:"bytes,7,opt,name=role,proto3" json:"role,omitempty" db:"role,omitempty"`
}