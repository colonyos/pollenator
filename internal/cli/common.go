package cli

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/pollinator/pkg/build"
	"github.com/colonyos/pollinator/pkg/colonies"
	"github.com/colonyos/pollinator/pkg/project"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func checkIfDirExists(dirPath string) error {
	fileInfo, err := os.Stat(dirPath)
	if err == nil {
		if fileInfo.IsDir() {
			return errors.New(dirPath + " already exists")
		}
	}
	return nil
}

func checkIfDirIsEmpty(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	return errors.New("Current directory is not empty, try create a new direcory and retry")
}

func parseEnv() {
	var err error
	ColoniesServerHostEnv := os.Getenv("COLONIES_SERVER_HOST")
	if ColoniesServerHostEnv != "" {
		ColoniesServerHost = ColoniesServerHostEnv
	}

	ColoniesServerPortEnvStr := os.Getenv("COLONIES_SERVER_PORT")
	if ColoniesServerPortEnvStr != "" {
		ColoniesServerPort, err = strconv.Atoi(ColoniesServerPortEnvStr)
		CheckError(err)
	}

	ColoniesTLSEnv := os.Getenv("COLONIES_SERVER_TLS")
	if ColoniesTLSEnv == "true" {
		ColoniesUseTLS = true
		ColoniesInsecure = false
	} else if ColoniesTLSEnv == "false" {
		ColoniesUseTLS = false
		ColoniesInsecure = true
	}

	VerboseEnv := os.Getenv("COLONIES_VERBOSE")
	if VerboseEnv == "true" {
		Verbose = true
	} else if VerboseEnv == "false" {
		Verbose = false
	}

	if ColonyName == "" {
		ColonyName = os.Getenv("COLONIES_COLONY_NAME")
	}
	if ColonyName == "" {
		CheckError(errors.New("Unknown Colony name"))
	}

	if PrvKey == "" {
		PrvKey = os.Getenv("COLONIES_PRVKEY")
	}
	if PrvKey == "" {
		CheckError(errors.New("Unkown private key"))
	}

	DashboardURL = os.Getenv("COLONYOS_DASHBOARD_URL")
}

func CheckError(err error) {
	if err != nil {
		log.WithFields(log.Fields{"Error": err, "BuildVersion": build.BuildVersion, "BuildTime": build.BuildTime}).Error(err.Error())
		os.Exit(-1)
	}
}

func SyncAndGenerateFuncSpec(client *client.ColoniesClient) (*core.FunctionSpec, *project.Project) {
	projectFile := "./project.yaml"
	projectData, err := ioutil.ReadFile(projectFile)
	CheckError(err)

	proj := &project.Project{}
	err = yaml.Unmarshal([]byte(projectData), &proj)
	CheckError(err)

	// Sync all directories
	err = colonies.SyncDir("/src", client, ColonyName, PrvKey, proj, true)
	CheckError(err)
	err = colonies.SyncDir("/data", client, ColonyName, PrvKey, proj, true)
	CheckError(err)

	snapshotID, err := colonies.CreateSrcSnapshot(client, ColonyName, PrvKey, proj)
	CheckError(err)

	log.Debug("Generating function spec")
	funcSpec := colonies.CreateFuncSpec(ColonyName, proj, snapshotID)
	CheckError(err)

	return funcSpec, proj
}

func formatTimestamp(timestamp string) string {
	return strings.Replace(timestamp, "T", " ", 1)
}
