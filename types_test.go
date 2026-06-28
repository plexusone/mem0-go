package mem0

import "testing"

func TestRoleConstants(t *testing.T) {
	if RoleUser != "user" {
		t.Errorf("expected 'user', got %q", RoleUser)
	}
	if RoleAssistant != "assistant" {
		t.Errorf("expected 'assistant', got %q", RoleAssistant)
	}
	if RoleSystem != "system" {
		t.Errorf("expected 'system', got %q", RoleSystem)
	}
}

func TestMemoryEventConstants(t *testing.T) {
	if MemoryEventAdd != "ADD" {
		t.Errorf("expected 'ADD', got %q", MemoryEventAdd)
	}
	if MemoryEventUpdate != "UPDATE" {
		t.Errorf("expected 'UPDATE', got %q", MemoryEventUpdate)
	}
	if MemoryEventDelete != "DELETE" {
		t.Errorf("expected 'DELETE', got %q", MemoryEventDelete)
	}
	if MemoryEventNoop != "NOOP" {
		t.Errorf("expected 'NOOP', got %q", MemoryEventNoop)
	}
}

func TestEntityTypeConstants(t *testing.T) {
	if EntityTypeUser != "user" {
		t.Errorf("expected 'user', got %q", EntityTypeUser)
	}
	if EntityTypeAgent != "agent" {
		t.Errorf("expected 'agent', got %q", EntityTypeAgent)
	}
	if EntityTypeApp != "app" {
		t.Errorf("expected 'app', got %q", EntityTypeApp)
	}
	if EntityTypeRun != "run" {
		t.Errorf("expected 'run', got %q", EntityTypeRun)
	}
}
