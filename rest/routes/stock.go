package routes

import (
	"encoding/csv"
	"fmt"
	"gudang-api-go-v1/models"
	"gudang-api-go-v1/rest/response"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSVRow struct {
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	Quantity    uint   `json:"quantity"`
}

func Stock(router *gin.RouterGroup, dbs map[string]*gorm.DB) {
	router.GET("/stocks", func(ctx *gin.Context) {
		db := dbs["main"].Session(&gorm.Session{})

		type StockWithProduct struct {
			models.Stock
			ProductCode   string `json:"product_code"`
			ProductName   string `json:"product_name"`
			ProductStatus string `json:"product_status"`
		}

		var stocksWithProducts []StockWithProduct

		if err := db.Table("stocks").
			Select("stocks.*, products.code as product_code, products.name as product_name, products.status as product_status").
			Joins("LEFT JOIN products ON stocks.product_id = products.id").
			Scan(&stocksWithProducts).Error; err != nil {
			response.Fatal(ctx, 500, err.Error())
			return
		}

		type StockDetail struct {
			Id            uint      `json:"id"`
			ProductId     uint      `json:"product_id"`
			Quantity      uint      `json:"quantity"`
			CreatedAt     time.Time `json:"created_at"`
			CreatedBy     uint      `json:"created_by"`
			ProductCode   string    `json:"product_code"`
			ProductName   string    `json:"product_name"`
			ProductStatus string    `json:"product_status"`
		}

		stockDetails := make([]StockDetail, len(stocksWithProducts))
		for i, s := range stocksWithProducts {
			stockDetails[i] = StockDetail{
				Id:            s.Id,
				ProductId:     s.ProductId,
				Quantity:      s.Quantity,
				CreatedAt:     s.CreatedAt,
				CreatedBy:     s.CreatedBy,
				ProductCode:   s.ProductCode,
				ProductName:   s.ProductName,
				ProductStatus: s.ProductStatus,
			}
		}

		response.OkWithTabledata(ctx, 200, "products loaded", stockDetails, gin.H{})
	})
	router.POST("/stocks/upload", func(ctx *gin.Context) {
		db := dbs["main"].Session(&gorm.Session{})

		file, err := ctx.FormFile("file")
		if err != nil {
			response.Fatal(ctx, 400, "please upload a csv file")
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			response.Fatal(ctx, 500, "error opening file: "+err.Error())
			return
		}
		defer openedFile.Close()

		reader := csv.NewReader(openedFile)
		header, err := reader.Read()
		if err != nil {
			response.Fatal(ctx, 500, "error reading CSV header: "+err.Error())
			return
		}

		if !validateHeader(header) {
			response.Fatal(ctx, 400, "invalid CSV format. expected columns: product_code, product_name, stock")
			return
		}

		var processed int
		var errors []string

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errors = append(errors, fmt.Sprintf("error reading row: %v", err))
				continue
			}

			row := CSVRow{
				ProductCode: strings.TrimSpace(record[0]),
				ProductName: strings.TrimSpace(record[1]),
				Quantity:    parseUint(record[2]),
			}

			if err := validateRow(row); err != nil {
				errors = append(errors, err.Error())
				continue
			}

			err = processStockRow(db, row)
			if err != nil {
				errors = append(errors, fmt.Sprintf("error processing row: %v", err))
				continue
			}

			processed++
		}

		if len(errors) > 0 {
			response.OkWithData(ctx, 200, "import completed with errors", gin.H{
				"processed": processed,
				"errors":    errors,
			})
			return
		}

		response.Ok(ctx, 200, fmt.Sprintf("successfully processed %d rows", processed))
	})
}

func processStockRow(db *gorm.DB, row CSVRow) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var product models.Product
	var stock models.Stock

	// Check if product exists by code
	if err := tx.Where("code = ?", row.ProductCode).First(&product).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}

		// Create new product if it doesn't exist
		product = models.Product{
			Code:      row.ProductCode,
			Name:      row.ProductName,
			Status:    "active",
			CreatedBy: 1,
		}

		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Check if stock exists for this product
	err := tx.Where("product_id = ?", product.Id).First(&stock).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	if stock.Id == 0 {
		// Create new stock entry
		stock = models.Stock{
			ProductId: product.Id,
			Quantity:  row.Quantity,
			CreatedBy: 1,
		}
		if err := tx.Create(&stock).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Add the new stock quantity to the existing stock
		newStock := stock.Quantity + row.Quantity
		if err := tx.Model(&stock).Update("quantity", newStock).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func validateHeader(header []string) bool {
	expectedHeader := []string{"product_code", "product_name", "quantity"}
	if len(header) != len(expectedHeader) {
		return false
	}
	for i, h := range header {
		if strings.ToLower(strings.TrimSpace(h)) != expectedHeader[i] {
			return false
		}
	}
	return true
}

func validateRow(row CSVRow) error {
	if row.ProductCode == "" {
		return fmt.Errorf("product code is required")
	}
	if row.ProductName == "" {
		return fmt.Errorf("product name is required")
	}
	if row.Quantity == 0 {
		return fmt.Errorf("stock must be greater than 0")
	}
	return nil
}

func parseUint(s string) uint {
	val, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	return uint(val)
}
