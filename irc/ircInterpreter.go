package irc

import (
	"log"
	"ollyster/files"
	"regexp"
	"strings"
)

func (this *IrcServer) IrcInterpreter(message string) {

	if matches, _ := regexp.MatchString("(?i)^PING :.*$", message); matches == true {
		log.Printf("[IRC] %s ", message)
		response := strings.Replace(message, "PING", "PONG", 1)
		log.Printf("[IRC] Sending back the -> %s", response)
		this.IrcCmd(response)

		return
	}

	// :nick!user@ip-address PRIVMSG your-nick :VERSION
	if matches, _ := regexp.MatchString("(?i)^:.*PRIVMSG.*VERSION$", message); matches == true {
		sinta := strings.Split(message, "!")
		log.Printf("[IRC] VERSION REQUEST from %s ", sinta[0])
		version := "NOTICE " + strings.TrimLeft(sinta[0], ":") + " : VERSION Ollyster DEV https://github.com/uriel-fanelli/ollyster"
		this.socket.Write([]byte(version))
		log.Printf("[IRC] Sending back the -> %s", version)
		return
	}

	// :nick!user@ip-address PRIVMSG your-nick :Message
	privMsgString := "(?i)^:.*!.*PRIVMSG.*" + this.nickname + " :.*$"
	if matches, _ := regexp.MatchString(privMsgString, message); matches == true {

		sinta := strings.Split(message, "!")
		sender := strings.TrimLeft(sinta[0], ":")
		// sender contains the sender nick

		field := strings.Split(message, "PRIVMSG")
		torn := strings.Split(field[1], ":")
		msg := strings.Join(torn[1:], ":")
		// msg contains the message after the 1st colon

		files.MyStream.WriteMsgPriv(sender, msg)

		log.Printf("[IRC] Private message from %s:  <%s>", sender, msg)
		return
	}

	// :nick!user@ip-address PRIVMSG #channel :Message
	chanMsgString := "(?i)^:.*!.*PRIVMSG.*" + this.channel + " :.*" + this.nickname + ".*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		sinta := strings.Split(message, "!")
		sender := strings.TrimLeft(sinta[0], ":")
		// sender contains the sender nick

		field := strings.Split(message, "PRIVMSG")
		torn := strings.Split(field[1], ":")
		msg := strings.Join(torn[1:], ":")
		// msg contains the message after the 1st colon

		log.Printf("[IRC] %s sent a message to %s:  <%s>", sender, this.channel, msg)
		files.MyStream.WriteMsgMention(sender, this.channel, msg)
		return
	}

	// :nick!user@ip-address PRIVMSG #channel :Message
	chanMsgString = "(?i)^:.*!.*PRIVMSG.*" + this.channel + " :.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		sinta := strings.Split(message, "!")
		sender := strings.TrimLeft(sinta[0], ":")
		// sender contains the sender nick

		field := strings.Split(message, "PRIVMSG")
		torn := strings.Split(field[1], ":")
		msg := strings.Join(torn[1:], ":")
		// msg contains the message after the 1st colon

		log.Printf("[IRC] %s sent a message to %s:  <%s>", sender, this.channel, msg)
		files.MyStream.WriteMsgGroup(sender, this.channel, msg)
		return
	}

	// :sinisalo.freenode.net 322 Ollyster #pld-git 3 :https://www.pld-linux.org/pld-git
	chanMsgString = "(?i)^:.*322 " + this.nickname + " #.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		field := strings.Split(message, " 322 "+this.nickname+" ") // field[0] is useless , field[1] contains the payload
		list := strings.Split(field[1], " ")                       // now list[0] contains the channel name, list[1] the number of members, list[2:] everything else
		desc := strings.TrimLeft(strings.Join(list[2:], " "), ":") // this is "everything else, rebuilt with spaces and then removed the ":"

		list_line := "<tr><td class=\"col-md-2\"><b>" + list[0] + "</b></td><td class=\"col-md-4\">" + desc + "</td></tr>"

		log.Printf("[IRC] Registering channel :  <%s>", list[0])
		files.MyStream.AddUniqueChannel(list_line)
		return
	}

	// :sinisalo.freenode.net 323 Ollyster :End of /LIST
	chanMsgString = "(?i)^:.*323 " + this.nickname + " :End.*LIST$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		log.Printf("[IRC] Saving list of channels....")
		files.MyStream.FlushChanList()
		log.Printf("[IRC] Channel list registered")
		return
	}

}
