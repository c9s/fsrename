package workers

import "testing"
import "github.com/stretchr/testify/assert"

func TestRegExpAction(t *testing.T) {
	act := NewRegExpActionWithPattern("\\.php$", ".txt")
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
