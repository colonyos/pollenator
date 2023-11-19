package colonies

import (
	"fmt"
	"time"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/fs"
	"github.com/colonyos/pollinator/pkg/project"
	log "github.com/sirupsen/logrus"
)

func SyncDir(dir string, client *client.ColoniesClient, colonyID string, executorPrvKey string, proj *project.Project, keepLocal bool) error {
	fsClient, err := fs.CreateFSClient(client, colonyID, executorPrvKey)
	if err != nil {
		return err
	}

	label := "/pollinator/" + proj.ProjectID + dir
	syncPlans, err := fsClient.CalcSyncPlans("./cfs"+dir, label, keepLocal)
	if err != nil {
		return err
	}

	counter := 0
	for _, syncPlan := range syncPlans {
		if len(syncPlan.LocalMissing) == 0 && len(syncPlan.RemoteMissing) == 0 && len(syncPlan.Conflicts) == 0 {
			counter++
		}
	}

	if counter == len(syncPlans) {
		log.WithFields(log.Fields{"Label": label, "SyncDir": "./cfs" + dir}).Debug("Synchronizing, nothing to do, already synchronized")
		return nil
	}

	for _, syncPlan := range syncPlans {
		err = fsClient.ApplySyncPlan(colonyID, syncPlan)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateSrcSnapshot(client *client.ColoniesClient, colonyID string, executorPrvKey string, proj *project.Project) (string, error) {
	snapshotID := core.GenerateRandomID()
	snapshot, err := client.CreateSnapshot(colonyID, "/pollinator/"+proj.ProjectID+"/src", snapshotID, executorPrvKey)
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{"SnapshotID": snapshot.ID, "Label": snapshot.Label}).Debug("Creating snapshot")

	return snapshot.ID, nil
}

func CreateFuncSpec(colonyID string, project *project.Project, snapshotID string) *core.FunctionSpec {
	maxRetries := 3
	env := make(map[string]string)
	env["PROJECT_DIR"] = "/cfs/" + project.ProjectID

	args := make([]interface{}, 0)
	kwargsArgs := make([]interface{}, 0)
	kwargsArgs = append(kwargsArgs, "/cfs/"+project.ProjectID+"/src/"+project.Environment.SourceFile)

	kwargs := make(map[string]interface{}, 1)
	kwargs["cmd"] = project.Environment.Cmd
	kwargs["docker-image"] = project.Environment.DockerImage
	kwargs["rebuild-image"] = project.Environment.RebuildImage
	kwargs["args"] = kwargsArgs

	var snapshots []core.SnapshotMount
	snapshot1 := core.SnapshotMount{
		Label:       "/pollinator/" + project.ProjectID + "/src",
		SnapshotID:  snapshotID,
		Dir:         "/" + project.ProjectID + "/src",
		KeepFiles:   false,
		KeepSnaphot: false}

	snapshots = append(snapshots, snapshot1)
	var syncdirs []core.SyncDirMount
	result := core.SyncDirMount{
		Label:     "/pollinator/" + project.ProjectID + "/result",
		Dir:       "/" + project.ProjectID + "/result",
		KeepFiles: false,
		ConflictResolution: core.ConflictResolution{
			OnStart: core.OnStart{KeepLocal: false},
			OnClose: core.OnClose{KeepLocal: true}}}
	syncdirs = append(syncdirs, result)

	data := core.SyncDirMount{
		Label:     "/pollinator/" + project.ProjectID + "/data",
		Dir:       "/" + project.ProjectID + "/data",
		KeepFiles: true,
		ConflictResolution: core.ConflictResolution{
			OnStart: core.OnStart{KeepLocal: false},
			OnClose: core.OnClose{KeepLocal: false}}}
	syncdirs = append(syncdirs, data)

	maxWaitTime := -1
	maxExecTime := project.Conditions.Walltime - 1
	funcSpec := core.CreateFunctionSpec(
		"",
		"execute",
		args,
		kwargs,
		colonyID,
		[]string{},
		project.Conditions.ExecutorType,
		maxWaitTime,
		maxExecTime,
		maxRetries,
		env,
		[]string{"test_name2"},
		5,
		"test_label")

	funcSpec.Filesystem = core.Filesystem{SnapshotMounts: snapshots, SyncDirMounts: syncdirs, Mount: "/cfs"}

	funcSpec.Conditions.Nodes = project.Conditions.Nodes
	funcSpec.Conditions.CPU = project.Conditions.CPU
	funcSpec.Conditions.ProcessesPerNode = project.Conditions.ProcessesPerNode
	funcSpec.Conditions.Memory = project.Conditions.Memory
	funcSpec.Conditions.GPU = core.GPU{Name: project.Conditions.GPU.Name, Count: project.Conditions.GPU.Count}
	funcSpec.Conditions.WallTime = int64(project.Conditions.Walltime)

	return funcSpec
}

func Follow(client *client.ColoniesClient, process *core.Process, executorPrvKey string, count int) error {
	log.WithFields(log.Fields{"ProcessID": process.ID}).Debug("Printing logs from process")
	var lastTimestamp int64
	lastTimestamp = 0
	for {
		logs, err := client.GetLogsByProcessIDSince(process.ID, count, lastTimestamp, executorPrvKey)
		if err != nil {
			return err
		}

		process, err := client.GetProcess(process.ID, executorPrvKey)
		if err != nil {
			return err
		}

		if len(logs) == 0 {
			time.Sleep(500 * time.Millisecond)
			if process.State == core.SUCCESS {
				log.WithFields(log.Fields{"ProcessID": process.ID}).Info("Process finished successfully")
				return nil
			}
			if process.State == core.FAILED {
				fmt.Println()
				log.WithFields(log.Fields{"ProcessID": process.ID}).Error("Process failed")
				return err
			}
			continue
		} else {
			for _, log := range logs {
				fmt.Print(log.Message)
			}
			lastTimestamp = logs[len(logs)-1].Timestamp
		}

	}
}
