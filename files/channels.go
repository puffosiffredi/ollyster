package files

import (
	"io/ioutil"
	"strings"
)

// saves the full list of channels and descriptions
func (this *ollysterSocial) FlushChanList() {

	// initialize the group file if it doesn't exists
	// periodically flushes the channelbuf there

	ioutil.WriteFile(this.channelname, []byte(this.DumpChanList()), 0755)

}

// Initializechanlist writes an initial entry for the channel
func (this *ollysterSocial) InitializeChanList() {

	this.OChannels = make(map[string]string)

	this.OChannels["#placeholder"] = "This is a placeholder, until the very list of channels arrives"

}

// AddUniqueChannel adds a channel to the channel buffer
func (this *ollysterSocial) AddUniqueChannel(channel_name string, channel_desc string) error {

	this.OChannels[channel_name] = channel_desc

	return nil

}

// returns the channel table out of a ChannelBuf which contains only the name of the channel.
func (this *ollysterSocial) DumpChanList() string {

	var tmp_buf string

	delete(this.OChannels, "#placeholder")

	for name, desc := range this.OChannels {

		col_1 := "<td class=\"col-md-2\"><b>" + name + "</b></td>"
		col_2 := "<td class=\"col-md-3\">" + desc + "</td>"
		col_3 := "<td class=\"col-md-1\"><a href=\"/get/addgroup/" + strings.Replace(name, "#", "@", -1) + "\"><span class=\"label label-primary\">Subscribe</span></a></td>"
		list_line := "<tr>" + col_1 + col_2 + col_3 + "</tr>\n"

		tmp_buf += list_line

	}

	return tmp_buf

}
