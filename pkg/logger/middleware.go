package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func RequestLogger(c *gin.Context) {
	defer c.Next()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}

	tmp := json.RawMessage{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	req := ""
	buf := bytes.NewBufferString(req)

	err = json.NewEncoder(buf).Encode(tmp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Body request: %s\n", tmp)
}
