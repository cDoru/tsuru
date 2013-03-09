// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"bytes"
	"encoding/json"
	"github.com/globocom/tsuru/cmd"
	. "launchpad.net/gocheck"
	"net/http"
)

func (s *S) TestAppInfo(c *C) {
	*AppName = "app1"
	var stdout, stderr bytes.Buffer
	result := `{"Name":"app1","CName":"","Ip":"myapp.tsuru.io","Framework":"php","Repository":"git@git.com:php.git","State":"dead", "Units":[{"Ip":"10.10.10.10","Name":"app1/0","State":"started"}, {"Ip":"9.9.9.9","Name":"app1/1","State":"started"}, {"Ip":"","Name":"app1/2","State":"pending"}],"Teams":["tsuruteam","crane"]}`
	expected := `Application: app1
Repository: git@git.com:php.git
Platform: php
Teams: tsuruteam, crane
Address: myapp.tsuru.io
Units:
+--------+---------+
| Unit   | State   |
+--------+---------+
| app1/0 | started |
| app1/1 | started |
| app1/2 | pending |
+--------+---------+

`
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppInfo{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppInfoNoUnits(c *C) {
	*AppName = "app1"
	var stdout, stderr bytes.Buffer
	result := `{"Name":"app1","Ip":"app1.tsuru.io","Framework":"php","Repository":"git@git.com:php.git","State":"dead", "Units":[],"Teams":["tsuruteam","crane"]}`
	expected := `Application: app1
Repository: git@git.com:php.git
Platform: php
Teams: tsuruteam, crane
Address: app1.tsuru.io

`
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppInfo{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppInfoWithoutArgs(c *C) {
	var stdout, stderr bytes.Buffer
	result := `{"Name":"secret","Ip":"secret.tsuru.io","Framework":"ruby","Repository":"git@git.com:php.git","State":"dead", "Units":[{"Ip":"10.10.10.10","Name":"secret/0","State":"started"}, {"Ip":"9.9.9.9","Name":"secret/1","State":"pending"}],"Teams":["tsuruteam","crane"]}`
	expected := `Application: secret
Repository: git@git.com:php.git
Platform: ruby
Teams: tsuruteam, crane
Address: secret.tsuru.io
Units:
+----------+---------+
| Unit     | State   |
+----------+---------+
| secret/0 | started |
| secret/1 | pending |
+----------+---------+

`
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &conditionalTransport{
		transport{
			msg:    result,
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			return req.URL.Path == "/apps/secret" && req.Method == "GET"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	fake := FakeGuesser{name: "secret"}
	guessCommand := GuessingCommand{G: &fake}
	command := AppInfo{guessCommand}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppInfoCName(c *C) {
	*AppName = "app1"
	var stdout, stderr bytes.Buffer
	result := `{"Name":"app1","Ip":"myapp.tsuru.io","CName":"yourapp.tsuru.io","Framework":"php","Repository":"git@git.com:php.git","State":"dead", "Units":[{"Ip":"10.10.10.10","Name":"app1/0","State":"started"}, {"Ip":"9.9.9.9","Name":"app1/1","State":"started"}, {"Ip":"","Name":"app1/2","State":"pending"}],"Teams":["tsuruteam","crane"]}`
	expected := `Application: app1
Repository: git@git.com:php.git
Platform: php
Teams: tsuruteam, crane
Address: yourapp.tsuru.io
Units:
+--------+---------+
| Unit   | State   |
+--------+---------+
| app1/0 | started |
| app1/1 | started |
| app1/2 | pending |
+--------+---------+

`
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppInfo{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppInfoInfo(c *C) {
	expected := &cmd.Info{
		Name:  "app-info",
		Usage: "app-info [--app appname]",
		Desc: `show information about your app.

If you don't provide the app name, tsuru will try to guess it.`,
		MinArgs: 0,
	}
	c.Assert((&AppInfo{}).Info(), DeepEquals, expected)
}

func (s *S) TestAppGrant(c *C) {
	*AppName = "games"
	var stdout, stderr bytes.Buffer
	expected := `Team "cobrateam" was added to the "games" app` + "\n"
	context := cmd.Context{
		Args:   []string{"cobrateam"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	command := AppGrant{}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: "", status: http.StatusOK}}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppGrantWithoutFlag(c *C) {
	var stdout, stderr bytes.Buffer
	expected := `Team "cobrateam" was added to the "fights" app` + "\n"
	context := cmd.Context{
		Args:   []string{"cobrateam"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	fake := &FakeGuesser{name: "fights"}
	command := AppGrant{GuessingCommand{G: fake}}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: "", status: http.StatusOK}}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppGrantInfo(c *C) {
	expected := &cmd.Info{
		Name:  "app-grant",
		Usage: "app-grant <teamname> [--app appname]",
		Desc: `grants access to an app to a team.

If you don't provide the app name, tsuru will try to guess it.`,
		MinArgs: 1,
	}
	c.Assert((&AppGrant{}).Info(), DeepEquals, expected)
}

func (s *S) TestAppRevoke(c *C) {
	*AppName = "games"
	var stdout, stderr bytes.Buffer
	expected := `Team "cobrateam" was removed from the "games" app` + "\n"
	context := cmd.Context{
		Args:   []string{"cobrateam"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	command := AppRevoke{}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: "", status: http.StatusOK}}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppRevokeWithoutFlag(c *C) {
	var stdout, stderr bytes.Buffer
	expected := `Team "cobrateam" was removed from the "fights" app` + "\n"
	context := cmd.Context{
		Args:   []string{"cobrateam"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	fake := &FakeGuesser{name: "fights"}
	command := AppRevoke{GuessingCommand{G: fake}}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: "", status: http.StatusOK}}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppRevokeInfo(c *C) {
	expected := &cmd.Info{
		Name:  "app-revoke",
		Usage: "app-revoke <teamname> [--app appname]",
		Desc: `revokes access to an app from a team.

If you don't provide the app name, tsuru will try to guess it.`,
		MinArgs: 1,
	}
	c.Assert((&AppRevoke{}).Info(), DeepEquals, expected)
}

func (s *S) TestAppList(c *C) {
	var stdout, stderr bytes.Buffer
	result := `[{"Ip":"10.10.10.10","Name":"app1","Units":[{"Name":"app1/0","State":"started"}]}]`
	expected := `+-------------+-------------------------+-------------+
| Application | Units State Summary     | Address     |
+-------------+-------------------------+-------------+
| app1        | 1 of 1 units in-service | 10.10.10.10 |
+-------------+-------------------------+-------------+
`
	context := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppList{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppListUnitIsntStarted(c *C) {
	var stdout, stderr bytes.Buffer
	result := `[{"Ip":"10.10.10.10","Name":"app1","Units":[{"Name":"app1/0","State":"pending"}]}]`
	expected := `+-------------+-------------------------+-------------+
| Application | Units State Summary     | Address     |
+-------------+-------------------------+-------------+
| app1        | 0 of 1 units in-service | 10.10.10.10 |
+-------------+-------------------------+-------------+
`
	context := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppList{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppListCName(c *C) {
	var stdout, stderr bytes.Buffer
	result := `[{"Ip":"10.10.10.10","CName":"app1.tsuru.io","Name":"app1","Units":[{"Name":"app1/0","State":"started"}]}]`
	expected := `+-------------+-------------------------+---------------+
| Application | Units State Summary     | Address       |
+-------------+-------------------------+---------------+
| app1        | 1 of 1 units in-service | app1.tsuru.io |
+-------------+-------------------------+---------------+
`
	context := cmd.Context{
		Args:   []string{},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	client := cmd.NewClient(&http.Client{Transport: &transport{msg: result, status: http.StatusOK}}, nil, manager)
	command := AppList{}
	err := command.Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(stdout.String(), Equals, expected)
}

func (s *S) TestAppListInfo(c *C) {
	expected := &cmd.Info{
		Name:    "app-list",
		Usage:   "app-list",
		Desc:    "list all your apps.",
		MinArgs: 0,
	}
	c.Assert(AppList{}.Info(), DeepEquals, expected)
}

func (s *S) TestAppListIsACommand(c *C) {
	var _ cmd.Command = AppList{}
}

func (s *S) TestAppRestart(c *C) {
	*AppName = "handful_of_nothing"
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			return req.URL.Path == "/apps/handful_of_nothing/restart" && req.Method == "GET"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&AppRestart{}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "Restarted")
}

func (s *S) TestAppRestartWithoutTheFlag(c *C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			return req.URL.Path == "/apps/motorbreath/restart" && req.Method == "GET"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	fake := &FakeGuesser{name: "motorbreath"}
	err := (&AppRestart{GuessingCommand{G: fake}}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "Restarted")
}

func (s *S) TestAppRestartInfo(c *C) {
	expected := &cmd.Info{
		Name:  "restart",
		Usage: "restart [--app appname]",
		Desc: `restarts an app.

If you don't provide the app name, tsuru will try to guess it.`,
		MinArgs: 0,
	}
	c.Assert((&AppRestart{}).Info(), DeepEquals, expected)
}

func (s *S) TestAppRestartIsACommand(c *C) {
	var _ cmd.Command = &AppRestart{}
}

func (s *S) TestSetCName(c *C) {
	*AppName = "death"
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Args:   []string{"death.evergrey.mycompany.com"},
	}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			var m map[string]string
			err := json.NewDecoder(req.Body).Decode(&m)
			c.Assert(err, IsNil)
			return req.URL.Path == "/apps/death" &&
				req.Method == "POST" &&
				m["cname"] == "death.evergrey.mycompany.com"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&SetCName{}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "cname successfully defined.\n")
}

func (s *S) TestSetCNameWithoutTheFlag(c *C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Args:   []string{"corey.evergrey.mycompany.com"},
	}
	fake := &FakeGuesser{name: "corey"}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			var m map[string]string
			err := json.NewDecoder(req.Body).Decode(&m)
			c.Assert(err, IsNil)
			return req.URL.Path == "/apps/corey" &&
				req.Method == "POST" &&
				m["cname"] == "corey.evergrey.mycompany.com"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&SetCName{GuessingCommand{G: fake}}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "cname successfully defined.\n")
}

func (s *S) TestSetCNameFailure(c *C) {
	var stdout, stderr bytes.Buffer
	*AppName = "masterplan"
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
		Args:   []string{"masterplan.evergrey.mycompany.com"},
	}
	trans := &transport{msg: "Invalid cname", status: http.StatusPreconditionFailed}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&SetCName{}).Run(&context, client)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Invalid cname")
}

func (s *S) TestSetCNameInfo(c *C) {
	expected := &cmd.Info{
		Name:    "set-cname",
		Usage:   "set-cname <cname> [--app appname]",
		Desc:    `defines a cname for your app.`,
		MinArgs: 1,
	}
	c.Assert((&SetCName{}).Info(), DeepEquals, expected)
}

func (s *S) TestSetCNameIsACommand(c *C) {
	var _ cmd.Command = &SetCName{}
}

func (s *S) TestUnsetCName(c *C) {
	*AppName = "death"
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			var m map[string]string
			err := json.NewDecoder(req.Body).Decode(&m)
			c.Assert(err, IsNil)
			return req.URL.Path == "/apps/death" &&
				req.Method == "POST" &&
				m["cname"] == ""
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&UnsetCName{}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "cname successfully undefined.\n")
}

func (s *S) TestUnsetCNameWithoutTheFlag(c *C) {
	var (
		called         bool
		stdout, stderr bytes.Buffer
	)
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	fake := &FakeGuesser{name: "corey"}
	trans := &conditionalTransport{
		transport{
			msg:    "Restarted",
			status: http.StatusOK,
		},
		func(req *http.Request) bool {
			called = true
			var m map[string]string
			err := json.NewDecoder(req.Body).Decode(&m)
			c.Assert(err, IsNil)
			return req.URL.Path == "/apps/corey" &&
				req.Method == "POST" &&
				m["cname"] == ""
		},
	}
	client := cmd.NewClient(&http.Client{Transport: trans}, nil, manager)
	err := (&UnsetCName{GuessingCommand{G: fake}}).Run(&context, client)
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(stdout.String(), Equals, "cname successfully undefined.\n")
}

func (s *S) TestUnsetCNameInfo(c *C) {
	expected := &cmd.Info{
		Name:    "unset-cname",
		Usage:   "unset-cname [--app appname]",
		Desc:    `unsets the current cname of your app.`,
		MinArgs: 0,
	}
	c.Assert((&UnsetCName{}).Info(), DeepEquals, expected)
}

func (s *S) TestUnsetCNameIsACommand(c *C) {
	var _ cmd.Command = &UnsetCName{}
}