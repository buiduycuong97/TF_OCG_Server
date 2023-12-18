package variant_handle

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/product_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func CreateVariantHandler(w http.ResponseWriter, r *http.Request) {
	var variant models.Variant

	firebaseFilePath, error := GetFilePath()
	if error != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, error)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Limit your maxMultipartMemory
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
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
	productId, err := strconv.Atoi(r.PostFormValue("productId"))
	if err != nil {
		fmt.Println("Error parsing productId:", err)
		return
	}
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

	variant.ProductID = int32(productId)
	variant.Price = int32(price)
	variant.CountInStock = int32(countInStock)
	variant.OptionValue1 = int32(op1)
	variant.OptionValue2 = int32(op2)
	variant.Image = imageURL

	variantRes, err := dbms.CreateVariant(&variant)
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

	res.JSON(w, http.StatusCreated, variantRes)
}
