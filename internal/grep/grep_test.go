package grep_test

import (
	"testing"

	"github.com/user/greplive/internal/grep"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := grep.New([]string{"[invalid"}, true)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_EmptyPatterns_Disabled(t *testing.T) {
	m, err := grep.New([]string{""}, false)
	if err != nil {
		t.Fatal(err)
	}
	if m.Enabled() {
		t.Fatal("expected disabled matcher")
	}
}

func TestMatch_NoPatterns_AlwaysTrue(t *testing.T) {
	m, _ := grep.New(nil, false)
	if !m.Match("anything") {
		t.Fatal("expected true with no patterns")
	}
}

func TestMatch_AnyLogic(t *testing.T) {
	m, _ := grep.New([]string{"foo", "bar"}, true)
	if !m.Match("foo baz") {
		t.Error("expected match on foo")
	}
	if m.Match("baz qux") {
		t.Error("expected no match")
	}
}

func TestMatch_AllLogic(t *testing.T) {
	m, _ := grep.New([]string{"foo", "bar"}, false)
	if !m.Match("foo bar") {
		t.Error("expected match when both present")
	}
	if m.Match("foo only") {
		t.Error("expected no match when bar missing")
	}
}

func TestGroups_NamedCaptures(t *testing.T) {
	m, err := grep.New([]string{`(?P<level>\w+) (?P<msg>.+)`}, true)
	if err != nil {
		t.Fatal(err)
	}
	groups := m.Groups("ERROR something went wrong")
	if groups["level"] != "ERROR" {
		t.Errorf("got level=%q", groups["level"])
	}
	if groups["msg"] != "something went wrong" {
		t.Errorf("got msg=%q", groups["msg"])
	}
}

func TestGroups_NoMatch_ReturnsNil(t *testing.T) {
	m, _ := grep.New([]string{"foo"}, true)
	if m.Groups("bar") != nil {
		t.Fatal("expected nil groups")
	}
}
