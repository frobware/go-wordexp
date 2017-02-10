package wordexp_test

import (
	"os"

	wordexp "github.com/frobware/go-wordexp"
	gc "gopkg.in/check.v1"
)

type WordExpansionSuite struct{}

var _ = gc.Suite(&WordExpansionSuite{})

func (s *WordExpansionSuite) TestWordExpansionBadCharacterError(c *gc.C) {
	_, err := wordexp.Expand(`testdata/|`)
	c.Check(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, "bad character")
}

func (s *WordExpansionSuite) TestWordExpansionBadVariableError(c *gc.C) {
	_, err := wordexp.Expand(`${TestWordExpansionBadVariableError}`)
	c.Check(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, "bad variable")
}

func (s *WordExpansionSuite) TestWordExpansionCmdSubstitutionError(c *gc.C) {
	_, err := wordexp.Expand(`$(uname -a)`)
	c.Check(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, "command execution not allowed")
}

func (s *WordExpansionSuite) TestWordExpansionOutOfMemoryError(c *gc.C) {
	allocator := wordexp.NewOutOfMemoryAllocator()
	expander := wordexp.NewWordExpander(allocator)
	_, err := expander.Expand("testdata/wordexp/*.cfg")
	c.Assert(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, "not enough memory to store the result")
}

func (s *WordExpansionSuite) TestWordExpansionShellSyntaxError(c *gc.C) {
	_, err := wordexp.Expand(`testdata/"`)
	c.Check(err, gc.NotNil)
	c.Assert(err, gc.ErrorMatches, "shell syntax error in words")
}

func (s *WordExpansionSuite) TestWordExpansionForNonExistentFiles(c *gc.C) {
	results, err := wordexp.Expand("testdata/wordexp/*.nonexistent")
	c.Assert(err, gc.IsNil)
	c.Assert(results, gc.HasLen, 1)
	c.Check(results[0], gc.Equals, "testdata/wordexp/*.nonexistent")
}

func (s *WordExpansionSuite) TestWordExpansionForExpectedFiles(c *gc.C) {
	results, err := wordexp.Expand("testdata/wordexp/*.cfg")
	c.Assert(err, gc.IsNil)
	c.Assert(results, gc.HasLen, 3)
	c.Check(results[0], gc.Equals, "testdata/wordexp/a.cfg")
	c.Check(results[1], gc.Equals, "testdata/wordexp/b.cfg")
	c.Check(results[2], gc.Equals, "testdata/wordexp/c.cfg")
}

func (s *WordExpansionSuite) TestWordExpansionNotAnError(c *gc.C) {
	err := wordexp.NewWordExpansionError(0)
	c.Assert(err, gc.NotNil)
	c.Check(err, gc.ErrorMatches, "not an error")
}

func (s *WordExpansionSuite) TestWordExpansionUnknownError(c *gc.C) {
	err := wordexp.NewWordExpansionError(-1)
	c.Assert(err, gc.NotNil)
	c.Check(err, gc.ErrorMatches, "unrecognised error code")
}

func (s *WordExpansionSuite) TestWordExpansionNoMatches(c *gc.C) {
	results, err := wordexp.Expand("testdata/wordexp/*.go")
	c.Assert(err, gc.IsNil)
	c.Assert(results, gc.HasLen, 1)
	c.Assert(results[0], gc.Equals, "testdata/wordexp/*.go")
	_, statError := os.Stat(results[0])
	c.Assert(os.IsNotExist(statError), gc.Equals, true)
}
