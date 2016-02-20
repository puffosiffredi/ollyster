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

	// :nick!user@ip-address PRIVMSG your-nick :Message
	// PRIVATE MESSAGE
	privMsgString := "(?i)^:.*!.*PRIVMSG.*" + this.nickname + " :.*$"
	if matches, _ := regexp.MatchString(privMsgString, message); matches == true {

		re, _ := regexp.Compile("(?i)^:(.*)!.*[ ]+PRIVMSG[ ]+.*[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)

		files.MyStream.WriteMsgPriv(match[1], match[2])

		log.Printf("[IRC] Private message from %s:  <%s>", match[1], match[2])
		return
	}

	// :nick!user@ip-address PRIVMSG #channel :Message
	// MENTION
	chanMsgString := "(?i)^:.*!.*PRIVMSG[ ]+ #.*[ ]+:.*" + this.nickname + ".*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		re, _ := regexp.Compile("(?i)^:(.*)!.*[ ]+PRIVMSG[ ]+(#.*)[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)

		log.Printf("[IRC] %s Mentioned you in %s:  <%s>", match[1], match[2], match[3])
		files.MyStream.WriteMsgMention(match[1], match[2], match[3])
		return
	}

	// :nick!user@ip-address PRIVMSG #channel :Message
	// MESSAGE ONLY
	chanMsgString = "(?i)^:.*!.*[ ]+PRIVMSG[ ]+#.*[ ]+:.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		re, _ := regexp.Compile("(?i)^:(.*)!.*[ ]+PRIVMSG[ ]+(#.*)[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)

		log.Printf("[IRC] %s sent a message to %s:  <%s>", match[1], match[2], match[3])
		files.MyStream.WriteMsgGroup(match[1], match[2], match[3])
		return
	}

	// :sinisalo.freenode.net 322 Ollyster #pld-git 3 :https://www.pld-linux.org/pld-git
	// CHANNEL LIST ITEM
	chanMsgString = "(?i)^:.*[ ]+322[ ]+" + this.nickname + "[ ]+#.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		re, _ := regexp.Compile("(?i)^:.*[ ]+322[ ]+.*[ ]+(#.*)[ ]+[0-9]+[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)
		list_line := "<tr><td class=\"col-md-2\"><b>" + match[1] + "</b></td><td class=\"col-md-4\">" + match[2] + "</td></tr>"

		log.Printf("[IRC] Registering channel :  <%s>", match[1])
		files.MyStream.AddUniqueChannel(list_line)
		return
	}

	// :Loweel!~loweel@p2003004C6815B300D25099FFFE17D56C.dip0.t-ipconnect.de NOTICE Ollyster :Notizione ma di quelli con le palle
	// NOTICE
	chanMsgString = "(?i)^:.*!.*NOTICE[ ]+" + this.nickname + "[ ]+:.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {

		re, _ := regexp.Compile("(?i)^:(.*)!.*NOTICE[ ]+.*[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)

		log.Printf("[IRC] %s sent a NOTICE :  <%s>", match[1], match[2])
		files.MyStream.WriteNotice(match[1], match[2])
		return
	}

	// :sinisalo.freenode.net 353 Ollyster = #social :gregoriosw_vp nullwarp asumu xmpp-gnu msava arctanx k0nsl cha_ron ascarpino jerrykan Sazius chimo dualbus n4mu cow_2001 atari-frosch molgrum alanz Stig_Atle BeS vinzv pztrn rec0de AlexanderS @ChanServ tonnerkiller kromonos nobody rolfrb
	// LIST OF USERS
	chanMsgString = "(?i)^:.*353[ ]+" + this.nickname + "[ ]+.[ ]+#.*[ ]+:.*$"
	if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {
		re, _ := regexp.Compile("(?i)^:.*353[ ]+(.*)[ ]+.[ ]+(#.*)[ ]+:(.*)$")
		match := re.FindStringSubmatch(message)
		log.Printf("[IRC] List of channel for user %s , channel %s : %s", match[1], match[2], match[3])

		names := strings.Split(match[3], " ")

		for _, name := range names {
			userLine := "<span class=\"badge\">" + name + "</span>"

			if strings.Index(name, "+") == 0 {
				userLine = "<span class=\"badge\">" + strings.TrimLeft(name, "+") + "</span>"
			}

			if strings.Index(name, "@") == 0 {
				userLine = "<span class=\"badge\">" + strings.TrimLeft(name, "@") + "</span>"
			}

			files.MyStream.AddUniqueUser(userLine)

		}

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

	log.Printf("[IRC]->  <%s> ", message)

}
