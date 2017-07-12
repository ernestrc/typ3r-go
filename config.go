package typ3r

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path"

	ini "gopkg.in/ini.v1"
)

const confFile = ".typ3rrc"
const tokenKey = "token"
const serverURLKey = "url"
const userKey = "user"

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

func gatherConfig(file *ini.File) error {
	const n = 3
	var err error
	var global *ini.Section
	var line []byte

	if global, err = file.GetSection(""); err != nil {
		return err
	}

	k := [n]string{serverURLKey, tokenKey, userKey}
	d := [n]string{defaultURL, "", "token"}
	q := [n]string{
		"What typ3r server do you want to use?",
		"What's your session token? If you don't have one create one by visiting https://typ3r.com/api/index.php/createtoken with your browser",
		"What's your user name?",
	}

	reader := bufio.NewReader(os.Stdin)

	for i, question := range q {
		fmt.Print(question)
		if d[i] != "" {
			fmt.Printf(" (default=%s): ", d[i])
		} else {
			fmt.Print(": ")
		}
		if line, _, err = reader.ReadLine(); err != nil {
			return err
		}
		a := string(line)
		if a != "" {
			d[i] = a
		}
		if _, err = global.NewKey(k[i], d[i]); err != nil {
			return err
		}
	}

	return nil
}

func newConfigFile() (file *ini.File, err error) {
	file = ini.Empty()

	if err = gatherConfig(file); err != nil {
		return
	}

	if err = file.SaveTo(confPath); err != nil {
		return
	}

	return file, nil
}

// LoadConfig loads a typ3r's client user configuration
func LoadConfig() (cfg *Config, err error) {
	var parsed *ini.File
	var created = false

	// create if config doesn't exist
	if _, err = os.Stat(confPath); err != nil {
		created = true
		if parsed, err = newConfigFile(); err != nil {
			return
		}
	} else if parsed, err = ini.Load(confPath); err != nil {
		return
	}

	cfg = new(Config)

	cfg.serverURL = getOption(rootSection, parsed, serverURLKey, defaultURL)

	if (*cfg).token, err = getConfig(rootSection, parsed, tokenKey); err != nil {
		goto clean
	}

	if cfg.user, err = getConfig(rootSection, parsed, userKey); err != nil {
		goto clean
	}

	return cfg, nil

clean:
	if created {
		os.Remove(confPath)
	}
	return nil, err
}
