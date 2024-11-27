package controller

import (
	"go-fiber-user-management/database"
	"go-fiber-user-management/model"
	"go-fiber-user-management/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUsers(c *fiber.Ctx) error {
	var userData []model.User
	// Mengambil data pengguna dengan urutan berdasarkan ID
	if err := database.DB.Order("id").Find(&userData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}

	// Menyiapkan respons DTO untuk setiap pengguna
	var usersResponse []model.UserResponseDTO
	for _, user := range userData {
		userResponse := model.UserResponseDTO{
			ID:          user.ID,
			Email:       user.Email,
			Fullname:    user.Fullname,
			Address:     user.Address,
			Gender:      user.Gender,
			PhoneNumber: user.PhoneNumber,
		}
		usersResponse = append(usersResponse, userResponse)
	}

	// Mengembalikan data pengguna dalam bentuk JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"data":    usersResponse,
	})
}

func GetDetailUser(c *fiber.Ctx) error {
	// Ambil parameter ID dari URL
	id := c.Params("id")

	// Cari data pengguna berdasarkan ID
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	// Membuat DTO respons untuk pengguna
	userResponse := model.UserResponseDTO{
		ID:          user.ID,
		Email:       user.Email,
		Fullname:    user.Fullname,
		Address:     user.Address,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
	}

	// Kembalikan respons dengan detail pengguna
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User details fetched successfully",
		"data":    userResponse,
	})
}

func CreateUser(c *fiber.Ctx) error {
	// Parsing input ke struct DTO
	var userRequest model.UserRequestDTO
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validasi input
	if userRequest.Email == "" || userRequest.Password == "" || len(userRequest.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and password are required, and password must be at least 6 characters",
		})
	}

	// Cek apakah email sudah ada
	var existingUser model.User
	if err := database.DB.Where("email = ?", userRequest.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Hash password
	hashedPassword := utils.GeneratePassword(userRequest.Password)

	// Buat model user
	userModel := model.User{
		Email:        userRequest.Email,
		PasswordHash: hashedPassword,
		Fullname:     userRequest.Fullname,
		Address:      userRequest.Address,
		Gender:       userRequest.Gender,
		PhoneNumber:  userRequest.PhoneNumber,
	}

	// Simpan data ke database
	response := database.DB.Create(&userModel)
	if response.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   response.Error.Error(),
		})
	}

	// Buat DTO untuk respons
	userResponse := model.UserResponseDTO{
		ID:          userModel.ID,
		Email:       userModel.Email,
		Fullname:    userModel.Fullname,
		Address:     userModel.Address,
		Gender:      userModel.Gender,
		PhoneNumber: userModel.PhoneNumber,
	}

	// Kembalikan respons sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    userResponse,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// Parsing input ke struct DTO
	var userRequest model.UserRequestDTO
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Ambil ID dari parameter URL
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID is required",
		})
	}

	// Temukan user berdasarkan ID
	var dataUser model.User
	result := database.DB.First(&dataUser, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user data",
			"error":   result.Error.Error(),
		})
	}

	// Perbarui data pengguna langsung pada objek `dataUser`
	dataUser.Email = userRequest.Email
	dataUser.Fullname = userRequest.Fullname
	dataUser.Address = userRequest.Address
	dataUser.Gender = userRequest.Gender
	dataUser.PhoneNumber = userRequest.PhoneNumber

	// Hanya set PasswordHash jika password baru disediakan
	if userRequest.Password != "" {
		hashedPassword := utils.GeneratePassword(userRequest.Password)
		dataUser.PasswordHash = hashedPassword
	}

	// Simpan perubahan ke database
	result = database.DB.Save(&dataUser)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   result.Error.Error(),
		})
	}

	// Buat response DTO tanpa password
	userResponse := model.UserResponseDTO{
		ID:          dataUser.ID,
		Email:       dataUser.Email,
		Fullname:    dataUser.Fullname,
		Address:     dataUser.Address,
		Gender:      dataUser.Gender,
		PhoneNumber: dataUser.PhoneNumber,
	}

	// Kembalikan respons sukses dengan data tanpa password
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    userResponse,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	// Ambil ID dari parameter
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID is required",
		})
	}

	// Cari data user berdasarkan ID
	var userData model.User
	result := database.DB.First(&userData, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user data",
			"error":   result.Error.Error(),
		})
	}

	// Jalankan penghapusan data
	if err := database.DB.Delete(&userData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	// Kembalikan respons sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"data":    userData,
	})
}
