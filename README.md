# ollyster
WEB-BASED IRC client, with a Social Flavour.

My challenge here was to proof something I realized time ago: a "social network" 
is just a web interface to a program which allows people to message each others.

If I am right, it is possible to make a social network out of every protocol which
allows users to message each others. Ostatus is an example, and this is why GnuSocial 
exists. 

Movim is the name of a social network built on top of Xmpp (Jabber).

I think we can have social-like clients for mail or NNTP protocol also, but I choose IRC,
since it is robust, widely used, I repeat **very** used (who says "IRC is dead" is just a stupid, 
go and check instead of lecturing your ignorance! ).

I will use the messaging protocol as a backbone, and the CTCP extension to add social-like
behaviors.

So:

    Phase 1: Web IRC client (Plan IRC).
    Phase 2: Wev IRC client + CTCP (Social features, like icons and so).


To install it:

just go into a folder, and _git clone https://github.com/uriel-fanelli/ollyster.git_

Then enter /etc and edit both the config file  (ollyster.conf) and the profile (profile.conf).

**ALL FIELDS ARE MANDATORY ** (I will make it more tolerant later)

Maybe you want to change the avatar also, in static/avatars/default.png 

Run "go build"
Run ./ollyster

Point your browser to the port you setup in ollyster.conf as _webport_, like  localhost:_webport_

Done.

##TODO

1. Groups page 
2. Form to edit configuration via web and not manually.
3. Form to edit profile via web and not manually.
4. Form to answer the messages on the "Timeline" page.
5. CTCP commands to get icons, email, website, and so, from the other ollyster clients.

