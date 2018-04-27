package mesos_singularity

import singularity "github.com/lenfree/go-singularity"

func clientConn(m interface{}) *singularity.Client {
	return m.(*Conn).sclient
}
