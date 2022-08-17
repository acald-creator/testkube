package tests

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	apiv1 "github.com/kubeshop/testkube/pkg/api/v1/client"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/ui"
)

const WatchInterval = 2 * time.Second

func NewRunTestCmd() *cobra.Command {
	var (
		name                     string
		image                    string
		watchEnabled             bool
		binaryArgs               []string
		variables                map[string]string
		secretVariables          map[string]string
		variablesFile            string
		downloadArtifactsEnabled bool
		downloadDir              string
		envs                     map[string]string
		secretEnvs               map[string]string
		selectors                []string
		concurrencyLevel         int
		httpProxy, httpsProxy    string
		executionLabels          map[string]string
	)

	cmd := &cobra.Command{
		Use:     "test <testName>",
		Aliases: []string{"t"},
		Short:   "Starts new test",
		Long:    `Starts new test based on Test Custom Resource name, returns results to console`,
		Run: func(cmd *cobra.Command, args []string) {

			paramsFileContent := ""
			if variablesFile != "" {
				b, err := os.ReadFile(variablesFile)
				ui.ExitOnError("reading variables file", err)
				paramsFileContent = string(b)
			}

			envs, err := cmd.Flags().GetStringToString("env")
			ui.WarnOnError("getting envs", err)

			variables, err := common.CreateVariables(cmd)
			ui.WarnOnError("getting variables", err)

			executorArgs, err := prepareExecutorArgs(binaryArgs)
			ui.ExitOnError("getting args", err)

			var executions []testkube.Execution
			client, namespace := common.GetClient(cmd)
			options := apiv1.ExecuteTestOptions{
				ExecutionVariables:            variables,
				ExecutionVariablesFileContent: paramsFileContent,
				ExecutionLabels:               executionLabels,
				Args:                          executorArgs,
				SecretEnvs:                    secretEnvs,
				HTTPProxy:                     httpProxy,
				HTTPSProxy:                    httpsProxy,
				Envs:                          envs,
				Image:                         image,
			}

			switch {
			case len(args) > 0:
				testName := args[0]
				namespacedName := fmt.Sprintf("%s/%s", namespace, testName)

				_, err = client.GetTest(testName)
				if err != nil {
					ui.UseStderr()
					ui.Errf("Can't get test with name '%s'. Test does not exist in namespace '%s'", testName, namespace)
					ui.Debug(err.Error())
					os.Exit(1)
				}

				execution, err := client.ExecuteTest(testName, name, options)
				ui.ExitOnError("starting test execution "+namespacedName, err)
				executions = append(executions, execution)
			case len(selectors) != 0:
				selector := strings.Join(selectors, ",")
				executions, err = client.ExecuteTests(selector, concurrencyLevel, options)
				ui.ExitOnError("starting test executions "+selector, err)
			default:
				ui.Failf("Pass Test name or labels to run by labels ")
			}

			var hasErrors bool
			for _, execution := range executions {
				printExecutionDetails(execution)

				if execution.ExecutionResult != nil && execution.ExecutionResult.ErrorMessage != "" {
					hasErrors = true
				}

				if execution.Id != "" {
					if watchEnabled && len(args) > 0 {
						watchLogs(execution.Id, client)
					}

					execution, err = client.GetExecution(execution.Id)
					ui.ExitOnError("getting recent execution data id:"+execution.Id, err)
				}

				uiPrintStatus(execution)

				if execution.Id != "" {
					if downloadArtifactsEnabled {
						DownloadArtifacts(execution.Id, downloadDir, client)
					}

					uiShellWatchExecution(execution.Id)
				}

				uiShellGetExecution(execution.Id)
			}

			if hasErrors {
				ui.ExitOnError("executions contain failed on errors")
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "execution name, if empty will be autogenerated")
	cmd.Flags().StringVarP(&image, "image", "", "", "execution variable passed to executor")
	cmd.Flags().StringVarP(&variablesFile, "variables-file", "", "", "variables file path, e.g. postman env file - will be passed to executor if supported")
	cmd.Flags().StringToStringVarP(&variables, "variable", "v", map[string]string{}, "execution variable passed to executor")
	cmd.Flags().StringToStringVarP(&secretVariables, "secret-variable", "s", map[string]string{}, "execution secret variable passed to executor")
	cmd.Flags().StringArrayVarP(&binaryArgs, "args", "", []string{}, "executor binary additional arguments")
	cmd.Flags().BoolVarP(&watchEnabled, "watch", "f", false, "watch for changes after start")
	cmd.Flags().StringVar(&downloadDir, "download-dir", "artifacts", "download dir")
	cmd.Flags().BoolVarP(&downloadArtifactsEnabled, "download-artifacts", "d", false, "downlaod artifacts automatically")
	cmd.Flags().StringToStringVarP(&envs, "env", "", map[string]string{}, "envs in a form of name1=val1 passed to executor")
	cmd.Flags().StringToStringVarP(&secretEnvs, "secret", "", map[string]string{}, "secret envs in a form of secret_name1=secret_key1 passed to executor")
	cmd.Flags().StringSliceVarP(&selectors, "label", "l", nil, "label key value pair: --label key1=value1")
	cmd.Flags().IntVar(&concurrencyLevel, "concurrency", 10, "concurrency level for multiple test execution")
	cmd.Flags().StringVar(&httpProxy, "http-proxy", "", "http proxy for executor containers")
	cmd.Flags().StringVar(&httpsProxy, "https-proxy", "", "https proxy for executor containers")
	cmd.Flags().StringToStringVarP(&executionLabels, "execution-label", "", nil, "execution-label key value pair: --execution-label key1=value1")

	return cmd
}

func uiPrintStatus(execution testkube.Execution) {
	result := execution.ExecutionResult

	if result == nil {
		return
	}

	ui.NL()

	switch true {
	case result.IsQueued():
		ui.Warn("Test queued for execution")

	case result.IsRunning():
		ui.Warn("Test execution started")

	case result.IsPassed():
		ui.Info(result.Output)
		duration := execution.EndTime.Sub(execution.StartTime)
		ui.Success("Test execution completed with success in " + duration.String())

	case result.IsFailed():
		ui.UseStderr()
		ui.Warn("Test test execution failed:\n")
		ui.Errf(result.ErrorMessage)
		ui.Info(result.Output)
		os.Exit(1)
	}

	ui.NL()
}

func uiShellGetExecution(id string) {
	ui.ShellCommand(
		"Use following command to get test execution details",
		"kubectl testkube get execution "+id,
	)

	ui.NL()
}

func uiShellWatchExecution(id string) {
	ui.ShellCommand(
		"Watch test execution until complete",
		"kubectl testkube watch execution "+id,
	)

	ui.NL()
}
