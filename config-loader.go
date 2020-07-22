package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	serverConfig map[string]string
)

func lgServerConfigListLoad(configURL string, svLoc string) {
	configLogger := log.New(os.Stdout, "Config Loader: ", log.LstdFlags)
	serverConfig = make(map[string]string)
	serverConfig["whois"] = "disabled"
	serverConfig["nslookup"] = "disabled"
	serverConfig["ping"] = "disabled"
	serverConfig["icmp"] = "disabled"
	serverConfig["tracert"] = "disabled"
	serverConfig["webcontrol"] = "disabled"
	serverConfig["tcp"] = "disabled"
	serverConfig["IPv4"] = "disabled"
	serverConfig["IPv6"] = "disabled"
	serverConfig["curl"] = "disabled"

	if configURL == "" {
		configLogger.Fatalln("Config url is not defined. please define with --config-url")
	} else {
		if strings.Contains(configURL, "--svloc") {
			configLogger.Fatalln("Config url is is empty. Please define with --config-url https://yoururl/server.json")
		}
	}
	if svLoc == "" {
		configLogger.Println("server Location is not defined. Please define with --svloc")
	}

	if svLoc == "" || configURL == "" {
		configLogger.Println("\033[1;31mAll settings Enabled.\033[0m")
	} else {
		resp, err := http.Get(configURL)
		if err != nil {
			configLogger.Fatalln(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			configLogger.Fatalln("Config server response is ok. Response : " + resp.Status)
		}
		RemoteConfigJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			configLogger.Fatalln(err)
		}

		currentServerName := strings.Split(svLoc, ".")

		for i := -1; i < len(currentServerName); i++ {
			var temp string
			for t := 0; t <= i; t++ {
				temp = temp + ".Servers." + currentServerName[t]
			}
			if i == len(currentServerName)-1 {
				temp2 := strings.Replace(temp, ".", "", 1)
				serverURLexpexted := gjson.Get(string(RemoteConfigJSON), temp2+".Url").String()
				serverListUnExpexted := gjson.Get(string(RemoteConfigJSON), temp2+".Servers").String()
				if serverURLexpexted != "" && serverListUnExpexted == "" {
					serverConfig["ThisServerURL"] = serverURLexpexted
				} else {
					configLogger.Fatalln("Given server location is not a server.")
				}

			}
			temp = temp + ".ServerConfig"
			temp = strings.Replace(temp, ".", "", 1)
			//configLogger.Println("\nCurrent loc " + temp)
			if temp != "" {
				for k := range serverConfig {
					if gjson.Get(string(RemoteConfigJSON), temp+"."+k).String() != "" {
						serverConfig[k] = gjson.Get(string(RemoteConfigJSON), temp+"."+k).String()
						//fmt.Printf("%s %s\n", k, serverConfig[k])
					}
				}
			}

		}
		configLogger.Println("Setting loaded for " + svLoc)
		for k, v := range serverConfig {
			configLogger.Println(k, v)
		}
	}

	// Get frontend server address from config url
	parsedConfigURL, err := url.Parse(configURL)
	if err != nil {
		panic(err)
	}
	if *allowedReferersConfig == "" {
		configLogger.Println("Referer Domain is not given, you can set with --referers. System is only allows incoming requests from ", parsedConfigURL.Host)
		allowedReferers = []string{parsedConfigURL.Host} //parsedConfigURL.Host
	} else {
		allowedReferers = strings.Split(*allowedReferersConfig, ",")
		configLogger.Println("Allowed referer domains ", allowedReferers)
	}
}
