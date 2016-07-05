package server

import "github.com/gin-gonic/gin"

// ListenHTTP start http server
func ListenHTTP(port string) {
	logger.Infof("listen http %s", port)
	if conf.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
		c.Next()
	})

	// router
	r.OPTIONS("/*cors", func(c *gin.Context) {})

	if err := r.Run(port); err != nil {
		logger.Criticalf("Cannot start http server: %s", err)
	}
	Exit()
}
