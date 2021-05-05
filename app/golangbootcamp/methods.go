package main

import "fmt"

type User struct {
    Id             int
    Name, Location string
}

func (u *User) Greetings() string {
    return fmt.Sprintf("Hi %s from %s",
        u.Name, u.Location)
}

type Player struct {
    *User
    GameId int
}


func NewPlayer(id int, name, location string, gameId int) *Player {
    return &Player{
        User: &User{id, name, location},
        GameId: gameId,
    }
}
func main() {
    p := NewPlayer(42, "Matt", "LA", 90404)
    fmt.Println(p.Greetings())
}
