package variant_handle

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/product_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

const (
	CredentialsPath       = "double2c-firebase-adminsdk-64fpu-cb7acf1b93.json"
	FirebaseStorageBucket = "double2c.appspot.com"
)

func UpdateVariantQuantityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	variantID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = UpdateVariantCountInStock(int32(variantID), int32(quantity))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Variant quantity updated successfully"})
}

func UpdateVariantCountInStock(variantID int32, newQuantity int32) error {
	variant := &models.Variant{}
	err := dbms.GetVariantById(variant, variantID)
	if err != nil {
		return errors.New("Failed to get variant")
	}

	if variant.CountInStock < newQuantity {
		return errors.New("Not enough quantity remaining")
	}

	variant.CountInStock -= newQuantity

	err = dbms.UpdateVariant(variant, variantID)
	if err != nil {
		return errors.New("Failed to update variant quantity")
	}

	return nil
}

func UpdateVariantQuantityWithIncrease(variantID int32, quantityToIncrease int32) error {
	variant := &models.Variant{}
	err := dbms.GetVariantById(variant, variantID)
	if err != nil {
		return errors.New("Failed to get variant")
	}

	variant.CountInStock += quantityToIncrease

	err = dbms.UpdateVariant(variant, variantID)
	if err != nil {
		return errors.New("Failed to update variant quantity")
	}

	return nil
}

func GetFilePath() (string, error) {
	absPath, err := filepath.Abs(CredentialsPath)
	if err != nil {
		return "Fail to get file path", err
	} else {
		return absPath, nil
	}
}

func UpdateVariantByAdmin(w http.ResponseWriter, r *http.Request) {
	variant := &models.Variant{}
	vars := mux.Vars(r)
	variantID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = dbms.GetVariantById(variant, int32(variantID))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	firebaseFilePath, error := GetFilePath()
	if error != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, error)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Limit your maxMultipartMemory
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		if err != http.ErrMissingFile {
			res.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		} else {
			countInStock, err := strconv.Atoi(r.PostFormValue("countInStock"))
			if err != nil {
				fmt.Println("Error parsing countInStock:", err)
				return
			}
			price, err := strconv.Atoi(r.PostFormValue("price"))
			if err != nil {
				fmt.Println("Error parsing price:", err)
				return
			}
			op1, err := strconv.Atoi(r.PostFormValue("optionValue1"))
			if err != nil {
				fmt.Println("Error parsing optionValue1:", err)
				return
			}
			op2, err := strconv.Atoi(r.PostFormValue("optionValue2"))
			if err != nil {
				fmt.Println("Error parsing optionValue2:", err)
				return
			}

			variant.Title = r.PostFormValue("title")
			variant.Price = int32(price)
			variant.CountInStock = int32(countInStock)
			variant.OptionValue1 = int32(op1)
			variant.OptionValue2 = int32(op2)

			err = dbms.UpdateVariantByAdmin(variant, int32(variantID))
			if err != nil {
				res.ERROR(w, http.StatusBadRequest, err)
				return
			}

			var product models.Product
			product, err = dbms.GetProductByID(variant.ProductID)

			err = utils.DeleteProductFromCache(product_handle.RedisClient, product.Handle)
			if err != nil {
				log.Println("Xóa sản phẩm trong cache thất bại: ", err)
			}

			res.JSON(w, http.StatusOK, variant)
		}
	} else {
		defer file.Close()

		// Khởi tạo Storage client
		ctx := context.Background()
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(firebaseFilePath))
		if err != nil {
			res.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		defer client.Close()

		// Tạo đường dẫn trong bucket để lưu file
		objectPath := "images/" + handler.Filename

		// Tạo đối tượng để ghi dữ liệu vào Firebase Storage
		wc := client.Bucket(FirebaseStorageBucket).Object(objectPath).NewWriter(ctx)
		if wc == nil {
			// Xử lý lỗi khi không lấy được đối tượng
			fmt.Printf("Error obtaining object writer: %v\n", err)
			return
		}

		// Copy dữ liệu từ file nhận được từ client vào đối tượng ghi (writer)
		if _, err := io.Copy(wc, file); err != nil {
			res.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		err = wc.Close()
		if err != nil {
			res.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		// Lấy thông tin chi tiết về tệp đã lưu trữ
		attrs := wc.Attrs()
		if attrs == nil {
			res.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		// Tạo URL để truy cập file đã lưu trong Firebase Storage
		imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", attrs.Bucket, attrs.Name)

		countInStock, err := strconv.Atoi(r.PostFormValue("countInStock"))
		if err != nil {
			fmt.Println("Error parsing countInStock:", err)
			return
		}
		price, err := strconv.Atoi(r.PostFormValue("price"))
		if err != nil {
			fmt.Println("Error parsing price:", err)
			return
		}
		op1, err := strconv.Atoi(r.PostFormValue("optionValue1"))
		if err != nil {
			fmt.Println("Error parsing optionValue1:", err)
			return
		}
		op2, err := strconv.Atoi(r.PostFormValue("optionValue2"))
		if err != nil {
			fmt.Println("Error parsing optionValue2:", err)
			return
		}

		variant.Price = int32(price)
		variant.CountInStock = int32(countInStock)
		variant.OptionValue1 = int32(op1)
		variant.OptionValue2 = int32(op2)
		variant.Image = imageURL

		err = dbms.UpdateVariantByAdmin(variant, int32(variantID))
		if err != nil {
			res.ERROR(w, http.StatusBadRequest, err)
			return
		}
		res.JSON(w, http.StatusOK, variant)

		var product models.Product
		product, err = dbms.GetProductByID(variant.ProductID)

		err = utils.DeleteProductFromCache(product_handle.RedisClient, product.Handle)
		if err != nil {
			log.Println("Xóa sản phẩm trong cache thất bại: ", err)
		}

	}

}
