package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdNode = &cobra.Command{
		Use:   "node",
		Short: "manage a opensvc cluster node",
	}

	cmdNodeCapabilities = &cobra.Command{
		Use:     "capabilities",
		Short:   "scan and list what the node is capable of",
		Aliases: []string{"capa", "caps", "cap", "ca"},
	}
	cmdNodeCompliance = &cobra.Command{
		Use:     "compliance",
		Short:   "node configuration manager commands",
		Aliases: []string{"compli", "comp", "com", "co"},
	}
	cmdNodeComplianceAttach = &cobra.Command{
		Use:     "attach",
		Short:   "attach modulesets and rulesets to the node",
		Aliases: []string{"attac", "atta", "att", "at"},
	}
	cmdNodeComplianceDetach = &cobra.Command{
		Use:     "detach",
		Short:   "detach modulesets and rulesets from the node",
		Aliases: []string{"detac", "deta", "det", "de"},
	}
	cmdNodeComplianceList = &cobra.Command{
		Use:     "list",
		Short:   "list modules, modulesets and rulesets available",
		Aliases: []string{"lis", "li", "ls", "l"},
	}
	cmdNodeComplianceShow = &cobra.Command{
		Use:     "show",
		Short:   "show modules, modulesets, rulesets, modules, attachments",
		Aliases: []string{"sho", "sh", "s"},
	}
	cmdNodePrint = &cobra.Command{
		Use:     "print",
		Short:   "print node discover information",
		Aliases: []string{"prin", "pri", "pr"},
	}
	cmdNodePush = &cobra.Command{
		Use:   "push",
		Short: "push node discover information to the collector",
	}
	cmdNodeRelay = &cobra.Command{
		Use:   "relay",
		Short: "relay subsystem commands",
	}
	cmdNodeEdit     = newCmdNodeEdit()
	cmdNodeValidate = newCmdNodeValidate()
)

func init() {
	root.AddCommand(cmdNode)
	cmdNode.AddCommand(cmdNodeCapabilities)
	cmdNodeCapabilities.AddCommand(
		newCmdNodeCapabilitiesList(),
		newCmdNodeCapabilitiesScan(),
	)
	cmdNode.AddCommand(cmdNodeCompliance)
	cmdNodeCompliance.AddCommand(
		cmdNodeComplianceAttach,
		cmdNodeComplianceDetach,
		cmdNodeComplianceShow,
		cmdNodeComplianceList,
		newCmdNodeComplianceEnv(),
		newCmdNodeComplianceAuto(),
		newCmdNodeComplianceCheck(),
		newCmdNodeComplianceFix(),
		newCmdNodeComplianceFixable(),
	)
	cmdNodeComplianceAttach.AddCommand(
		newCmdNodeComplianceAttachModuleset(),
		newCmdNodeComplianceAttachRuleset(),
	)
	cmdNodeComplianceDetach.AddCommand(
		newCmdNodeComplianceDetachModuleset(),
		newCmdNodeComplianceDetachRuleset(),
	)
	cmdNodeComplianceShow.AddCommand(
		newCmdNodeComplianceShowRuleset(),
		newCmdNodeComplianceShowModuleset(),
	)
	cmdNodeComplianceList.AddCommand(
		newCmdNodeComplianceListModules(),
		newCmdNodeComplianceListModuleset(),
		newCmdNodeComplianceListRuleset(),
	)
	cmdNodeEdit.AddCommand(
		newCmdNodeEditConfig(),
	)
	cmdNode.AddCommand(
		cmdNodeEdit,
		cmdNodePrint,
		cmdNodePush,
		cmdNodeRelay,
		cmdNodeValidate,
		newCmdNodeAbort(),
		newCmdNodeChecks(),
		newCmdNodeClear(),
		newCmdNodeDoc(),
		newCmdNodeDelete(),
		newCmdNodeDrain(),
		newCmdNodeDrivers(),
		newCmdNodeLogs(),
		newCmdNodeLs(),
		newCmdNodePRKey(),
		newCmdNodeFreeze(),
		newCmdNodeGet(),
		newCmdNodeEvents(),
		newCmdNodeEval(),
		newCmdNodeRegister(),
		newCmdNodeSet(),
		newCmdNodeSysreport(),
		newCmdNodeUnfreeze(),
		newCmdNodeUnset(),
	)
	cmdNodePrint.AddCommand(
		newCmdNodePrintConfig(),
		newCmdNodePrintSchedule(),
	)
	cmdNodePush.AddCommand(
		newCmdNodePushAsset(),
		newCmdNodePushDisks(),
		newCmdNodePushPatch(),
		newCmdNodePushPkg(),
	)
	cmdNodeValidate.AddCommand(
		newCmdNodeValidateConfig(),
	)
	cmdNodeRelay.AddCommand(
		newCmdNodeRelayStatus(),
	)

}
