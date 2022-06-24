package cmd

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/commands"
)

var (
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "manage a opensvc cluster node",
	}

	nodeComplianceCmd = &cobra.Command{
		Use:     "compliance",
		Short:   "node configuration expectations analysis and application",
		Aliases: []string{"compli", "comp", "com", "co"},
	}
	nodeComplianceAttachCmd = &cobra.Command{
		Use:     "attach",
		Short:   "attach modulesets and rulesets to the node.",
		Aliases: []string{"attac", "atta", "att", "at"},
	}
	nodeComplianceDetachCmd = &cobra.Command{
		Use:     "detach",
		Short:   "detach modulesets and rulesets from the node.",
		Aliases: []string{"detac", "deta", "det", "de"},
	}
	nodeComplianceListCmd = &cobra.Command{
		Use:     "list",
		Short:   "list modules, modulesets and rulesets available",
		Aliases: []string{"lis", "li", "ls", "l"},
	}
	nodeComplianceShowCmd = &cobra.Command{
		Use:     "show",
		Short:   "show states: current moduleset and ruleset attachments, modules last check",
		Aliases: []string{"sho", "sh", "s"},
	}
	nodePrintCmd = &cobra.Command{
		Use:     "print",
		Short:   "print node",
		Aliases: []string{"prin", "pri", "pr"},
	}
	nodePushCmd = &cobra.Command{
		Use:   "push",
		Short: "data pushing commands",
	}
	nodeScanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scan node",
	}
	nodeEditCmd = &cobra.Command{
		Use:     "edit",
		Short:   "edition command group",
		Aliases: []string{"edi", "ed"},
	}
	nodeValidateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "validation command group",
		Aliases: []string{"validat", "valida", "valid", "val"},
	}

	cmdNodeChecks                    commands.CmdNodeChecks
	cmdNodeComplianceAttachModuleset commands.CmdNodeComplianceAttachModuleset
	cmdNodeComplianceDetachModuleset commands.CmdNodeComplianceDetachModuleset
	cmdNodeComplianceAttachRuleset   commands.CmdNodeComplianceAttachRuleset
	cmdNodeComplianceDetachRuleset   commands.CmdNodeComplianceDetachRuleset
	cmdNodeComplianceAuto            commands.CmdNodeComplianceAuto
	cmdNodeComplianceCheck           commands.CmdNodeComplianceCheck
	cmdNodeComplianceFix             commands.CmdNodeComplianceFix
	cmdNodeComplianceFixable         commands.CmdNodeComplianceFixable
	cmdNodeComplianceShowRuleset     commands.CmdNodeComplianceShowRuleset
	cmdNodeComplianceShowModuleset   commands.CmdNodeComplianceShowModuleset
	cmdNodeComplianceListModules     commands.CmdNodeComplianceListModules
	cmdNodeComplianceListModuleset   commands.CmdNodeComplianceListModuleset
	cmdNodeComplianceListRuleset     commands.CmdNodeComplianceListRuleset
	cmdNodeComplianceEnv             commands.CmdNodeComplianceEnv
	cmdNodeDoc                       commands.NodeDoc
	cmdNodeDelete                    commands.NodeDelete
	cmdNodeDrivers                   commands.NodeDrivers
	cmdNodeEditConfig                commands.NodeEditConfig
	cmdNodeLogs                      commands.NodeLogs
	cmdNodeLs                        commands.NodeLs
	cmdNodeGet                       commands.NodeGet
	cmdNodeEval                      commands.NodeEval
	cmdNodePrintCapabilities         commands.NodePrintCapabilities
	cmdNodePrintConfig               commands.NodePrintConfig
	cmdNodePushAsset                 commands.NodePushAsset
	cmdNodePushDisks                 commands.NodePushDisks
	cmdNodePushPatch                 commands.NodePushPatch
	cmdNodePushPkg                   commands.NodePushPkg
	cmdNodeRegister                  commands.CmdNodeRegister
	cmdNodeScanCapabilities          commands.NodeScanCapabilities
	cmdNodeSet                       commands.NodeSet
	cmdNodeSysreport                 commands.CmdNodeSysreport
	cmdNodeUnset                     commands.NodeUnset
	cmdNodeValidateConfig            commands.NodeValidateConfig
)

func init() {
	root.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeComplianceCmd)
	nodeComplianceCmd.AddCommand(nodeComplianceAttachCmd)
	nodeComplianceCmd.AddCommand(nodeComplianceDetachCmd)
	nodeComplianceCmd.AddCommand(nodeComplianceShowCmd)
	nodeComplianceCmd.AddCommand(nodeComplianceListCmd)
	nodeCmd.AddCommand(nodeEditCmd)
	nodeCmd.AddCommand(nodePrintCmd)
	nodeCmd.AddCommand(nodePushCmd)
	nodeCmd.AddCommand(nodeScanCmd)
	nodeCmd.AddCommand(nodeValidateCmd)

	cmdNodeChecks.Init(nodeCmd)
	cmdNodeComplianceEnv.Init(nodeComplianceCmd)
	cmdNodeComplianceAttachModuleset.Init(nodeComplianceAttachCmd)
	cmdNodeComplianceAttachRuleset.Init(nodeComplianceAttachCmd)
	cmdNodeComplianceDetachModuleset.Init(nodeComplianceDetachCmd)
	cmdNodeComplianceDetachRuleset.Init(nodeComplianceDetachCmd)
	cmdNodeComplianceAuto.Init(nodeComplianceCmd)
	cmdNodeComplianceCheck.Init(nodeComplianceCmd)
	cmdNodeComplianceFix.Init(nodeComplianceCmd)
	cmdNodeComplianceFixable.Init(nodeComplianceCmd)
	cmdNodeComplianceShowRuleset.Init(nodeComplianceShowCmd)
	cmdNodeComplianceShowModuleset.Init(nodeComplianceShowCmd)
	cmdNodeComplianceListModules.Init(nodeComplianceListCmd)
	cmdNodeComplianceListModuleset.Init(nodeComplianceListCmd)
	cmdNodeComplianceListRuleset.Init(nodeComplianceListCmd)
	cmdNodeDoc.Init(nodeCmd)
	cmdNodeDelete.Init(nodeCmd)
	cmdNodeDrivers.Init(nodeCmd)
	cmdNodeEditConfig.Init(nodeEditCmd)
	cmdNodeLogs.Init(nodeCmd)
	cmdNodeLs.Init(nodeCmd)
	cmdNodeGet.Init(nodeCmd)
	cmdNodeEval.Init(nodeCmd)
	cmdNodePrintCapabilities.Init(nodePrintCmd)
	cmdNodePrintConfig.Init(nodePrintCmd)
	cmdNodePushAsset.Init(nodePushCmd)
	cmdNodePushAsset.InitAlt(nodeCmd)
	cmdNodePushDisks.Init(nodePushCmd)
	cmdNodePushDisks.InitAlt(nodeCmd)
	cmdNodePushPatch.Init(nodePushCmd)
	cmdNodePushPatch.InitAlt(nodeCmd)
	cmdNodePushPkg.Init(nodePushCmd)
	cmdNodePushPkg.InitAlt(nodeCmd)
	cmdNodeRegister.Init(nodeCmd)
	cmdNodeScanCapabilities.Init(nodeScanCmd)
	cmdNodeSet.Init(nodeCmd)
	cmdNodeSysreport.Init(nodeCmd)
	cmdNodeUnset.Init(nodeCmd)
	cmdNodeValidateConfig.Init(nodeValidateCmd)

}
