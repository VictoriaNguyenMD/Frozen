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
**PRIVMSG:** Send a message to another user or a channel

## Breakdown: Analogy

### Disclaimer
This is a very simplified version of a basic IRC Server. Because the "supported commands" were interpreted as stand-alone words that need to be on a single line, no parsing for the user input was needed for this project. Some other users, such as `https://github.com/sayakura/Go-IRC-server` and `https://github.com/jjaniec/Frozen/blob/dev/README.md`, had differing interpretations of the project's pdf.
