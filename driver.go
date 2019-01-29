package main

import (
  "fmt"
  "os"
  "bufio"
  "bytes"
  "io/ioutil"
  "encoding/json"
)

//USER DATASTRUCTURES
//key userid, value is JSON object (
type Storage struct {
  IdMap map[string]interface{}
}

//key userid, value is slice of Storage obj
type Connection struct {
  SocialGraph map[string][]Storage
}

var E = Event{}
var User1 = Storage{}
var User2 = Storage{}
var Bacon = Connection{}

//JSON DATATYPES
type Event struct {
  To struct {
    Id string `json:"id"`
    Name string `json:"name"`
  } `json:"to"`
  From struct {
    Id string `json:"id"`
    Name string `json:"name"`
  } `json:"from"`
  Timestamp int `json:"timestamp"`
  AreFriends bool `json:"areFriends"`
}

type Back struct {
  Id string `json:"id"`
  Name string `json:"name"`
}

//USER FUNCTIONS
func InitMap(u * Storage) {
  u.IdMap = make(map[string]interface{})
}

func InitStores() []Storage{
  return make([]Storage,1,1)
}
func InitGraph(c * Connection) {
  c.SocialGraph = make(map[string][]Storage)
}

//USER METHODS
func (s *Storage) Place(id string, e Event) {
  if(id == e.To.Id) { //recieving friend or unfriend req
    _, ok := s.IdMap[e.To.Id]

    if e.AreFriends && !ok {
      //fmt.Println(e.To.Name, " recieve req from ",e.From.Id)
      s.IdMap[e.From.Id] = e.From
    } else if !e.AreFriends && ok {
      //fmt.Println(e.To.Name, " recieve unfriend req from ",e.From.Name)
      delete(s.IdMap, e.From.Id)
    }
  } else if(id == e.From.Id) { //sending friend or unfriend req
    _, ok := s.IdMap[e.From.Id]

    if e.AreFriends && !ok {
      //fmt.Println(e.From.Name, " sent req to ", e.To.Id)
      s.IdMap[e.To.Id] = e.To
    } else if !e.AreFriends && ok {
      //fmt.Println(e.From.Name, " sent unfriend req to ", e.To.Name)
      delete(s.IdMap, e.To.Id)
    }
  }
}

func (c * Connection) AddPrev(front, k string) {
  c.SocialGraph[k] = append(c.SocialGraph[k], Storage{})
  InitMap(&c.SocialGraph[k][1])
  c.SocialGraph[k][1].IdMap[k] = c.SocialGraph[k][0].IdMap[front]
}

type DirectFriends  struct {
  Friends []interface{} `json:"friends"`
}

type MutualFriends struct {
  Mutuals []interface{} `json:"mutuals"`
}

type DegreeFriends struct {
  Degrees []int `json:"degrees"`
}


//FUNCTIONS TO WRITE JSON FILES
func WriteFriends(s *Storage) {
  var df DirectFriends
  df.Friends = make([]interface{}, len(s.IdMap), len(s.IdMap))

  i := 0
  for _, v := range s.IdMap {
    df.Friends[i] = v
    i++
  }

  directJson, _ := json.Marshal(df)
  ioutil.WriteFile("friends.json", directJson, 0644)
}

func WriteMutuals(s1 *Storage, s2 *Storage) {
  var mutf MutualFriends
  mutf.Mutuals = make([]interface{},0)

  for k, v := range s1.IdMap {
    _, ok := s2.IdMap[k]
    if ok {
      mutf.Mutuals = append(mutf.Mutuals,v)
    }
  }

  mutualJson, _ := json.Marshal(mutf)
  ioutil.WriteFile("mutuals.json", mutualJson, 0644)
}


