
package events

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	
)

const InitProcNsDir = "/proc/1/ns"

// InitNamespacesEvent collect the init process namespaces and create event from them.
func InitNamespacesEvent() trace.Event {
	initNamespacesDef := Definitions.Get(InitNamespaces)
	initNamespacesArgs := getInitNamespaceArguments()
	initNamespacesEvent := trace.Event{
		Timestamp:   int(time.Now().UnixNano()),
		ProcessName: "tracee-ebpf",
		EventID:     int(InitNamespaces),
		EventName:   initNamespacesDef.Name,
		ArgsNum:     len(initNamespacesArgs),
		Args:        initNamespacesArgs,
	}
	return initNamespacesEvent
}

// getInitNamespaceArguments Fetch the namespaces of the init process and parse them into event arguments.
func getInitNamespaceArguments() []trace.Argument {
	initNamespaces := fetchInitNamespaces()
	eventDefinition := Definitions.Get(InitNamespaces)
	initNamespacesArgs := make([]trace.Argument, len(eventDefinition.Params))
	for i, arg := range initNamespacesArgs {
		arg.ArgMeta = eventDefinition.Params[i]
		arg.Value = initNamespaces[arg.Name]
		initNamespacesArgs[i] = arg
	}
	return initNamespacesArgs
}

// fetchInitNamespaces fetch the namespaces values from the /proc/1/ns directory
func fetchInitNamespaces() map[string]uint32 {
	initNamespacesMap := make(map[string]uint32)
	namespaceValueReg := regexp.MustCompile(":[[[:digit:]]*]")
	namespacesLinks, _ := os.ReadDir(InitProcNsDir)
	for _, namespaceLink := range namespacesLinks {
		linkString, _ := os.Readlink(filepath.Join(InitProcNsDir, namespaceLink.Name()))
		trim := strings.Trim(namespaceValueReg.FindString(linkString), "[]:")
		namespaceNumber, _ := strconv.ParseUint(trim, 10, 32)
		initNamespacesMap[namespaceLink.Name()] = uint32(namespaceNumber)
	}
	return initNamespacesMap
}

// ExistingContainersEvents returns a list of events for each existing container
func ExistingContainersEvents(containers *containers.Containers, enrich bool) []trace.Event {
	var events []trace.Event
	def := Definitions.Get(ExistingContainer)
	for id, info := range containers.GetContainers() {
		container := runtime.ContainerMetadata{}
		if enrich {
			container, _ = containers.EnrichCgroupInfo(uint64(id))
		}
		args := []trace.Argument{
			{ArgMeta: def.Params[0], Value: info.Runtime.String()},
			{ArgMeta: def.Params[1], Value: info.Container.ContainerId},
			{ArgMeta: def.Params[2], Value: info.Ctime.UnixNano()},
			{ArgMeta: def.Params[3], Value: container.Image},
			{ArgMeta: def.Params[4], Value: container.Name},
			{ArgMeta: def.Params[5], Value: container.Pod.Name},
			{ArgMeta: def.Params[6], Value: container.Pod.Namespace},
			{ArgMeta: def.Params[7], Value: container.Pod.UID},
			{ArgMeta: def.Params[8], Value: container.Pod.Sandbox},
		}
		existingContainerEvent := trace.Event{
			Timestamp:   int(time.Now().UnixNano()),
			ProcessName: "travel website",
			EventID:     int(ExistingContainer),
			EventName:   def.Name,
			ArgsNum:     len(args),
			Args:        args,
		}
		events = append(events, existingContainerEvent)
	}
	return events
}
