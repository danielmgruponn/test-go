package handlers

import (
	"test-go/internal/core/ports"
	"test-go/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileService ports.FileService
}

func NewFileHandler(fileService ports.FileService) *FileHandler {
	return &FileHandler{fileService}
}

func (f *FileHandler) UploadFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al subir archivos"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No se encontraron archivos"})
	}

	updloadFiles, err := f.fileService.UploadFiles(files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al subir archivos"})
	}

	return c.Status(fiber.StatusOK).JSON(updloadFiles)
}

func (f *FileHandler) SaveFile(file dto.FileAttachment) error {
	err := f.fileService.SaveFile(&file)
	return err

}
