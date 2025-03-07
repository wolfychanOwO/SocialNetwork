package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.POST("/register", func(c *gin.Context) {
        proxyRequest(c, "http://user-service:8081/register")
    })

    r.POST("/authenticate", func(c *gin.Context) {
        proxyRequest(c, "http://user-service:8081/authenticate")
    })

    r.PUT("/user", func(c *gin.Context) {
        proxyRequest(c, "http://user-service:8081/user")
    })

    r.GET("/user/:login", func(c *gin.Context) {
        proxyRequest(c, "http://user-service:8081/user/"+c.Param("login"))
    })

    if err := http.ListenAndServe(":8080", r); err != nil {
        panic(err)
    }
}

func proxyRequest(c *gin.Context, url string) {
    client := &http.Client{}
    req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    req.Header = c.Request.Header
    resp, err := client.Do(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

