package game

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/danmrichards/dessego/internal/transport/encoding/form"
)

func (s *Server) initCharacterHandler() http.HandlerFunc {
	type initCharacterReq struct {
		CharacterID string `form:"characterID"`
		Index       int    `form:"index"`
		Version     int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"init character request from %q to %q", r.RemoteAddr, s.l.Addr(),
		)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Demon's Souls sends it's request body as an AES encrypted version of
		// a standard HTTP form POST.

		// TODO: Move to dependency of the game server.
		// Check out the super secure AES key. The real game used this...
		c, err := aes.NewCipher([]byte("11111111222222223333333344444444"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Block-chain mode decrypter, pulling out the initialisation vector
		// (IV) from the body.
		cbc := cipher.NewCBCDecrypter(c, b[:aes.BlockSize])

		// Decrypt everything in the body, following the IV.
		req := make([]byte, len(b[aes.BlockSize:]))
		cbc.CryptBlocks(req, b[aes.BlockSize:])

		// Trim the trailing characters.
		trim := int(req[len(req)-1])
		req = req[:len(req)-trim]

		// Can now use the body as a normal HTTP form.
		v, err := url.ParseQuery(string(req))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var icr initCharacterReq
		if err = form.NewDecoder(v).Decode(&icr); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", icr)

		// TODO: Create player in DB
		// TODO: Add player to active list
	}
}
