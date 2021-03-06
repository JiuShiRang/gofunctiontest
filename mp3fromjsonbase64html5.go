package gofunctiontest

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
)

// This function intends to produce some mp3 files from base64-encoded strings which are read from a json file at first.
// The json file was made from the website https://aidn.jp/mikutap/ , from its track.js and main.js when I used Network functionality of chrome brower.
// I find this function unnecessary for my purpose, but I will still take it as a test.
// At last, I want to play them one by one
func Mp3FromBase64(jsonfile, outputdir string) error {
	// First, to get strings from json file
	// Second, use /Regular Expression/  to extract code for each pair
	// Third,  produce file for each pair
	// So I need a function which take Step 2-3 for each pair

	log.Logger = log.With().Caller().Logger()

	mp3pair := make(map[string]string)
	bytes, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		// fmt.Println("ReadFile: ", err.Warn())
		log.Warn().Str("err", err.Error()).Str("filepath", jsonfile).Msg("ioutil.ReadFile")
		return err
	}

	if err := json.Unmarshal(bytes, &mp3pair); err != nil {
		// fmt.Println("Unmarshal: ", err.Warn())
		log.Warn().Str("err", err.Error()).Str("filepath", jsonfile).Msg("json.Unmarshal")
		return err
	}
	log.Info().Int("Num", len(mp3pair)).Msg("Strings have been read")

	// sync/waitgroup seems to be a little faster, I donot use it though
	c := make(chan error, len(mp3pair))
	for k, v := range mp3pair {
		go generateFileForPair(k, v, outputdir, c)
	}
	for i := 0; i < len(mp3pair); i++ {
		<-c
		// for i := range c {
		// log.Info().Str("Warn", i.Error()).Msg("")
	}
	log.Info().Msg("Files have been generated")

	return nil
}

// To extract a base64-encoded string and write a file at the given dir
func generateFileForPair(fileName string, strEncoded string, dirOutput string, c chan error) {
	strEncoded = strEncoded[strings.Index(strEncoded, ",")+1:]
	// log.Info().Str("base64", strEncoded).Int("num", len(strEncoded)).Msg("")

	// The situation in which + & / are replaced by - & _ means urlSafeEncode instead of stdEncode
	sDec, err := base64.StdEncoding.DecodeString(strEncoded)
	if err != nil {
		log.Warn().Str("Warn", err.Error()).Msg("base64.StdEncoding.DecodeString")
		c <- err
		return
	}

	if err := ioutil.WriteFile(dirOutput+"/"+fileName, sDec, 0644); err != nil {
		log.Warn().Str("Warn", err.Error()).Msg("ioutil.WriteFile")
		c <- err
		return
	}

	c <- nil
}
