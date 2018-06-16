package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/wybiral/shellgame/levels"
	"github.com/wybiral/shellgame/pkg/container"
	"io"
	"log"
)

// Handler for new SSH sessions.
func sessionHandler(s ssh.Session) {
	ctx := s.Context()
	// Get level from context
	level, ok := ctx.Value("level").(levels.Level)
	if !ok {
		log.Println("context missing level")
		return
	}
	// Get container client
	cl, err := container.NewClient()
	if err != nil {
		log.Println(err)
		return
	}
	// Create container for level image
	c, err := cl.Create(ctx, level.Image)
	if err != nil {
		log.Println(err)
		return
	}
	// Attach to container for IO
	conn, err := c.Attach(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	// Cleanup
	defer func() {
		conn.Close()
		c.Kill(ctx)
		c.Remove(ctx)
	}()
	// Start container
	err = c.Start(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	// Pump IO between SSH session and container
	go func() {
		io.Copy(conn, s)
	}()
	io.Copy(s, conn)
}

func main() {
	addr := "127.0.0.1:2222"
	levels := levels.GetAll()
	s := &ssh.Server{
		Addr: addr,
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			ctx.SetValue("password", password)
			user := ctx.User()
			log.Println(user, password)
			level, ok := levels[user]
			if !ok {
				return false
			}
			ctx.SetValue("level", level)
			return password == level.Password
		},
	}
	s.Handle(sessionHandler)
	log.Println("SSH server listening at ", addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
