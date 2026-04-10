package image

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMIMETypeFromPath(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "png", path: "a.png", want: "image/png"},
		{name: "jpg", path: "a.jpg", want: "image/jpeg"},
		{name: "jpeg", path: "a.jpeg", want: "image/jpeg"},
		{name: "webp", path: "a.webp", want: "image/webp"},
		{name: "unsupported", path: "a.gif", wantErr: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := mimeTypeFromPath(testCase.path)
			if testCase.wantErr {
				if err == nil {
					t.Fatal("mimeTypeFromPath returned nil error")
				}
				return
			}

			if err != nil {
				t.Fatalf("mimeTypeFromPath returned error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("mime type = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestBuildGenerateContentPartsWithoutInputImage(t *testing.T) {
	req := GenerateImageRequest{
		Prompt: "turn this into watercolor",
	}

	parts, err := buildGenerateContentParts(req)
	if err != nil {
		t.Fatalf("buildGenerateContentParts returned error: %v", err)
	}

	if len(parts) != 1 {
		t.Fatalf("parts length = %d, want %d", len(parts), 1)
	}
	if parts[0].Text != "turn this into watercolor" {
		t.Fatalf("text part = %q, want %q", parts[0].Text, "turn this into watercolor")
	}
}

func TestBuildGenerateContentPartsWithInputImage(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "source.png")

	err := os.WriteFile(inputPath, []byte("fake-image-bytes"), 0o644)
	if err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	req := GenerateImageRequest{
		Prompt:      "turn this into watercolor",
		InputImages: []string{inputPath},
	}

	parts, err := buildGenerateContentParts(req)
	if err != nil {
		t.Fatalf("buildGenerateContentParts returned error: %v", err)
	}

	if len(parts) != 2 {
		t.Fatalf("parts length = %d, want %d", len(parts), 2)
	}
	if parts[0].InlineData == nil {
		t.Fatal("first part inline data is nil")
	}
	if parts[0].InlineData.MIMEType != "image/png" {
		t.Fatalf("mime type = %q, want %q", parts[0].InlineData.MIMEType, "image/png")
	}
	if string(parts[0].InlineData.Data) != "fake-image-bytes" {
		t.Fatalf("image bytes = %q, want %q", string(parts[0].InlineData.Data), "fake-image-bytes")
	}
	if parts[1].Text != "turn this into watercolor" {
		t.Fatalf("text part = %q, want %q", parts[1].Text, "turn this into watercolor")
	}
}
