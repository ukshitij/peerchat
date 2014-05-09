package dht

import "testing"
import "runtime"
import "fmt"
import "time"
import "math"
import "strconv"

// Signal failures with the following:
// t.Fatalf("error message here")

const localIp = "127.0.0.1"

func TestBasic(t *testing.T) {
	/*
		TestBasic:
		1) Starts two nodes
		2) Introduces node1 to node2
		3) Nodes send messages
		
		We verify the messages are not lost
		and arrive unaltered. 
	*/
	runtime.GOMAXPROCS(4)

	port1 := ":4444"
	port2 := ":5555"
	username1 := "Alice"
	username2 := "Frans"

	// user1 starts the Peerchat network, and
	// user2 joins by bootstrapping
	user1 := Register(username1, localIp + port1, "")
	time.Sleep(time.Millisecond * 50)
	user2 := Register(username2, localIp + port2, localIp + port1)

	time.Sleep(time.Millisecond * 50)

	// tests that we can find both users!
	u1_ip := user2.node.FindUser(username1)
	assertEqual(t, u1_ip, localIp+port1)
	u2_ip := user1.node.FindUser(username2)
	assertEqual(t, u2_ip, localIp+port2)
	
	// users exchange messages
	user1.SendMessage(username2, "Hi Frans! Wanna play squash?")
	time.Sleep(time.Second * 1)
	user2.SendMessage(username1, "Sure Alice, what time?")
	
	// kill user nodes
	user1.node.Dead <- true
	user2.node.Dead <- true
}


func registerMany(num_users int) map[string]*User{
	users := make(map[string]*User)

	bootstrap := ""

	for i :=0; i < num_users; i++{
		username := strconv.Itoa(i)
		ipAddr := localIp + ":" + strconv.Itoa(i + 7000)
		user := Register(username, ipAddr, bootstrap)
		bootstrap = localIp + ":" + strconv.Itoa(i + 7000)
		time.Sleep(time.Millisecond * 5)
		users[username] = user
	}

	return users

}

func TestManyRegistrations(t *testing.T) {
	
	users := registerMany(40)
	time.Sleep(time.Second)
	for _, user := range users{
		user.node.AnnounceUser(user.name, user.node.IpAddr)
	}
	time.Sleep(time.Second)
	for _, user := range users{
		fmt.Println(user.name, user.node.kv)
		for targetUsername, targetUser := range users{
			targetIp := user.node.FindUser(targetUsername)
			assertEqual(t, targetIp, targetUser.node.IpAddr)
			fmt.Println("Correct")
		}
	}
	
	for _, user := range users {
		user.node.Dead <- true
	}
}


func assertEqual(t *testing.T, out, ans interface{}) {
    if out != ans {
        t.Fatalf("wanted %v, got %v", ans, out)
    }
}

func isEqualRE(entry1 []RoutingEntry, entry2 []RoutingEntry) bool{
	if len(entry1) != len(entry2){
		return false
	}
	for i, v := range entry1{
		if v != entry2[i] {
			return false
		}
	}
	return true
}

func TestCommonUnit(t *testing.T) {
    //common unit tests

    //Sha1 Test
    assertEqual(t, Sha1("abc"), Sha1("abc"))
    if Sha1("fjkels") == Sha1("qwewqi") {
        t.Fatalf("Sha1 collision")
    }
    //reference Sha1 computed at www.sha1-online.com and lowest 8 bytes converted to decimal
    assertEqual(t, Sha1("Forrest"), ID(10556789446649181072))
    assertEqual(t, Sha1("testing testing 123"), ID(16871972680281001427))

    //find_n
    a := ID(0)
    b := ID(1)
    c := ID(math.MaxUint64)
    d := ID(1 << 15)
    assertEqual(t, find_n(a, b), uint(63))
    assertEqual(t, find_n(a, c), uint(0))
    assertEqual(t, find_n(a, d), uint(48))
}

func TestDhtNodeUnit(t *testing.T) {
    //DhtNode Unit Tests

    //moveToEnd Test
    id := Sha1("hi")
    in0 := []RoutingEntry{RoutingEntry{"a", id}, RoutingEntry{"b", id}, RoutingEntry{"c", id}}
    ans1 := []RoutingEntry{in0[1], in0[2], in0[0]}  //moveToEnd(in1, 0)
    ans2 := []RoutingEntry{in0[0], in0[2], in0[1]}  //moveToEnd(in1, 1)
    out1 := make([]RoutingEntry, 3)
    out2 := make([]RoutingEntry, 3)
    out3 := make([]RoutingEntry, 3)
    copy(out1, in0)
    copy(out2, in0)
    copy(out3, in0)
    moveToEnd(out1, 0)
    if !isEqualRE(ans1, out1) {
        t.Fatalf("wanted %v, got %v", ans1, out1)
    }
    moveToEnd(out2, 1)
    if !isEqualRE(ans2, out2) {
        t.Fatalf("wanted %v, got %v", ans2, out2)
    }
    moveToEnd(out3, 2)
    if !isEqualRE(in0, out3) {
        t.Fatalf("wanted %v, got %v, in0, out3")
    }
    //
    fmt.Println()


}
