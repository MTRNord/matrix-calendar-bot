package main

import (
	"encoding/json"
	"io"
	"os"
)

type config struct {
	MatrixBot configMatrixBot `json:"matrix_bot"`
	SQLiteURI string          `json:"sqlite_uri"`
}

type configMatrixBot struct {
	Homeserver string `json:"homeserver"`
	AccountID  string `json:"account_id"`
	Token      string `json:"token"`
}

type loadConfigError struct {
	create bool  // If the error occured while trying to create the file.
	err    error // Original error.
}

func (l loadConfigError) Error() string {
	return "cannot load config; " + l.err.Error()
}

var defaultConfig = config{
	MatrixBot: configMatrixBot{
		Homeserver: "https://example.org",
		AccountID:  "@calendarbot:remi.im",
		Token:      "",
	},
	SQLiteURI: "matrix-caldav-bot.db",
}

// loadConfig unmarhsals the contents of the file with given filename as JSON, which
// is returned.
// If there is no file with the given filename the file will be created containing
// defaultConfig as JSON. isNew is then set to true.
func loadConfig(filename string) (cfg config, isNew bool, err error) {
	f, err := os.Open(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return config{}, false, loadConfigError{false, err}
		}

		// Configuration file probably doesn't exist.
		// Create it with defaultConfig in JSON as contents.
		err = createConfig(filename, defaultConfig)
		if err != nil {
			return defaultConfig, true, loadConfigError{true, err}
		}

		return defaultConfig, true, nil
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&cfg)
	if err != nil && err != io.EOF {
		return cfg, isNew, loadConfigError{false, err}
	}

	return cfg, isNew, err
}

func createConfig(filename string, cfg config) (err error) {
	data, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
