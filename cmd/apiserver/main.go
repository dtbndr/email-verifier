package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	emailVerifier "github.com/AfterShip/email-verifier"
)

func GetEmailVerification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	verifier := emailVerifier.NewVerifier().
		EnableSMTPCheck().
		HelloName("solventfunding.com").
		FromEmail("max@solventfunding.com").
		ConnectTimeout(30 * time.Second)

	// Enable the API verifier for Yahoo (also covers AOL)
	if err := verifier.EnableAPIVerifier("yahoo"); err != nil {
		// This setup is unlikely to fail, but we log it just in case.
		log.Printf("Warning: Failed to enable Yahoo API verifier: %v", err)
	}

	log.Println("DEBUG: Before verifier.Verify()")
	ret, err := verifier.Verify(ps.ByName("email"))
	log.Println("DEBUG: After verifier.Verify()")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !ret.Syntax.Valid {
		_, _ = fmt.Fprint(w, "email address syntax is invalid")
		return
	}

	bytes, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, _ = fmt.Fprint(w, string(bytes))

}

func main() {
	router := httprouter.New()

	router.GET("/v1/:email/verification", GetEmailVerification)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
