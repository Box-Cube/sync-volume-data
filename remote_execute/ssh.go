/*
Copyright 2021 Box-Cube

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package remote

import (
	"errors"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

type Cli struct {
	user       string
	pwd        string
	addr       string
	client     *gossh.Client
	session    *gossh.Session
	LastResult string
}

func (c *Cli) Connect() (*Cli, error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []gossh.AuthMethod{gossh.Password(c.pwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	client, err := gossh.Dial("tcp", c.addr, config)
	if nil != err {
		return c, err
	}
	c.client = client
	return c, nil
}

func (c *Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	stderr, err := session.StderrPipe()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := session.Start(shell); err != nil {
		return "", err
	}
	errMsg, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", err
	} else if len(errMsg) > 0 {
		return "", errors.New(string(errMsg))
	}

	//get only a row as expected
	stdoutMsg, _ := ioutil.ReadAll(stdout)
	if err := session.Wait(); err != nil {
		return "", err
	}
	return string(stdoutMsg), nil
}

func NewCli(user, pwd, addr string) *Cli {
	return &Cli{
		user: user,
		pwd:  pwd,
		addr: addr,
	}
}
