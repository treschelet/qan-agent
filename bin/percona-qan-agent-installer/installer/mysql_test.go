/*
   Copyright (c) 2016, Percona LLC and/or its affiliates. All rights reserved.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package installer_test

import (
	"github.com/percona/go-mysql/dsn"
	i "github.com/percona/qan-agent/bin/percona-qan-agent-installer/installer"
	. "gopkg.in/check.v1"
)

type MySQLTestSuite struct {
}

var _ = Suite(&MySQLTestSuite{})

// --------------------------------------------------------------------------

func (s *MySQLTestSuite) TestMakeGrant(t *C) {
	dsnTest := dsn.DSN{
		Username: "new-user",
		Password: "some pass",
	}

	dsnTest.Hostname = "localhost"
	maxOpenConnections := int64(1)
	got := i.MakeGrant(dsnTest, maxOpenConnections)
	expect := []string{
		"SET SESSION old_passwords=0",
		"GRANT SUPER, PROCESS, USAGE, SELECT, RELOAD ON *.* TO 'new-user'@'localhost' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
		"GRANT UPDATE, DELETE, DROP ON performance_schema.* TO 'new-user'@'localhost' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
	}
	t.Check(got, DeepEquals, expect)

	dsnTest.Hostname = "127.0.0.1"
	got = i.MakeGrant(dsnTest, maxOpenConnections)
	expect = []string{
		"SET SESSION old_passwords=0",
		"GRANT SUPER, PROCESS, USAGE, SELECT, RELOAD ON *.* TO 'new-user'@'127.0.0.1' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
		"GRANT UPDATE, DELETE, DROP ON performance_schema.* TO 'new-user'@'127.0.0.1' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
	}
	t.Check(got, DeepEquals, expect)

	dsnTest.Hostname = "10.1.1.1"
	got = i.MakeGrant(dsnTest, maxOpenConnections)
	expect = []string{
		"SET SESSION old_passwords=0",
		"GRANT SUPER, PROCESS, USAGE, SELECT, RELOAD ON *.* TO 'new-user'@'%' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
		"GRANT UPDATE, DELETE, DROP ON performance_schema.* TO 'new-user'@'%' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
	}
	t.Check(got, DeepEquals, expect)

	dsnTest.Hostname = ""
	dsnTest.Socket = "/var/lib/mysql.sock"
	got = i.MakeGrant(dsnTest, maxOpenConnections)
	expect = []string{
		"SET SESSION old_passwords=0",
		"GRANT SUPER, PROCESS, USAGE, SELECT, RELOAD ON *.* TO 'new-user'@'localhost' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
		"GRANT UPDATE, DELETE, DROP ON performance_schema.* TO 'new-user'@'localhost' IDENTIFIED BY 'some pass' WITH MAX_USER_CONNECTIONS 1",
	}
	t.Check(got, DeepEquals, expect)
}
