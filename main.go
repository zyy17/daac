package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Dashboard struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Title   string   `yaml:"title"`
	Style   *Style   `yaml:"style"`
	Groups  []*Group `yaml:"groups"`
}

type Style struct {
	Reference string `yaml:"reference"`
}

type Group struct {
	Name   string   `yaml:"name"`
	Title  string   `yaml:"title"`
	IsRow  bool     `yaml:"isRow"`
	Panels []*Panel `yaml:"panels"`
}

type Panel struct {
	Name        string   `yaml:"name"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Queries     []*Query `yaml:"queries"`
}

type Query struct {
	Expr string `yaml:"expr"`
}

type DashboardRenderVars struct {
	DashboardTitle string
	DashboardUID   string

	// FIXME(zyy17): should be support multiple panels.
	PanelTitle string
	PanelExpr  string
	PanelType  string
}

func main() {
	var (
		dashboardInputFile  string
		dashboardOutputFile string
	)

	flag.StringVar(&dashboardInputFile, "f", "", "input dashboard yaml file")
	flag.StringVar(&dashboardOutputFile, "o", "", "output dashboard json file")
	flag.Parse()

	ds, err := parseDashboard(dashboardInputFile)
	if err != nil {
		log.Fatalf("failed to parse dashboard: %v", err)
	}

	if err = render(ds, dashboardOutputFile); err != nil {
		log.Fatalf("failed to render dashboard: %v", err)
	}
}

func parseDashboard(dashboardFile string) (*Dashboard, error) {
	data, err := os.ReadFile(dashboardFile)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = yaml.Unmarshal(data, &dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

// TODO(zyy17): For the serious usage, we need to design a more complicated render and layout engine.
func render(ds *Dashboard, output string) error {
	if len(ds.Style.Reference) == 0 {
		return fmt.Errorf("style reference is empty")
	}

	tmpl, err := template.ParseFiles(ds.Style.Reference)
	if err != nil {
		return err
	}

	var drv DashboardRenderVars
	drv.DashboardTitle = ds.Title
	drv.DashboardUID = uuid.New().String()
	drv.PanelTitle = ds.Groups[0].Panels[0].Title
	drv.PanelExpr = ds.Groups[0].Panels[0].Queries[0].Expr

	// TODO(zyy17): it's just a trick...
	if ds.Groups[0].IsRow {
		drv.PanelType = "row"
	} else {
		drv.PanelType = "timeseries"
	}

	// Ouptut to file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}

	if err = tmpl.Execute(f, drv); err != nil {
		return err
	}

	return nil
}
