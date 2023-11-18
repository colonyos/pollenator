package project

import (
	"errors"
	"os"

	"github.com/colonyos/colonies/pkg/core"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Conditions struct {
	ExecutorType     string `yaml:"executorType"`
	Nodes            int    `yaml:"nodes"`
	ProcessesPerNode int    `yaml:"processesPerNode"`
	CPU              string `yaml:"cpu"`
	Memory           string `yaml:"mem"`
	Walltime         int    `yaml:"walltime"`
	GPU              GPU    `yaml:"gpu"`
}

type GPU struct {
	Count int    `yaml:"count"`
	Name  string `yaml:"name"`
}

type Environment struct {
	DockerImage  string `yaml:"docker"`
	RebuildImage bool   `yaml:"rebuildImage"`
	Cmd          string `yaml:"cmd"`
	SourceFile   string `yaml:"source"`
}

type Project struct {
	ProjectID   string      `yaml:"projectid"`
	Conditions  Conditions  `yaml:"conditions"`
	Environment Environment `yaml:"environment"`
}

func GenerateProjectConfig(executorType string) error {
	projectFile := "./project.yaml"
	_, err := os.Stat(projectFile)
	if err == nil {
		return errors.New(projectFile + " already exists")
	}

	log.WithFields(log.Fields{
		"Filename": projectFile}).
		Info("Generating")

	project := &Project{}
	cond := Conditions{
		ExecutorType:     executorType,
		Nodes:            1,
		ProcessesPerNode: 1,
		CPU:              "1000m",
		Memory:           "1000M",
		Walltime:         600}

	env := Environment{
		DockerImage:  "python:3.12-rc-bookworm",
		RebuildImage: false,
		Cmd:          "python3",
		SourceFile:   "src/main.py",
	}

	projectID := core.GenerateRandomID()

	project.ProjectID = projectID
	project.Conditions = cond
	project.Environment = env

	file, err := os.Create(projectFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func GenerateProjectData() error {
	dataFilename := "./cfs/data/hello.txt"
	_, err := os.Stat(dataFilename)
	if err == nil {
		return errors.New(dataFilename + " already exists")
	}

	log.WithFields(log.Fields{
		"Filename": dataFilename}).
		Info("Generating")

	data := `Hello world!`

	dataFile, err := os.Create(dataFilename)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	_, err = dataFile.Write([]byte(data))
	if err != nil {
		return err
	}

	src := `import os

projdir = str(os.environ.get("PROJECT_DIR"))
processid = os.environ.get("COLONIES_PROCESS_ID")

file = open(projdir + "/data/hello.txt", 'r')
contents = file.read()
print(contents)

result_dir = projdir + "/result/"
os.makedirs(result_dir, exist_ok=True)

file = open(result_dir + "/result.txt", "w")
file.write("Hello, World!")
file.close()`
	srcFilename := "./cfs/src/main.py"
	_, err = os.Stat(srcFilename)
	if err == nil {
		return errors.New(srcFilename + " already exists")
	}

	log.WithFields(log.Fields{
		"Filename": srcFilename}).
		Info("Generating")

	sourceFile, err := os.Create(srcFilename)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	_, err = sourceFile.Write([]byte(src))
	return err
}
