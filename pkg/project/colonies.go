package project

import "github.com/colonyos/colonies/pkg/core"

func CreateFuncSpec(colonyID string, project *Project, snapshotID string) *core.FunctionSpec {
	maxRetries := 3
	env := make(map[string]string)

	args := make([]interface{}, 1)
	kwargsArgs := make([]interface{}, 0)
	kwargsArgs = append(kwargsArgs, "/cfs/"+project.ProjectID+"/src"+project.Environment.SourceFile)

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
	maxExecTime := project.Conditions.Walltime + 60
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
	funcSpec.Conditions.WallTime = funcSpec.Conditions.WallTime

	return funcSpec
}
