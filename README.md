# Frozen

## Running the Program

### Starting the serve:
```bash
git clone https://github.com/VictoriaNguyenMD/Frozen.git frozen
cd frozen
go run frozen.go
```

### Adding users to the server:
For each user, open a new terminal.
```bash
nc localhost 6667
```

## Supported Commands
**PASS NICK USER:** Initial authentication for the user 

**NICK:** Change user nickname

**JOIN:** Join a channel

**PART:** Leave a channel

**NAMES:** List all users connected to the server

**LIST:** List all channels in the server

**LIST_NAMES:** List the users in each channel

**PRIVMSG:** Send a message to another user or a channel

**HELP:** List the commands the server supports

## Implementation

### Variables
Two global variables were used, `global_users` and `global_channel`, to store all the users who had created an account on the server and what channels exist on the server. We want to store all the accounts created on the server so that the user is able to log back into their account. In addition, by having the two global arrays `global_users` and `global_channel`, it is easier to later use the `NAMES` and `LIST` command to list all the users and channels on the server. To be more precise, it would be better to note if the user is currently connected or not connected (aka logged-out) to the server. 

A `user struct` was created to store all the user information. The `user struct` contains a person's username, nickname, password, and connection ID. To mimic the creation of a user joining the server, you would need to open a new terminal. As a result, each terminal is a new user. Therefore, to distinguish between the users, a unique terminal ID (`conn net.Conn`) will be assigned to each user in the `user struct`. 

A `channel struct` was created to store all the information pertaining to a channel. The `channel struct` will have the channel's name, a description of the channel, and the users who have joined the channel. It is important to keep track of the users in a particular channel to potentially note how many users the channel has and be able to write a message to the channel so all the other users could see the message. By having the `users` array within the `channel struct`, it is easier to iterate through a specific channel and write a message to each user within the array. Another option is to have an array of channels a particular used joined in the `user struct`. 

A `server struct` was created to store information relating to the server; however, this is not necessary since there is only one server. The port number for the server will be 6667 since IRC typically use transmission control protocol (TCP), a type of acknowledgement system in which the users need to provide information of receiving the message (eg. "Read at 1:43 pm"), as its transport protocol. The TCP port for IRC traffic is 6667.

### Logical Flow
The `tcp` protocol and port number `:6667` was passed into `start_server()` function where `net.Listen()` starts the server. Once the server is running, `server.Accept()` will listen to any incoming users or people who try to login to the server. When the server was accessed by a user, `handle_connection()` function runs to handle the new user and output the server's commands on the user's terminal. A unique connection ID is assigned to the user when the user opens the server. As long as the programming is running, the server will listen into the user's terminal for any written messages/input. When an input is received, the `check_command()` function will run to determine if any of the input is one of the server's pre-defined commands (eg. PASS NICK USER, NICK, etc.). If any of the commands match via checking if the word is in a global string array `commands`, then a respective function will run:

**`PASS NICK USER:`** The server prompts the user to login with their username/password. For simplication purposes, only one chance if given. If the user has an incorrect username/password combination or does not have an account, the server prompts the user to create a new account and type in their username, nickname, and password. The user struct is stored in the `global_users` variable. No other commands will run until the user either logs into their account or creates a new account.

**`NICK:`** The server prompts the user to change the user's nickname. If the user's username matches with a user's username in `global_users`, then the user's nickname will be changed. The username was used to determine matching because the username cannot be changed. Another alternative is to look at the user's specific connection ID.

**`JOIN:`** The server will prompt the user to join a channel. If the channel does not exist in `global_channels`, then the a new channel will be created and the server will prompt the user for the channel's description. If the channel exist, then the user will be appended to the the channel's `user` array.

**`PART:`** The server will prompt the user to indicate which channel to leave. If the channel is in `global_channels`, then the code will slice the `users` array to remove the user and reassign the `global_channels`'s `users` array to the new slice. This is performed via `slice[:j], slice[j+1:]...`. 

**`NAMES:`** The code will iterate through `global_users` and print out the name of each user that created an account in the server contains.

**`LIST:`** The code will iterate through `global_channels` and print out the name of each channel the server contains. The code will still list the channel even if there are no users in the channel.

**`LIST_NAMES:`** The code will iterate through each user in each channel in `global channels` and print out the channel name and each user within that channel. This is perofrmed via two `for loops` to iterate through the `global channels` array and the `users` array within each channel.

**`PRIVMSG:`** The server will prompt the user to either send a global message to the channel or to a send a message to a specific user. If the user chooses global message, the code will iterate through all the users in the `global_channel` that matches the channel name, and display a channel announcement according to the sender's input. If the user chooses a private message, the recipient's specific connection ID will be obtained fron `get_conn_id()` where the `global_users` array will be searched for a matching username. A message will then be displayed on the recipient's terminal according to the sender's input.

**`MY_INFO:`** The server will display the username and the unique connection ID. The unique connection ID was previously specified when the `server.Accept()` listens for a new user in `start_server()` function.

**`HELP:`** The server will display instructions to login and the available commands for users to use.

## Analogy
To better understand how to do this project, think of Facebook as an Internet Relay Chat (IRC) system. 

1. You need to first make sure that Facebook is launched or even exists on the internet. For example, when we type www.Facebook.com, the url will show the Facebook page instead of an error page. Similarly, to launch our chat system, we type `go run frozen.go` to make sure that the server exists on our computer.

2. Now, we want to go to Facebook, so we need to connect to the website server by typing in the url, www.Facebook.com. In terms of our chat system, to connect to the server, we will type `nc localhost 6667`.

3. Once we are on the chat server, you need to login/create an account, similarly to how you would for Facebook.

4. You now how the option of changing your nickname, joining a channel (aka a chat room or Facebook group), leaving a channel, seeing who else signed up, seeing what channels there are, and sending a private message to another user or a global announcement to another channel.

#### Disclaimer
This is a very simplified version of a basic IRC Server that is not RFC compliant and did not use security protection. Because the "supported commands" were interpreted as stand-alone words that need to be on a single line, no parsing for the user input was needed for this project. Some other users, such as `https://github.com/sayakura/Go-IRC-server` and `https://github.com/jjaniec/Frozen/`, had differing interpretations of the project's pdf. The above IRC does not fulfill all of the grading requirements due to time constraints and the vagueness of the original pdf. If I were to re-code this, I would use golang's `chan` feature rather than global variables.
