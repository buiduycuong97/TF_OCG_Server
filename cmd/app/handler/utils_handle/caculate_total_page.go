package utils_handle

func CalculateTotalPages(totalCount int64, pageSize int32) int32 {
	totalPages := int32(totalCount) / pageSize
	if int32(totalCount)%pageSize != 0 {
		totalPages++
	}
	return totalPages
}
