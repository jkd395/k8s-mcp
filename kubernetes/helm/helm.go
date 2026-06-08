package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-mcp/kubernetes/client"
	"k8s-mcp/kubernetes/output"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func getActionConfig(ns string) (*action.Configuration, error) {
	cfg := &action.Configuration{}
	restConfig, err := client.GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get rest config: %v", err)
	}
	rg := &restClientGetter{restConfig: restConfig}
	if err := cfg.Init(rg, ns, "secret", log.Printf); err != nil {
		return nil, fmt.Errorf("failed to init helm config: %v", err)
	}
	return cfg, nil
}

func getSettings() *cli.EnvSettings {
	return cli.New()
}

type restClientGetter struct {
	restConfig *rest.Config
}

func (r *restClientGetter) ToRESTConfig() (*rest.Config, error) {
	return r.restConfig, nil
}

func (r *restClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	d, err := discovery.NewDiscoveryClientForConfig(r.restConfig)
	if err != nil {
		return nil, err
	}
	return memory.NewMemCacheClient(d), nil
}

func (r *restClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	d, err := r.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	return restmapper.NewDeferredDiscoveryRESTMapper(d), nil
}

func (r *restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return &simpleClientConfig{restConfig: r.restConfig}
}

type simpleClientConfig struct {
	restConfig *rest.Config
}

func (s *simpleClientConfig) RawConfig() (clientcmdapi.Config, error) {
	return clientcmdapi.Config{}, nil
}

func (s *simpleClientConfig) ClientConfig() (*rest.Config, error) {
	return s.restConfig, nil
}

func (s *simpleClientConfig) Namespace() (string, bool, error) {
	return "default", false, nil
}

func (s *simpleClientConfig) ConfigAccess() clientcmd.ConfigAccess {
	return nil
}

