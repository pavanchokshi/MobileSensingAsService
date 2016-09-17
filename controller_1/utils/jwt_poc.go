package main

import(
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	fmt.Println("JWT tutorial: Creating Token")
	 // Create the token
    token := jwt.New(jwt.SigningMethodHS256)
    // Set some claims
    token.Claims["foo"] = "bar"
    token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
    // Sign and get the complete encoded token as a string
    tokenString,_ := token.SignedString([]byte("this is dhruv"))
    fmt.Println(tokenString)

    fmt.Println("JWT tutorial: Verify Token")

    token, _ = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        
        //return fmt.Println(token)
         //Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return token.Header["kid"]
    })

}