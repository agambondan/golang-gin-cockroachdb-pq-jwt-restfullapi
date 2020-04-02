package controllers

import (
	"../models"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	_ "strconv"
	"time"
)

func (server *Server) ExportUserToExcel(c *gin.Context) {
	xlsx := excelize.NewFile()
	sheet1Name := "user"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	xlsx.SetCellValue(sheet1Name, "A1", "id")
	xlsx.SetCellValue(sheet1Name, "B1", "created_at")
	xlsx.SetCellValue(sheet1Name, "C1", "updated_at")
	xlsx.SetCellValue(sheet1Name, "D1", "deleted_at")
	xlsx.SetCellValue(sheet1Name, "E1", "full_name")
	xlsx.SetCellValue(sheet1Name, "F1", "phone_number")
	xlsx.SetCellValue(sheet1Name, "G1", "username")
	xlsx.SetCellValue(sheet1Name, "H1", "password")
	xlsx.SetCellValue(sheet1Name, "I1", "email")
	xlsx.SetCellValue(sheet1Name, "J1", "role_id")

	err := xlsx.AutoFilter(sheet1Name, "A1", "J1", "")
	if err != nil {
		log.Fatal(err)
	}

	var users models.User
	findAllUser, err := users.FindAllUser(server.DB)
	for i, each := range findAllUser {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), each.ID)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), each.CreatedAt.String())
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), each.UpdatedAt.String())
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), each.DeletedAt.String())
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), each.FullName)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), each.Username)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), each.Password)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), each.Email)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), each.RoleId)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), each.PhoneNumber)
	}

	err = xlsx.SaveAs("./assets/excel/user" + time.Now().Local().String() + ".xlsx")
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Export successfully with name user-" + time.Now().Local().String() + ".xlsx"})
}

func (server *Server) ImportExcelToUser(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + " hmm"})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		regex := after(basename, ".")
		if regex != "xlsx" && regex != "xls" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "file format must be xlsx or xls"})
			return
		}
		filename := filepath.Join("./assets/excel/", time.Now().Local().String()+basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + " asem"})
			return
		}
		// open file xlsx && Insert into database
		xlsx, err := excelize.OpenFile(filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + " asem2"})
		}
		var data []string
		rows, err := xlsx.GetRows("user")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		for i, row := range rows {
			if i != 0 {
				for j := 0; j < len(row); j++ {
					data = append(data, row[j])
					if j+1 == len(row) {
						stmt, err := server.DB.Prepare("INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
						if err != nil {
							c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
						}
						_, err = stmt.Exec(data[0], convertStringToDate(data[1]), convertStringToDate(data[2]), convertStringToDate(data[3]), data[4], data[5], data[6], data[7], data[8], data[9])
						if err != nil {
							c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
						}
						data = nil
					}
				}
			}
		}
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Filename)
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusAccepted, "message": "upload ok!", "data": gin.H{"files": filenames}})
}

func convertStringToDate(string string) *time.Time {
	if string == "NULL" || string == "0001-01-01 00:00:00 +0000 UTC" {
		return nil
	}
	layoutFormat := "2006-01-02 00:00:00 +0000 +0000"
	date, err := time.Parse(layoutFormat, string)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &date
}
