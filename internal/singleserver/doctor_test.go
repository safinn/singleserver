package singleserver

import "testing"

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
