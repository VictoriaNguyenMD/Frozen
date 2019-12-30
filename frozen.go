package main

import ( //Libraries to import
	"fmt" //Format library for println/printf
	"net" //Needed for creating the server. Listen(network, port)
	"log" //Prints error message even if the code breaks
	"bufio" //Input/Output using a buffer
	"strings" //Allow us to use string function
	"io" //Allows us to write to a particular user
)

type user struct {
	username string
	nickname string
	password string
	conn net.Conn
}

type connection struct {
	conn net.Conn //Maybe Connection ID?? Like TERMINAL A, B....
	server serv 
}

/*
** On a single server, there could be multiple channels. For example,
** think of the Slack application. There are multiple channels:
** - 42sv_global_random
** - 42sv_other_cantina, etc.
** Each channel can have a name, topic/description, and a differing
** amount of users
*/ 
type channel struct {
	name string
	description string
	users []user //This will have only the users in this channel
}


var global_users []user //This will store all user accounts on the server
var global_channels []channel //This will store all channels on the server

type serv struct {
	protocol string
	port string
	listener net.Listener //To later close this (defer)
}

/*
** Opens the server/location in which the clients will chat on.
** A server is a central computer that hosts data and other forms of resources
** A user will send and receive messages from the server 
** 
** Application:
** In terms of the chat room, think of it as User #1 wants to send a message
** to User #2. User #1 tells the server, "hey, I want to date User #2". Depending 
** on User #2's setting, the server could either deliver or keep the message. Thing of it
** as a warehouse/location/gateway
*/ 
func start_server(protocol string, port string) connection{
	fmt.Println("Launching server...")
	server, err := net.Listen(protocol, port) //Opens the connection. If it can't open the connection, err will have a random val that is not nil
	//net.Listen returns --> Listener, error
	var con connection

	if err != nil { //If there is an error with opening the "tcp" network with port value
		log.Println("Error:", err) //We use log instead of fmt because we want the message to print out even if the function breaks
		log.Println("There was an error starting the server...")
		return con
	}
	fmt.Println("Server Port", port)
	for { //This will keep on running until you kill or shut down the server. CMD + C
		c, err := server.Accept() //the server waits and accepts incoming connections/clients// Returns: net.Conn
		if err != nil {
			log.Println("Error:", err)
			log.Println("There was an error accepting...")
			return con
		}
		con = connection{conn: c} //If you're able to accept the connection properly, create a con "object" and assign c to conn
		con.server.listener = server //
		go handle_connection(con)
	}
}

func (c *connection) create_new_user(username *string){
	io.WriteString((*c).conn, "\nCreating a new account...\n")
	fmt.Println("\nCreating a new account...")
	io.WriteString((*c).conn, "Enter a username: ")
	usr, _ := bufio.NewReader((*c).conn).ReadString('\n')
	*username = usr[:len(usr)-1]

	io.WriteString((*c).conn, "Enter a nickname: ")
	nickname, _ := bufio.NewReader((*c).conn).ReadString('\n')
	nickname = nickname[:len(nickname)-1]

	io.WriteString((*c).conn, "Enter a password: ")
	password, _ := bufio.NewReader((*c).conn).ReadString('\n')
	password = password[:len(password)-1]

	new_user := user{username: *username, nickname: nickname, password: password, conn: c.conn}
	// fmt.Println("++++++++", c.conn)
	// fmt.Println("+++++++++++", new_user.conn)
	global_users = append(global_users, new_user)

	io.WriteString((*c).conn, "New account was successfully created...\n\n")

	fmt.Println("Username:", *username)
	fmt.Println("Nickname:", nickname)
	fmt.Println("Password:", password)
	fmt.Println("New account was successfully created...")
	fmt.Println("")

}

func (c connection) valid_user_pass(username, password string) int {
	for i := 0; i < len(global_users); i++ {
		if (global_users[i].username == username && global_users[i].password == password) {
			return 1
		}
	}
	return 0
}


func (c *connection) get_conn_id(username string) net.Conn {
	var recipient_conn_id net.Conn

	for i := 0; i < len(global_users); i++ {
		if (global_users[i].username == username) {
			recipient_conn_id = global_users[i].conn
			return recipient_conn_id 
		}
	}
	return (*c).conn
}

func get_username_from_conn(conn net.Conn) string {
	var curr_username string
	for i := 0; i < len(global_users); i++ {
		if (global_users[i].conn == conn) {
			curr_username = global_users[i].username
			return curr_username
		}
	}
	return curr_username
}

