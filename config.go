package typ3r

import (
	"fmt"
	"os"
	"os/user"
	"path"

	ini "gopkg.in/ini.v1"
)

const confFile = ".typ3rrc"

const rootSection = ""
const defaultURL = "https://typ3r.com/api/index.php"

var confPath string

// Config holds all typ3r's client user configuration
type Config struct {
	serverURL string
	token     string
	user      string
}

func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	confPath = path.Join(usr.HomeDir, confFile)
}

func getOption(sectionName string, file *ini.File, keyName string, def string) string {
	var err error
	var key *ini.Key
	var section *ini.Section

	section, err = file.GetSection(sectionName)

	if err != nil {
		return def
	}

	key, err = section.GetKey(keyName)

	if err != nil || key.Value() == "" {
		return def
	}

	return key.Value()
}

func getConfig(sectionName string, file *ini.File, keyName string) (cfg string, err error) {
	var key *ini.Key
	var section *ini.Section

	if section, err = file.GetSection(sectionName); err != nil {
		return
	}

	if key, err = section.GetKey(keyName); err != nil {
		return
	}

	cfg = key.Value()

	return
}

// Load loads a typ3r's client user configuration
func Load() (cfg *Config, err error) {
	var parsed *ini.File

	if parsed, err = ini.Load(confPath); err != nil {
		return
	}

	cfg = new(Config)

	cfg.serverURL = getOption(rootSection, parsed, "url", defaultURL)

	if (*cfg).token, err = getConfig(rootSection, parsed, "token"); err != nil {
		return
	}

	if cfg.user, err = getConfig(rootSection, parsed, "user"); err != nil {
		return
	}

	return cfg, nil
}
