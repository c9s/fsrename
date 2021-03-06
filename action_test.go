package fsrename

import "testing"
import "github.com/stretchr/testify/assert"

func TestUnderscoreAction(t *testing.T) {
	act := NewUnderscoreAction()
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/CamelCase.php")
	assert.Nil(t, err)
	act.Act(entry)
	assert.Equal(t, "tests/camel_case.php", entry.newpath)
}

func TestExtReplaceAction(t *testing.T) {
	act := NewExtReplaceAction("twig")
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/file.html")
	assert.Nil(t, err)
	assert.True(t, act.Act(entry))
	assert.Equal(t, "tests/file.twig", entry.newpath)
}

func TestExtReplaceAction2(t *testing.T) {
	act := NewExtReplaceAction("twig")
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/file.html.twig")
	assert.Nil(t, err)
	assert.True(t, act.Act(entry))
	assert.Equal(t, "tests/file.html.twig", entry.newpath)
}

func TestCamelCaseAction(t *testing.T) {
	act := NewCamelCaseAction("[-_]+")
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/foo-bar.php")
	assert.Nil(t, err)
	act.Act(entry)
	assert.Equal(t, "tests/fooBar.php", entry.newpath)
}

func TestRegExpAction(t *testing.T) {
	act := NewRegExpReplaceActionWithPattern("\\.php$", ".txt")
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/autoload.php")
	assert.Nil(t, err)
	act.Act(entry)
	assert.Equal(t, "tests/autoload.txt", entry.newpath)
}

func TestStrReplaceAction(t *testing.T) {
	act := NewStrReplaceAction(".php", ".txt", 1)
	assert.NotNil(t, act)
	entry, err := NewFileEntry("tests/autoload.php")
	assert.Nil(t, err)
	act.Act(entry)
	assert.Equal(t, "tests/autoload.txt", entry.newpath)
}