func (c *connection) do_cmd_pass_nick_user(username *string) {
	fmt.Println("Initial user authentication...")
	io.WriteString((*c).conn, "Enter your username: ")
	username_input, _ := bufio.NewReader((*c).conn).ReadString('\n')
	*username = username_input[:len(username_input)-1]

	io.WriteString((*c).conn, "Enter your password: ")
	password, _ := bufio.NewReader((*c).conn).ReadString('\n')
	password = password[:len(password)-1]

	if (c.valid_user_pass(*username, password) != 1) {
		io.WriteString((*c).conn, "\nIncorrect username or password...\n")
		c.create_new_user(username)
	} else {
		fmt.Printf("...successful.")
		io.WriteString((*c).conn, "Welcome back!\n\n")
	}
}

func (c *connection) do_cmd_nick(username string, conn_id net.Conn) {
	io.WriteString(conn_id, "Write a new nickname: ")
	new_nickname, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_nickname := new_nickname[:len(new_nickname)-1] //Removes extra leading and trailing white spaces

	for i := 0; i < len(global_users); i++ { //Changes all mentions of the username's nickname in each channel's users area
		if (global_users[i].username == username) {
			global_users[i].nickname = trimmed_nickname
			io.WriteString(conn_id, "...Nickname successfully changed...\n")
			io.WriteString(conn_id, "\n")
		}
	}
}

func (c *connection) do_cmd_list(conn_id net.Conn) {
	for i:= 0; i < len(global_channels); i++ {
		io.WriteString(conn_id, global_channels[i].name + "\n")
	}
	io.WriteString(conn_id, "\n")
}

func (c *connection) do_cmd_names(conn_id net.Conn) {
	for i:= 0; i < len(global_users); i++ {
		io.WriteString(conn_id, global_users[i].username + "\n")
	}
	io.WriteString(conn_id, "\n")
}

func remove(slice []user, username string) []user {
	for j:= 0; j < len(slice); j++ {
		if (slice[j].username == username) {
			return append(slice[:j], slice[j+1:]...) //We need the ... because we are appending another slice
		}
	}
	return slice
}

func (c *connection) do_cmd_part(username string) {
	conn_id := (*c).get_conn_id(username)

	io.WriteString(conn_id, "Channel to leave: ")
	channel, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_channel := channel[:len(channel)-1] //Removes extra leading and trailing white spaces
	found := false

	for i:= 0; i < len(global_channels); i++ {
		if (global_channels[i].name == trimmed_channel) {
			found = true
			global_channels[i].users = remove(global_channels[i].users, username)
		}
	}
	if (found == true) {
		io.WriteString(conn_id, "...Successfully left " + trimmed_channel + "\n\n")
	} else {
		io.WriteString(conn_id, "...User was not in " + trimmed_channel + "\n\n")
	}
}

func (c *connection) do_cmd_list_names(conn_id net.Conn) {
	fmt.Println(global_channels)
	for i:= 0; i < len(global_channels); i++ {
		io.WriteString(conn_id, "[" + global_channels[i].name + "]" + "\n")
		for j := 0 ; j < len(global_channels[i].users); j++ {
			io.WriteString(conn_id, global_channels[i].users[j].username + "\n")
		}
		io.WriteString(conn_id, "\n")
	}
}

func (c *connection) get_user(username string) (user, int){
	var ret_user user
	conn_id := (*c).get_conn_id(username)
	for i:= 0; i < len(global_users); i++ {
		if (global_users[i].username == username) {
			ret_user = global_users[i]
			return ret_user, 1
		}
	}
	io.WriteString(conn_id, "\n")
	return ret_user, 0
}

func (c *connection) append_new_channel(channel_name string, conn_id net.Conn) {
	var new_channel channel
	
	io.WriteString(conn_id, "...Channel does not exist. Creating channel...\n")
	io.WriteString(conn_id, "\nChannel name: " + channel_name + "\n")
	io.WriteString(conn_id, "Channel description: ")
	description, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_description := description[:len(description)-1]
	io.WriteString(conn_id, "\n")

	new_channel.name = channel_name
	new_channel.description = trimmed_description 
	global_channels = append(global_channels, new_channel)
}

func (c *connection) do_cmd_join(username string, trimmed_channel string) {
	conn_id := (*c).get_conn_id(username)
	found_channel := false
	for i:= 0; i < len(global_channels); i++ {
		if (global_channels[i].name == trimmed_channel) {
			found_channel = true
			ret_user, success := (*c).get_user(username)
			if (success == 1) {
				global_channels[i].users = append(global_channels[i].users, ret_user)
				io.WriteString(conn_id, "...Successfully joined " + trimmed_channel + "\n\n")
				return 
			}
		}
	}
	if (found_channel == false) {
		(*c).append_new_channel(trimmed_channel, conn_id)
		(*c).do_cmd_join(username, trimmed_channel)
	}
}

