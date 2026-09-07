package server

import "github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"

func multi(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 0 {
		return protocol.NewError("wrong number of arguments for 'MULTI'")
	}
	if client.inTransaction {
		return protocol.NewError("'MULTI' calls can not be nested")
	}
	client.inTransaction = true
	client.queue = nil

	return protocol.NewSimpleString("OK")
}

func discard(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 0 {
		return protocol.NewError("wrong number of arguments for 'DISCARD'")
	}
	if !client.inTransaction {
		return protocol.NewError("DISCARD was called without MULTI")
	}
	client.inTransaction = false
	client.queue = nil

	return protocol.NewSimpleString("OK")
}

func exec(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 0 {
		return protocol.NewError("wrong number of arguments for 'EXEC'")
	}
	if !client.inTransaction {
		return protocol.NewError("EXEC was called without MULTI")
	}
	queued := client.queue
	client.queue = nil
	client.inTransaction = false

	responses := make([]protocol.Value, 0, len(client.queue))
	for _, command := range queued {
		responses = append(responses, s.execute(client, command))
	}

	return protocol.NewArray(responses)
}
