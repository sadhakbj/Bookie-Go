package common

import "github.com/sadhakbj/bookie-go/src/internal/helpers"

type ResponseDTO struct {
	Success    bool                   `json:"success"`
	Data       any                    `json:"data"`
	Pagination helpers.PaginationInfo `json:"pagination"`
}
