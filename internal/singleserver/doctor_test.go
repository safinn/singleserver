package singleserver

import "testing"

func TestDoctorAppsReturnsAllWhenNoFilter(t *testing.T) {
	apps := []AppConfig{
		{Repo: "dvassallo/fullsend", Name: "fullsend"},
		{Repo: "smallbets/userbase-homepage", Name: "userbase-homepage"},
	}

	got, err := doctorApps(apps, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected all apps, got %d", len(got))
	}
}

func TestDoctorAppsFiltersByNameOrRepo(t *testing.T) {
	apps := []AppConfig{
		{Repo: "dvassallo/fullsend", Name: "fullsend"},
		{Repo: "smallbets/userbase-homepage", Name: "userbase-homepage"},
	}

	byName, err := doctorApps(apps, []string{"fullsend"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byName) != 1 || byName[0].Repo != "dvassallo/fullsend" {
		t.Fatalf("unexpected app selected by name: %#v", byName)
	}

	byRepo, err := doctorApps(apps, []string{"smallbets/userbase-homepage"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byRepo) != 1 || byRepo[0].Name != "userbase-homepage" {
		t.Fatalf("unexpected app selected by repo: %#v", byRepo)
	}
}

func TestDoctorAppsRejectsUnknownAndExtraArgs(t *testing.T) {
	apps := []AppConfig{{Repo: "dvassallo/fullsend", Name: "fullsend"}}

	if _, err := doctorApps(apps, []string{"missing"}); err == nil {
		t.Fatal("expected unknown app to fail")
	}
	if _, err := doctorApps(apps, []string{"fullsend", "extra"}); err == nil {
		t.Fatal("expected extra args to fail")
	}
}

func TestCloudflaredRoutes(t *testing.T) {
	config := &cloudflaredConfig{
		Ingress: []cloudflaredIngress{
			{Hostname: "Hooks.Example.com", Service: "http://127.0.0.1:8787"},
			{Hostname: "app.example.com", Service: "http://127.0.0.1:80"},
			{Service: "http_status:404"},
		},
	}

	routes := cloudflaredRoutes(config)
	if routes["hooks.example.com"] != "http://127.0.0.1:8787" {
		t.Fatalf("unexpected hook route: %#v", routes)
	}
	if routes["app.example.com"] != "http://127.0.0.1:80" {
		t.Fatalf("unexpected app route: %#v", routes)
	}
	if _, ok := routes[""]; ok {
		t.Fatal("fallback route should not be keyed")
	}
}

func TestAppsHaveHosts(t *testing.T) {
	if appsHaveHosts([]AppConfig{{Repo: "dvassallo/fullsend", Name: "fullsend"}}) {
		t.Fatal("expected no hosts")
	}
	if !appsHaveHosts([]AppConfig{{Repo: "dvassallo/fullsend", Name: "fullsend", Hosts: []string{"fullsend.game"}}}) {
		t.Fatal("expected hosts")
	}
}

func TestExpectedCloudflaredHosts(t *testing.T) {
	hosts := expectedCloudflaredHosts("Hooks.Example.com", []AppConfig{
		{Repo: "dvassallo/fullsend", Name: "fullsend", Hosts: []string{"Fullsend.Game", "fullsend.assetstacks.com"}},
		{Repo: "smallbets/userbase-homepage", Name: "userbase-homepage"},
	})

	for _, host := range []string{"hooks.example.com", "fullsend.game", "fullsend.assetstacks.com"} {
		if !hosts[host] {
			t.Fatalf("expected host %s in %#v", host, hosts)
		}
	}
	if hosts["userbase.com"] {
		t.Fatalf("unexpected host set: %#v", hosts)
	}
}

func TestStaleCloudflaredHosts(t *testing.T) {
	routes := map[string]string{
		"hooks.example.com": "http://127.0.0.1:8787",
		"old.example.com":   "http://127.0.0.1:80",
		"z.example.com":     "http://127.0.0.1:80",
	}
	expected := map[string]bool{"hooks.example.com": true}

	got := staleCloudflaredHosts(routes, expected)
	if len(got) != 2 || got[0] != "old.example.com" || got[1] != "z.example.com" {
		t.Fatalf("unexpected stale hosts: %#v", got)
	}
}

func TestFormatBytesGB(t *testing.T) {
	if got := formatBytesGB(1536 * 1024 * 1024); got != "1.5GB" {
		t.Fatalf("unexpected formatted size: %s", got)
	}
}

func TestLastDeployStatusFromJournalUsesMostRecentOutcome(t *testing.T) {
	journal := `
[deploy:fullsend-1] success total_ms=1200
[deploy:userbase-homepage-1] success total_ms=900
[deploy:fullsend-2] failed after 300ms: boom
`
	status, detail := lastDeployStatusFromJournal("fullsend", journal)
	if status != "failed" {
		t.Fatalf("unexpected status: %s", status)
	}
	if detail != "failed after 300ms: boom" {
		t.Fatalf("unexpected detail: %s", detail)
	}
}

func TestLastDeployStatusFromJournalReportsUnknown(t *testing.T) {
	status, detail := lastDeployStatusFromJournal("sillyface-games", "[server] ok")
	if status != "unknown" {
		t.Fatalf("unexpected status: %s", status)
	}
	if detail != "no recent deploy outcome" {
		t.Fatalf("unexpected detail: %s", detail)
	}
}

func TestCompactWhitespace(t *testing.T) {
	got := compactWhitespace(" M file\n?? other\n")
	if got != "M file ?? other" {
		t.Fatalf("unexpected value: %q", got)
	}
}
