package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RequestLogger(c *gin.Context) {
	defer c.Next()

	tmp := json.RawMessage{}
	err := c.ShouldBindJSON(&tmp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	req := ""
	buf := bytes.NewBufferString(req)

	err = json.NewEncoder(buf).Encode(tmp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}
	fmt.Printf("Body request: %s\n", tmp)
}