func (c *connection) do_cmd_privmsg_user() {
	io.WriteString(conn_id, "To user: ")
	username, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_username := username[:len(username)-1]

	io.WriteString(conn_id, "Message: ")
	message, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_message := message[:len(message)-1]

	current_username := get_username_from_conn(conn_id)
	recipient_conn_id := (*c).get_conn_id(trimmed_username)

	io.WriteString(conn_id, "\n")
	if (recipient_conn_id != nil) {

		io.WriteString(recipient_conn_id, "From user: " + current_username + "\n")
		io.WriteString(recipient_conn_id, "Message: " + trimmed_message + "\n")
		io.WriteString(recipient_conn_id, "\n")
	}
}

func (c *connection) do_cmd_privmsg_channel() {
	io.WriteString(conn_id, "To channel: ")
	channel, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_channel := channel[:len(channel)-1]
	found := 0

	for i:= 0; i < len(global_channels); i++ {
		if global_channels[i].name == trimmed_channel {
			found = 1		
		}
	}
	
	if found == 1 {
		io.WriteString(conn_id, "Message: ")
		message, _ := bufio.NewReader(conn_id).ReadString('\n')
		trimmed_message := message[:len(message)-1]
		for i:= 0; i < len(global_channels); i++ {
			for j := 0; j < len(global_channels[i].users); j++ {
				io.WriteString(global_channels[i].users[j].conn, "Channel Announcement: " + trimmed_message + "\n")
			}
		}
		io.WriteString(conn_id, "\n")
	}
}

func (c *connection) do_cmd_privmsg(conn_id net.Conn) {
	io.WriteString(conn_id, "User or Channel? ")
	decision, _ := bufio.NewReader(conn_id).ReadString('\n')
	trimmed_decision := decision[:len(decision)-1]
	io.WriteString(conn_id, "\n")

	if trimmed_decision == "User" {
		c.do_cmd_privmsg_user()
	} else if trimmed_decision == "Channel" {
		c.do_cmd_privmsg_channel()
	} else {
		io.WriteString(conn_id, "\nInvalid option.\n\n")
	}
}

var conn_id net.Conn
var commands = [10]string{"PASS NICK USER", "NICK", "JOIN", "PART", "NAMES", "LIST", "LIST_NAMES", "PRIVMSG", "MY_INFO", "HELP"}

func is_command(message string) int {
	for i:= 0; i < len(commands); i++ {
		if message == commands[i] {
			return 1
		}
	}
	return 0
}

func (c *connection) check_command(message string, username *string) {
	trimmed_message := strings.TrimSpace(message) //Remove any leading and trailing whitespaces from messages
	if (trimmed_message != "PASS NICK USER" && *username == "") {
		io.WriteString((*c).conn, "You are not logged in. Please type \"PASS NICK USER\" to log in.\n")
	} else  if (is_command(trimmed_message) == 1){
		switch trimmed_message {
		case "PASS NICK USER": //Initial authentication for the user
			io.WriteString((*c).conn, "\n->Initial authentication for the user...\n")
			(*c).do_cmd_pass_nick_user(username)
			fmt.Println(*username)
		case "NICK": //Change nickname
			conn_id = (*c).get_conn_id(*username) //To make sure we have the same conn id based on username
			fmt.Println(conn_id, "<->", (*c).conn)
			io.WriteString(conn_id, "\n->Change nickname...\n")
			(*c).do_cmd_nick(*username, conn_id)
		case "JOIN": //Make user join a channel
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Joining a channel...\n")
			io.WriteString(conn_id, "Channel to join: ")
			channel, _ := bufio.NewReader(conn_id).ReadString('\n')
			// io.WriteString(conn_id, "\nSuccesfully read user input" + channel + "\n")
			trimmed_channel := channel[:len(channel)-1] //Removes extra leading and trailing white spaces
			(*c).do_cmd_join(*username, trimmed_channel)
		case "PART": //Make user leave a channel
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n-> " + *username + " leaving a channel...\n")
			(*c).do_cmd_part(*username)
		case "NAMES": //List all users connected to the server
			// fmt.Println(*c)
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Listing all users in the server...\n")
			(*c).do_cmd_names(conn_id)
		case "LIST": //List all channels in the server
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Listing all channels in the server...\n")
			(*c).do_cmd_list(conn_id)
		case "LIST_NAMES": //List all users in the channels
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Listing all users in each channel...\n")
			(*c).do_cmd_list_names(conn_id)
		case "PRIVMSG": //Send a message to another user or channel
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Sending a message to another user/channel...\n")
			(*c).do_cmd_privmsg(conn_id)
		case "MY_INFO": //Get my info
			conn_id = (*c).get_conn_id(*username)
			io.WriteString(conn_id, "\n->Getting my information...\n")
			io.WriteString(conn_id, "My username: " + *username + "\n")
			fmt.Fprintln(conn_id, "My conn_id: ", conn_id)
			io.WriteString(conn_id, "\n")
		case "HELP": //provides a list of commands
			input_instructions(*c)
		}
		fmt.Println("Global Channels: ", global_channels)
		fmt.Println("Global Users: ", global_users)
		fmt.Println("")
	} else {
		io.WriteString((*c).conn, *username + ": " + trimmed_message + "\n")
	}
}