type releaseData struct {
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Version    int    `json:"version,omitempty"`
	Status     string `json:"status,omitempty"`
	Chart      string `json:"chart,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
	Updated    string `json:"updated,omitempty"`
}

type releaseHistoryData struct {
	Revision    int    `json:"revision"`
	Status      string `json:"status"`
	Chart       string `json:"chart"`
	AppVersion  string `json:"appVersion"`
	Description string `json:"description"`
	Updated     string `json:"updated"`
}

func releaseSummary(r *release.Release) releaseData {
	return releaseData{
		Name:       r.Name,
		Namespace:  r.Namespace,
		Version:    r.Version,
		Status:     r.Info.Status.String(),
		Chart:      r.Chart.Name(),
		AppVersion: r.Chart.AppVersion(),
		Updated:    r.Info.LastDeployed.Local().Format(time.RFC3339),
	}
}

func buildValues(valuesStr, setStr string) (map[string]interface{}, error) {
	merged := make(map[string]interface{})

	if valuesStr != "" {
		vals, err := chartutil.ReadValues([]byte(valuesStr))
		if err != nil {
			return nil, fmt.Errorf("failed to parse values: %v", err)
		}
		for k, v := range vals {
			merged[k] = v
		}
	}

	if setStr != "" {
		setVals := make(map[string]interface{})
		if err := strvals.ParseInto(setStr, setVals); err != nil {
			return nil, fmt.Errorf("failed to parse --set values: %v", err)
		}
		for k, v := range setVals {
			merged[k] = v
		}
	}

	return merged, nil
}

func ListHelmReleases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := request.GetString("namespace", "")
	outFmt := request.GetString("output", "")
	allNS := request.GetBool("allNamespaces", false)

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}
	listAction := action.NewList(cfg)
	listAction.AllNamespaces = allNS
	releases, err := listAction.Run()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error listing releases: %v", err)), nil
	}

	if outFmt != "" {
		out, err := output.Format(outFmt, releases)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(out), nil
	}

	var out []releaseData
	for _, r := range releases {
		out = append(out, releaseSummary(r))
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func GetHelmRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}
	outFmt := request.GetString("output", "")

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}
	getAction := action.NewGet(cfg)
	r, err := getAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting release %s/%s: %v", ns, name, err)), nil
	}

	if outFmt != "" {
		out, err := output.Format(outFmt, r)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error marshalling: %v", err)), nil
		}
		return mcp.NewToolResultText(out), nil
	}

	b, _ := json.MarshalIndent(releaseSummary(r), "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func GetHelmReleaseValues(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}
	getAction := action.NewGetValues(cfg)
	vals, err := getAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting values for %s/%s: %v", ns, name, err)), nil
	}

	b, _ := json.MarshalIndent(vals, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func InstallHelmRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}
	chart := request.GetString("chart", "")
	chartVersion := request.GetString("version", "")
	valuesStr := request.GetString("values", "")
	setStr := request.GetString("set", "")

	if chart == "" {
		return mcp.NewToolResultText("Provide chart name or path"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	settings := getSettings()

	installAction := action.NewInstall(cfg)
	installAction.ReleaseName = name
	installAction.Namespace = ns
	if chartVersion != "" {
		installAction.Version = chartVersion
	}

	cp, err := installAction.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error locating chart %s: %v", chart, err)), nil
	}

	loaded, err := loader.Load(cp)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading chart: %v", err)), nil
	}

	vals, err := buildValues(valuesStr, setStr)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error processing values: %v", err)), nil
	}

	r, err := installAction.Run(loaded, vals)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error installing release %s/%s: %v", ns, name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Release %s/%s installed (version %d)", r.Namespace, r.Name, r.Version)), nil
}

func UpgradeHelmRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}
	chart := request.GetString("chart", "")
	chartVersion := request.GetString("version", "")
	valuesStr := request.GetString("values", "")
	setStr := request.GetString("set", "")
	reuseValues := request.GetBool("reuseValues", true)

	if chart == "" {
		return mcp.NewToolResultText("Provide chart name or path"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	settings := getSettings()

	upgradeAction := action.NewUpgrade(cfg)
	upgradeAction.ReuseValues = reuseValues
	if chartVersion != "" {
		upgradeAction.Version = chartVersion
	}

	cp, err := upgradeAction.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error locating chart %s: %v", chart, err)), nil
	}

	loaded, err := loader.Load(cp)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading chart: %v", err)), nil
	}

	vals, err := buildValues(valuesStr, setStr)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error processing values: %v", err)), nil
	}

	r, err := upgradeAction.Run(name, loaded, vals)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error upgrading release %s/%s: %v", ns, name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Release %s/%s upgraded (version %d)", r.Namespace, r.Name, r.Version)), nil
}

func UninstallHelmRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	uninstallAction := action.NewUninstall(cfg)
	_, err = uninstallAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error uninstalling release %s/%s: %v", ns, name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Release %s/%s uninstalled", ns, name)), nil
}

func RollbackHelmRelease(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	rollbackAction := action.NewRollback(cfg)
	if rev := request.GetInt("revision", 0); rev > 0 {
		rollbackAction.Version = int(rev)
	}

	if err := rollbackAction.Run(name); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error rolling back %s/%s: %v", ns, name, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Release %s/%s rolled back", ns, name)), nil
}

func GetHelmReleaseHistory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	historyAction := action.NewHistory(cfg)
	if max := request.GetInt("max", 0); max > 0 {
		historyAction.Max = int(max)
	}

	releases, err := historyAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting history for %s/%s: %v", ns, name, err)), nil
	}

	var out []releaseHistoryData
	for _, r := range releases {
		out = append(out, releaseHistoryData{
			Revision:    r.Version,
			Status:      r.Info.Status.String(),
			Chart:       r.Chart.Name() + "-" + r.Chart.AppVersion(),
			AppVersion:  r.Chart.AppVersion(),
			Description: r.Info.Description,
			Updated:     r.Info.LastDeployed.Local().Format(time.RFC3339),
		})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func GetHelmReleaseManifest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	getAction := action.NewGet(cfg)
	r, err := getAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting release %s/%s: %v", ns, name, err)), nil
	}

	return mcp.NewToolResultText(r.Manifest), nil
}

func GetHelmReleaseNotes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide name for release"), nil
	}
	ns, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultText("Provide namespace for release"), nil
	}

	cfg, err := getActionConfig(ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error initializing helm: %v", err)), nil
	}

	getAction := action.NewGet(cfg)
	r, err := getAction.Run(name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error getting release %s/%s: %v", ns, name, err)), nil
	}

	notes := r.Info.Notes
	if strings.TrimSpace(notes) == "" {
		return mcp.NewToolResultText("(no notes)"), nil
	}
	return mcp.NewToolResultText(notes), nil
}

func ListHelmRepos(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	settings := getSettings()

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading repo file: %v", err)), nil
	}

	var out []map[string]string
	for _, re := range f.Repositories {
		out = append(out, map[string]string{
			"name": re.Name,
			"url":  re.URL,
		})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func AddHelmRepo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide repo name"), nil
	}
	url, err := request.RequireString("url")
	if err != nil {
		return mcp.NewToolResultText("Provide repo URL"), nil
	}

	entry := &repo.Entry{
		Name: name,
		URL:  url,
	}

	settings := getSettings()

	r, err := repo.NewChartRepository(entry, getter.All(settings))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error creating repo: %v", err)), nil
	}

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		f = repo.NewFile()
	}

	f.Add(entry)

	if _, err := r.DownloadIndexFile(); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Warning: repo added but index download failed: %v", err)), nil
	}

	if err := f.WriteFile(settings.RepositoryConfig, 0644); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error writing repo file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Repo %s (%s) added and index downloaded", name, url)), nil
}

func RemoveHelmRepo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultText("Provide repo name"), nil
	}

	settings := getSettings()

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading repo file: %v", err)), nil
	}

	if !f.Remove(name) {
		return mcp.NewToolResultText(fmt.Sprintf("Repo %s not found", name)), nil
	}

	if err := f.WriteFile(settings.RepositoryConfig, 0644); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error writing repo file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Repo %s removed", name)), nil
}

func UpdateHelmRepos(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	settings := getSettings()

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading repo file: %v", err)), nil
	}

	repoName := request.GetString("name", "")
	var results []string

	for _, entry := range f.Repositories {
		if repoName != "" && entry.Name != repoName {
			continue
		}

		cr, err := repo.NewChartRepository(entry, getter.All(settings))
		if err != nil {
			results = append(results, fmt.Sprintf("%s: error creating repo: %v", entry.Name, err))
			continue
		}

		if _, err := cr.DownloadIndexFile(); err != nil {
			results = append(results, fmt.Sprintf("%s: error downloading index: %v", entry.Name, err))
			continue
		}

		results = append(results, fmt.Sprintf("%s: index updated", entry.Name))
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No repositories found"), nil
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func SearchHelmRepo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyword, err := request.RequireString("keyword")
	if err != nil {
		return mcp.NewToolResultText("Provide search keyword"), nil
	}

	settings := getSettings()

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error loading repo file: %v", err)), nil
	}

	type chartResult struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		AppVersion  string `json:"appVersion"`
		Description string `json:"description"`
		Repo        string `json:"repo"`
	}

	var results []chartResult
	keywordLower := strings.ToLower(keyword)

	for _, entry := range f.Repositories {
		idxFile, err := repo.LoadIndexFile(filepath.Join(settings.RepositoryCache, helmpath.CacheIndexFile(entry.Name)))
		if err != nil {
			continue
		}

		for chartName, chartVersions := range idxFile.Entries {
			if len(chartVersions) == 0 {
				continue
			}
			latest := chartVersions[0]
			if strings.Contains(strings.ToLower(chartName), keywordLower) ||
				strings.Contains(strings.ToLower(latest.Description), keywordLower) ||
				strings.Contains(strings.ToLower(latest.AppVersion), keywordLower) {
				results = append(results, chartResult{
					Name:        chartName,
					Version:     latest.Version,
					AppVersion:  latest.AppVersion,
					Description: latest.Description,
					Repo:        entry.Name,
				})
			}
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No charts found matching %q", keyword)), nil
	}

	b, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}
