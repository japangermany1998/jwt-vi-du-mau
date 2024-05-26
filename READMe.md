## Điều kiện
- Khởi tạo server postgres có giá trị sau:
  ```
	Addr:     ":5432",
	User:     "postgres",
	Password: "123456",
	Database: "postgres",
  ```
- Tải [Postman](https://www.postman.com/downloads/) để chạy request

## Các bước chạy thử JWT.

### Bước 1:

Đầu tiên đăng ký user bằng cách gọi request sau và truyền dữ liệu JSON:

```
POST localhost:8080/register
```
Xem kết quả trả ra trên POSTMAN
![](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41an4c51co7eb8irtog)

Code của đoạn đăng ký này:
```go
func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		log.Print(err)
		return err
	}

	if data["password"] != data["passwordconfirm"] {   //Kiểm tra password có match ko
		ctx.Status(400)
		return ctx.JSON(map[string]string{
			"message": "password doesn't match",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)   //Mã hóa pass
	user := model.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  string(password),
	}

	_, err := DB.Model(&user).Insert()   //Insert database
	if err != nil {
		panic(err)
	}
	return ctx.JSON(user)
}
```

### Bước 2:
Đăng nhập bằng cách gọi request sau và truyền dữ liệu JSON:
```
POSt localhost:8080/user/login
```
Và kết quả trả ra token đăng nhập và cookies được lưu
![](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41aqms51co7eb8irtp0)
![](https://media.techmaster.vn/api/api/bm0tmgk51co4vclgfvu0/c41auks51co7eb8irtpg)

Code của phần này:
```go
func Login(ctx *fiber.Ctx) error {
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.JSON(err)
	}

	var user model.User
	err := DB.Model(&user).Where("email = ?", data["email"]).First()

	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(err)
	}

	if user.Id == 0 {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "not found",
		})
	}

	//So sánh password nhập vào với password băm trong database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		ctx.Status(400)

		return ctx.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := util.GenerateJWT(strconv.Itoa(user.Id)) //Tạo token đăng nhập
	if err != nil {
		log.Print(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{ 
		Name:     "jwt",  
		Value:    token, 
		Expires:  time.Now().Add(time.Hour * 24), //thời hạn 1 ngày
		HTTPOnly: true,  //Giới hạn quyền truy cập
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(token)
}
```

### Bước 3:
Thử một request lấy thông tin user đăng nhập
```
GET  localhost:8080/user
```
![](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41b04451co7eb8irtq0)

Code của phần này:
```go
func User(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")	//lấy value cookie

	issuer, _ := util.ParseJWT(cookie)	//Parse token lấy thông tin id của user đăng nhập

	var user model.User

	DB.Model(&user).Where("id = ?", issuer).Relation("Posts").Select()

	return ctx.JSON(user)
}
```

### Bước 4: 
Logout bằng request sau:
```
GET localhost:8080/logout
```
![](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41b1rs51co7eb8irtqg)
Code của phần này.
```go
func Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Logout success",
	})
}
```
Thử lại request lấy thông tin user đăng nhập
![d](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41b37k51co7eb8irtr0)

## Các function generate token và parse token
[Bấm vào đây để xem đầy đủ](./util/jwt.go):
```go
const SecretKey = "secret"

func GenerateJWT(issuer string) (string, error){	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	return claims.SignedString([]byte(SecretKey))
}
```
```go
func ParseJWT(cookie string) (string,error){
	token, err := jwt.ParseWithClaims(cookie,&jwt.StandardClaims{},func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"),nil
	})

	if err != nil || !token.Valid{
		return "",err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer,nil
}
```

## Middleware
[Bấm vào đây để xem đầy đủ](./util/jwt.go):
```go
func IsAuthenticated(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")

	if _, err := util.ParseJWT(cookie); err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	return ctx.Next()
}
```
Chạy thử một request /api trong trạng thái chưa đăng nhập:
![](https://media.techmaster.vn/api/static/bm0tmgk51co4vclgfvu0/c41b6j451co7eb8irtrg)

