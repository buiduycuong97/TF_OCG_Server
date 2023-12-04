package review_handle

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"time"
)

// Hàm cập nhật GetListReviewByProductID
func GetListReviewByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("ProductID is required"))
		return
	}

	pid, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	productID := int32(pid)

	variantIDs, err := dbms.GetVariantIDsByProductID(productID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var reviews []models.Review
	for _, variantID := range variantIDs {
		variantReviews, err := dbms.GetReviewsByVariantID(variantID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		reviews = append(reviews, variantReviews...)
	}

	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})

	var responseReviews []map[string]interface{}
	for _, review := range reviews {
		user, err := dbms.GetUserByID(review.UserID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		responseReview := map[string]interface{}{
			"reviewID":   review.ReviewID,
			"userID":     review.UserID,
			"variantID":  review.VariantID,
			"rating":     review.Rating,
			"comment":    review.Comment,
			"created_at": review.CreatedAt.Format(time.RFC3339),
			"user": map[string]interface{}{
				"userID":   user.UserID,
				"userName": user.UserName,
			},
		}

		responseReviews = append(responseReviews, responseReview)
	}

	res.JSON(w, http.StatusOK, responseReviews)
}
