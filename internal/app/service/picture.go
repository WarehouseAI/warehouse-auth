package service

import (
	"auth-service/internal/app/dataservice"
	e "auth-service/internal/pkg/errors/http"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func UploadImage(pic *multipart.FileHeader, picture dataservice.PictureInterface, logger *logrus.Logger) (string, error) {
	picPayload, err := pic.Open()

	if err != nil {
		return "", e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
	}

	url, fileErr := picture.UploadFile(picPayload, fmt.Sprintf("avatar.%s%s", uuid.Must(uuid.NewV4()).String(), filepath.Ext(pic.Filename)))

	if fileErr != nil {
		return "", e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
	}

	defer picPayload.Close()
	return url, nil
}
