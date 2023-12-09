package types

import "strconv"

type GetSets struct {
	Limit  uint64
	Offset uint64
}

func ValidateGetSets(limit_str, offset_str string) (limit, offset uint64, res bool) {
	if limit_str == "" || offset_str == "" {
		return
	}

	limit_int, errl := strconv.Atoi(limit_str)
	offset_int, erro := strconv.Atoi(offset_str)

	if errl != nil || erro != nil {
		return
	}

	if limit_int < 1 || offset_int < 0 {
		return
	}

	res = true
	limit = uint64(limit_int)
	offset = uint64(offset_int)

	return
}
