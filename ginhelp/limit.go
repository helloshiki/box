package ginhelp

import (
	"bytes"
	"fmt"
	"errors"
)

func ParseLimit(s string) (int, int, error) {
	var offset, limit uint
	buf := bytes.NewBufferString(s)
	n , err := fmt.Fscanf(buf, "%d,%d", &offset, &limit)
	if err != nil {
		return 0, 0, errors.New("Invalid Params")
	}

	if n != 2 {
		return 0, 0, errors.New("Invalid Params")
	}

	if limit == 0 || limit > 1000 {
		return 0, 0, errors.New("Invalid Params")
	}

	return int(offset), int(limit) , nil
}