func WriteDegrees(conn *Connection, src string, dest string) {
  var degf DegreeFriends
  var front string
  isPath := false
  degrees := 0

  //fifo queue implemented using go splice
  q := []string{src}

  //check to see if a path is even possible
  var slen, dlen, exists bool
  _, oks := conn.SocialGraph[src]
  _, okd := conn.SocialGraph[dest]
  if oks && okd {
    slen = len(conn.SocialGraph[src][0].IdMap) > 0
    dlen = len(conn.SocialGraph[dest][0].IdMap) > 0
  }

  exists = (oks && okd && slen && dlen)

  //set previous for src to itself
  if exists {
    conn.AddPrev(src,src)
  }

  //depth first search across social graph from src to dest
  for len(q) > 0 && exists {
    //pop
    front, q =  q[0], q[1:]
    m := conn.SocialGraph[front][0]

    for k, _ := range m.IdMap {
      //push
      if k == dest {
        //clear queue and exit loops
        q = q[0:0]
        isPath = true
        //rewind

        conn.AddPrev(front, k)

        curr := k
        //backtrack path from dest to src and count degrees of sep
        for curr != src {
          var buf bytes.Buffer
          var b Back

          enc := json.NewEncoder(&buf)
          if err := enc.Encode(conn.SocialGraph[curr][1].IdMap[curr]); err != nil {
           fmt.Println(err.Error())
          }
          if err := json.Unmarshal(buf.Bytes(), &b); err != nil {
           fmt.Println(err.Error())
          }
          curr = b.Id

          degrees++
        }

        break
      } else if len(conn.SocialGraph[k]) < 2 {
        //append friend's data to slice for purpose of both
        //visited flag and mem for backtrack
        conn.AddPrev(front,k)

        //insert key into back of queue
        q = append(q,k)
      }
    }
    //delete(conn.SocialGraph, front)
  }

  degf.Degrees = make([]int,1,1)

  if degf.Degrees[0] = degrees; !isPath {
    degf.Degrees[0] = -1 //no intersection
  }

  degreeJson, _ := json.Marshal(degf)
  ioutil.WriteFile("bacon.json", degreeJson, 0644)
}

func LoadEvents(args []string) {
  file, err := os.Open(args[0])

  defer file.Close()

  if err != nil {
    fmt.Println(err.Error())
    fmt.Println("Unreadable File")
    file.Close()

    os.Exit(3)
    return
  }

  s := bufio.NewScanner(file)

  //stores json elems
  var uid1 = User1
  var uid2 = User2
  var uidBacon = Bacon

  //init maps for use
  InitMap(&uid1)
  if len(args) >= 3 {
    InitMap(&uid2)
  }
  if len(args) == 4 && args[3] == "bacon" {
    InitGraph(&uidBacon)
  }

  //scan jsons and input data into maps 
  for(s.Scan()) {
    //stores current json event
    var v = E

    if errr := json.Unmarshal(s.Bytes(), &v); errr != nil {
      fmt.Println(errr.Error())
      panic(errr)
    }

    switch {
    case len(args) < 2 :
      fmt.Println("two few arguements")
      file.Close()
      os.Exit(3)

    //friends of one user
    case len(args) < 3 :
      uid1.Place(args[1],v)

    //mutual friends b/w two users
    case len(args) < 4 :
      uid1.Place(args[1],v)
      uid2.Place(args[2],v)

    //degrees of sep b/w two users
    case len(args) < 5 && args[3] == "bacon" :
      if _, ok := uidBacon.SocialGraph[v.To.Id]; !ok {
        uidBacon.SocialGraph[v.To.Id] = InitStores()
        InitMap(&uidBacon.SocialGraph[v.To.Id][0])
      }
      if _, ok := uidBacon.SocialGraph[v.From.Id]; !ok {
        uidBacon.SocialGraph[v.From.Id] = InitStores()
        InitMap(&uidBacon.SocialGraph[v.From.Id][0])
      }
      uidBacon.SocialGraph[v.To.Id][0].Place(v.To.Id,v)
      uidBacon.SocialGraph[v.From.Id][0].Place(v.From.Id,v)
      //case bacon

    default:
      fmt.Println("Please enter the correct arguements")
      file.Close()
      os.Exit(3)
    }//switch
  }

  //write data to new json file
  switch {
    case len(args) < 3 :
      parse.WriteFriends(&uid1)
      fmt.Println("returning friends")
    case len(args) < 4 :
      parse.WriteMutuals(&uid1,&uid2)
      fmt.Println("returning mutual friends")
    case len(args) < 5 && args[3] == "bacon" :
      parse.WriteDegrees(&uidBacon,args[1],args[2])
      fmt.Println("returning bacon")
    default:
  }

  //return data
}

func main() {
  //grab arguements
  addargs := os.Args[1:]

  //run program
  LoadEvents(addargs)
}