func test_case(con *connection) {
	tinder := channel{name: "Tinder", description: "Tinder_Description"}
	east := channel{name: "East", description: "East_Description"}
	global_channels = append(global_channels, tinder)
	global_channels = append(global_channels, east)
	u1 := user{username: "Femi", nickname: "Fem", password: "1234"}
	u2 := user{username: "Yeonuk", nickname: "Roomie", password: "0000"}
	global_channels[0].users = append(global_channels[0].users, u1)
	global_channels[0].users = append(global_channels[0].users, u2)
	global_users = append(global_users, u1)
	global_users = append(global_users, u2)

	// fmt.Println(con)
	// fmt.Println(con.server)
	// fmt.Println(con.server.channels)
	// fmt.Println(con.server.channels[0])
	// fmt.Printf("\n")
	// fmt.Println(u1)
	// fmt.Println(u2)
}

func input_instructions(c connection) {
	io.WriteString(c.conn, "\n--------Frozen Rush --------------------------\n")
	io.WriteString(c.conn, "\nBasic Input Instructions\n")
	io.WriteString(c.conn, "1. Login using: PASS NICK USER\n")
	io.WriteString(c.conn, "2. Perform another command\n")

	io.WriteString(c.conn, "\nCommands:\n")
	io.WriteString(c.conn, "	PASS NICK USER: login/create account\n")
	io.WriteString(c.conn, "	NICK: change nickname\n")
	io.WriteString(c.conn, "	JOIN: join a channel\n")
	io.WriteString(c.conn, "	PART: leave a channel\n")
	io.WriteString(c.conn, "	NAMES: show all users connected to the server\n")
	io.WriteString(c.conn, "	LIST: show all channels in the server\n")
	io.WriteString(c.conn, "	LIST_NAMES: show all users in each channel in the server\n")
	io.WriteString(c.conn, "	PRIVMSG: send a message to another user/channel\n")
	io.WriteString(c.conn, "	MY_INFO: show username and conn_id\n")

	io.WriteString(c.conn, "\nAdditional Commands:\n")
	io.WriteString(c.conn, "	HELP: provides a list of server commands\n\n")
}

func handle_connection(con connection) {
	defer con.conn.Close() //This will defer closing the Conn (aka Terminal) until everything else runs in this function
	input_instructions(con) //Runs a function that shows the commands for the server
	username := ""
	for {
		message, err := bufio.NewReader(con.conn).ReadString('\n') //The return for ReadString is string, err. 
		if err != nil {
			log.Println("Error:", err)
			log.Println("Data does not end with a new line...")
			return 
		}
		message = message[:len(message)-1]
		fmt.Println("\nMessage Received:", message)
		fmt.Println("")
		con.check_command(message, &username)		
	}
}

func stop_server(server net.Listener) {
	global_channels = nil
	global_users = nil
	if (server != nil) {
		server.Close()
	}
}

func show_users_info(chan_id channel) {
	fmt.Println("Chat Name:", chan_id.name)
	fmt.Println("Description:", chan_id.description)
	fmt.Println("")
	for i := 0; i < len(chan_id.users); i++ {
		fmt.Printf("[User #%d]\n", i+1);
		fmt.Println("Username:", chan_id.users[i].username)
		fmt.Println("Nickname:", chan_id.users[i].nickname)
		fmt.Println("Password:", chan_id.users[i].password)
		fmt.Println("")
	}
}

func main() {
	s := serv{protocol: "tcp", port: ":6667"}
	c := start_server(s.protocol, s.port)
	defer stop_server(c.server.listener) //Will execute only after the surrounding function returns
}