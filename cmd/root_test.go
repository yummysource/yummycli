package cmd

import (
	"testing"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

func TestRootCommandIncludesCoreGroups(t *testing.T) {
	f := cmdutil.NewDefault()

	root := NewRootCommand(f)

	expected := []string{"auth", "gemini", "image", "version", "video"}

	assertSubcommands(t, root, expected)
}

func TestGeminiCommandIncludesPhaseOneAliases(t *testing.T) {
	f := cmdutil.NewDefault()
	root := NewRootCommand(f)
	gemini, _, err := root.Find([]string{"gemini"})
	if err != nil {
		t.Fatalf("find gemini command: %v", err)
	}

	expected := []string{"init", "nanobanana", "veo"}
	assertSubcommands(t, gemini, expected)

	for _, check := range []struct {
		name      string
		canonical string
	}{
		{name: "init", canonical: "auth init --provider gemini"},
		{name: "nanobanana", canonical: "image generate --provider gemini --preset nano-banana"},
		{name: "veo", canonical: "video generate --provider gemini"},
	} {
		command, _, err := gemini.Find([]string{check.name})
		if err != nil {
			t.Fatalf("find gemini %s command: %v", check.name, err)
		}

		got := command.Annotations["canonical"]

		if got != check.canonical {
			t.Fatalf("canonical target for %s = %q, want %q", check.name, got, check.canonical)
		}
	}
}

func TestImageCommandIncludesPhaseOneActions(t *testing.T) {
	f := cmdutil.NewDefault()
	root := NewRootCommand(f)
	image, _, err := root.Find([]string{"image"})
	if err != nil {
		t.Fatalf("find image command: %v", err)
	}

	expected := []string{"generate"}

	assertSubcommands(t, image, expected)
}

func TestAuthCommandIncludesCoreActions(t *testing.T) {
	f := cmdutil.NewDefault()

	root := NewRootCommand(f)
	auth, _, err := root.Find([]string{"auth"})
	if err != nil {
		t.Fatalf("find auth command: %v", err)
	}

	expected := []string{"init", "list", "remove", "status"}
	assertSubcommands(t, auth, expected)
}
