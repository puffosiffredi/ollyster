/*
 * Copyright (c) 2015, Shinya Yagyu
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package node

import (
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/shingetsu-gou/shingetsu-gou/util"
)

const (
	defaultNodes = 5 // Nodes keeping in node list
	shareNodes   = 5 // Nodes having the file
)

//ManagerConfig contains params for NodeManager struct.
type ManagerConfig struct {
	Lookup    string
	Fmutex    *sync.RWMutex
	NodeAllow *util.RegexpList
	NodeDeny  *util.RegexpList
	Myself    *Myself
	InitNode  *util.ConfList
}

//Manager represents the map that maps datfile to it's source node list.
type Manager struct {
	*ManagerConfig
	isDirty bool
	nodes   map[string]Slice //map[""] is nodelist
	mutex   sync.RWMutex
}

//NewManager read the file and returns NodeManager obj.
func NewManager(cfg *ManagerConfig) *Manager {
	r := &Manager{
		ManagerConfig: cfg,
		nodes:         make(map[string]Slice),
	}
	err := util.EachKeyValueLine(cfg.Lookup, func(key string, value []string, i int) error {
		var nl Slice
		for _, v := range value {
			if v == "" {
				continue
			}
			nn, err := newNode(v)
			if err != nil {
				log.Println("line", i, "in lookup.txt,err=", err, v)
				continue
			}
			nl = append(nl, nn)
		}
		r.nodes[key] = nl
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return r
}

//getFromList returns number=n in the nodelist.
func (lt *Manager) getFromList(n int) *Node {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	if lt.ListLen() == 0 {
		return nil
	}
	return lt.nodes[""][n]
}

//NodeLen returns size of all nodes.
func (lt *Manager) NodeLen() int {
	ns := lt.getAllNodes()
	return ns.Len()
}

//ListLen returns size of nodelist.
func (lt *Manager) ListLen() int {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	return len(lt.nodes[""])
}

//GetNodestrSlice returns Nodestr of all nodes.
func (lt *Manager) GetNodestrSlice() []string {
	return lt.getAllNodes().getNodestrSlice()
}

//getAllNodes returns all nodes in table.
func (lt *Manager) getAllNodes() Slice {
	var n Slice
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	for _, v := range lt.nodes {
		for _, node := range v {
			n = append(n, node)
		}
	}
	return n.uniq()
}

//GetNodestrSliceInTable returns Nodestr slice of nodes associated datfile thread.
func (lt *Manager) GetNodestrSliceInTable(datfile string) []string {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	n := lt.nodes[datfile]
	return n.getNodestrSlice()
}

//Random selects # of min(all # of nodes,n) nodes randomly except exclude nodes.
func (lt *Manager) Random(exclude Slice, num int) []*Node {
	all := lt.getAllNodes()
	if exclude != nil {
		cand := make([]*Node, 0, len(all))
		m := exclude.toMap()
		for _, n := range all {
			if _, exist := m[n.Nodestr]; !exist {
				cand = append(cand, n)
			}
		}
		all = cand
	}
	n := all.Len()
	if num < n && num != 0 {
		n = num
	}
	r := make([]*Node, n)
	rs := rand.Perm(all.Len())
	for i := 0; i < n; i++ {
		r[i] = all[rs[i]]
	}
	return r
}

//AppendToTable add node n to table if it is allowd and list doesn't have it.
func (lt *Manager) AppendToTable(datfile string, n *Node) {
	lt.mutex.RLock()
	l := len(lt.nodes[datfile])
	lt.mutex.RUnlock()
	if ((datfile != "" && l < shareNodes) || (datfile == "" && l < defaultNodes)) &&
		n != nil && n.IsAllowed() && !lt.hasNodeInTable(datfile, n) {
		lt.mutex.Lock()
		lt.isDirty = true
		lt.nodes[datfile] = append(lt.nodes[datfile], n)
		lt.mutex.Unlock()
	}
}

//extendTable adds slice of nodes with check.
func (lt *Manager) extendToTable(datfile string, ns []*Node) {
	if ns == nil {
		return
	}
	for _, n := range ns {
		lt.AppendToTable(datfile, n)
	}
}

//appendToList add node n to nodelist if it is allowd and list doesn't have it.
func (lt *Manager) appendToList(n *Node) {
	lt.AppendToTable("", n)
}

//ReplaceNodeInList removes one node and say bye to the node and add n in nodelist.
//if len(node)>defaultnode
func (lt *Manager) ReplaceNodeInList(n *Node) *Node {
	lt.mutex.RLock()
	l := len(lt.nodes[""])
	lt.mutex.RUnlock()
	if !n.IsAllowed() || lt.hasNodeInTable("", n) {
		return nil
	}
	var old *Node
	if l >= defaultNodes {
		old := lt.getFromList(0)
		lt.RemoveFromList(old)
		old.bye()
	}
	lt.appendToList(n)
	return old
}

//extendToList adds node slice to nodelist.
func (lt *Manager) extendToList(ns []*Node) {
	lt.extendToTable("", ns)
}

//hasNode returns true if nodelist in all tables has n.
func (lt *Manager) hasNode(n *Node) bool {
	return len(lt.findNode(n)) > 0
}

//findNode returns datfile of node n, or -1 if not exist.
func (lt *Manager) findNode(n *Node) []string {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	var r []string
	for k := range lt.nodes {
		if lt.hasNodeInTable(k, n) {
			r = append(r, k)
		}
	}
	return r
}

//hasNodeInTable returns true if nodelist has n.
func (lt *Manager) hasNodeInTable(datfile string, n *Node) bool {
	return lt.findNodeInTable(datfile, n) != -1
}

//findNode returns location of node n, or -1 if not exist.
func (lt *Manager) findNodeInTable(datfile string, n *Node) int {
	return util.FindString(lt.GetNodestrSliceInTable(datfile), n.Nodestr)
}

//RemoveFromTable removes node n and return true if exists.
//or returns false if not exists.
func (lt *Manager) RemoveFromTable(datfile string, n *Node) bool {
	lt.mutex.Lock()
	defer lt.mutex.Unlock()
	i := 0
	if n != nil {
		i = util.FindString(lt.nodes[datfile].getNodestrSlice(), n.Nodestr)
	} else {
		for ii, nn := range lt.nodes[datfile] {
			if nn == nil {
				i = ii
				break
			}
		}
	}
	if i >= 0 {
		ln := len(lt.nodes[datfile])
		lt.nodes[datfile], lt.nodes[datfile][ln-1] = append(lt.nodes[datfile][:i], lt.nodes[datfile][i+1:]...), nil
		lt.isDirty = true
		return true
	}
	return false
}

//RemoveFromList removes node n from nodelist and return true if exists.
//or returns false if not exists.
func (lt *Manager) RemoveFromList(n *Node) bool {
	return lt.RemoveFromTable("", n)
}

//RemoveFromAllTable removes node n from all tables and return true if exists.
//or returns false if not exists.
func (lt *Manager) RemoveFromAllTable(n *Node) bool {
	del := false
	lt.mutex.RLock()
	for k := range lt.nodes {
		lt.mutex.RUnlock()
		del = del || lt.RemoveFromTable(k, n)
		lt.mutex.RLock()
	}
	lt.mutex.RUnlock()
	return del
}

//moreNodes gets another node info from each nodes in nodelist.
func (lt *Manager) moreNodes() {
	const retry = 5 // Times; Common setting

	no := 0
	count := 0
	all := lt.getAllNodes()
	for lt.ListLen() < defaultNodes {
		nn := all[no]
		newN, err := nn.getNode()
		if err == nil {
			if (lt.Myself.GetStatus() == Port0 && !lt.Myself.IsRelayed()) || lt.Join(newN) {
				all = append(all, newN)
				lt.appendToList(newN)
			}
		}
		if count++; count > retry {
			count = 0
			if no++; no >= len(all) {
				return
			}
		}
	}
}

//Initialize pings one of initNode except myself and added it if success,
//and get another node info from each nodes in nodelist.
func (lt *Manager) Initialize(rundir string) {
	var confile string
	con := "opened"
	if rundir != "" {
		confile = filepath.Join(rundir, "connection.dat")
		c, err := ioutil.ReadFile(confile)
		if err != nil {
			log.Println(err)
		}
		s := string(c)
		s = strings.Trim(s, "\r\n")
		if s == "uPnP" {
			log.Println("using uPnP as prevous.")
			lt.Myself.useUPnP()
			con = "uPnP"
		}
	}
	fn := []func([]*Node) (string, int){
		func(pingOK []*Node) (string, int) {
			log.Println("trying defaultport")
			lt.Myself.resetPort()
			return "opened", Normal
		},
		func(pingOK []*Node) (string, int) {
			log.Println("trying uPnP")
			lt.Myself.useUPnP()
			return "uPnP", Normal
		},
		func(pingOK []*Node) (string, int) {
			log.Println("trying relayed")
			lt.Myself.resetPort()
			seed, err := newNode(lt.InitNode.GetData()[0])
			if err != nil {
				log.Fatal(err)
			}
			<-lt.Myself.tryRelay(seed)
			return "relayed", Port0
		},
		func(pingOK []*Node) (string, int) {
			log.Println("failed to join")
			lt.Myself.setRelayServer(nil)
			for _, n := range pingOK {
				lt.appendToList(n)
			}
			return "failed", Port0
		},
	}
	stat := Normal
	for _, f := range fn {
		if ok, pingOK := lt.initialize(); !ok {
			con, stat = f(pingOK)
		} else {
			log.Println("success to join by", con)
			break
		}
	}
	if confile != "" {
		lt.Myself.setStatus(stat)
		err := ioutil.WriteFile(confile, []byte(con), 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func (lt *Manager) initialize() (bool, []*Node) {
	inodes := lt.Random(nil, defaultNodes)
	for _, i := range lt.InitNode.GetData() {
		nn, err := newNode(i)
		if err != nil {
			continue
		}
		inodes = append(inodes, nn)
	}
	var wg sync.WaitGroup
	pingOK := make([]*Node, 0, len(inodes))
	var mutex sync.Mutex
	for _, inode := range inodes {
		wg.Add(1)
		go func(inode *Node) {
			defer wg.Done()
			if _, err := inode.Ping(); err == nil {
				mutex.Lock()
				pingOK = append(pingOK, inode)
				mutex.Unlock()
				lt.Join(inode)
			}
		}(inode)
	}
	wg.Wait()
	if lt.ListLen() > 0 {
		lt.moreNodes()
	}
	log.Println("# of nodelist:", lt.ListLen())
	return lt.ListLen() > 0, pingOK
}

//Join tells n to join and adds n to nodelist if welcomed.
//if n returns another nodes, repeats it and return true..
//removes fron nodelist if not welcomed and return false.
func (lt *Manager) Join(n *Node) bool {
	const retryJoin = 2 // Times; Join network
	if n == nil {
		return false
	}
	flag := false
	if lt.hasNodeInTable("", n) || lt.Myself.IPPortPath().Nodestr == n.Nodestr {
		return false
	}
	for count := 0; count < retryJoin && lt.ListLen() < defaultNodes; count++ {
		extnode, err := n.join()
		if err == nil && extnode == nil {
			lt.appendToList(n)
			return true
		}
		if err == nil {
			lt.appendToList(n)
			n = extnode
			flag = true
		} else {
			lt.RemoveFromTable("", n)
			return flag
		}
	}
	return flag
}

//TellUpdate makes mynode info from node or dnsname or ip addr,
//and broadcast the updates of record id=id in cache c.datfile with stamp.
func (lt *Manager) TellUpdate(datfile string, stamp int64, id string, node *Node) {
	const updateNodes = 10

	tellstr := lt.Myself.toxstring()
	if node != nil {
		tellstr = node.toxstring()
	}
	msg := strings.Join([]string{"/update", datfile, strconv.FormatInt(stamp, 10), id, tellstr}, "/")

	ns := lt.Get(datfile, nil)
	ns = ns.Extend(lt.Get("", nil))
	ns = ns.Extend(lt.Random(ns, updateNodes))
	log.Println("telling #", len(ns))
	for _, n := range ns {
		_, err := n.Talk(msg, true, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

//Get returns rawnodelist associated with datfile
//if not found returns def
func (lt *Manager) Get(datfile string, def Slice) Slice {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	if v, exist := lt.nodes[datfile]; exist {
		nodes := make([]*Node, v.Len())
		copy(nodes, v)
		return nodes
	}
	return Slice(def)
}

//stringMap returns map of k=datfile, v=Nodestr of rawnodelist.
func (lt *Manager) stringMap() map[string][]string {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	result := make(map[string][]string)
	for k, v := range lt.nodes {
		if k == "" {
			continue
		}
		result[k] = v.getNodestrSlice()
	}
	return result
}

//Sync saves  k=datfile, v=Nodestr map to the file.
func (lt *Manager) Sync() {
	lt.mutex.RLock()
	isDirty := lt.isDirty
	lt.mutex.RUnlock()
	if isDirty {
		m := lt.stringMap()
		lt.Fmutex.Lock()
		defer lt.Fmutex.Unlock()
		err := util.WriteMap(lt.Lookup, m)
		if err != nil {
			log.Println(err)
		} else {
			lt.mutex.Lock()
			lt.isDirty = false
			lt.mutex.Unlock()
		}
	}
}

//NodesForGet returns nodes which has datfile cache , and that extends nodes to #searchDepth .
func (lt *Manager) NodesForGet(datfile string, searchDepth int) Slice {
	var ns, ns2 Slice
	ns = ns.Extend(lt.Get(datfile, nil))
	ns = ns.Extend(lt.Get("", nil))
	ns = ns.Extend(lt.Random(ns, 0))

	for _, n := range ns {
		if !n.Equals(lt.Myself.toNode()) && n.IsAllowed() {
			ns2 = append(ns2, n)
		}
	}
	if ns2.Len() > searchDepth {
		ns2 = ns2[:searchDepth]
	}
	return ns2
}

//Rejoin adds nodes in searchlist to nodelist if ping is ok and len(nodelist)<defaultNodes
//and doesn't have it's node.
//if ping is ng, removes node from searchlist.
func (lt *Manager) Rejoin() {
	all := lt.getAllNodes()
	for _, n := range all {
		if lt.ListLen() >= defaultNodes {
			return
		}
		lt.mutex.RLock()
		m := lt.nodes[""].toMap()
		_, has := m[n.Nodestr]
		lt.mutex.RUnlock()
		if has {
			continue
		}
		if _, err := n.Ping(); err == nil || !lt.Join(n) {
			lt.RemoveFromAllTable(n)
			lt.Sync()
		} else {
			lt.appendToList(n)
		}
	}
	log.Println("# of nodelist", lt.ListLen())
}

//PingAll pings to all nodes in nodelist.
//if ng, removes from nodelist.
func (lt *Manager) PingAll() {
	lt.mutex.RLock()
	var ns Slice
	for _, n := range lt.nodes[""] {
		ns = append(ns, n)
	}
	lt.mutex.RUnlock()
	var wg sync.WaitGroup
	for _, n := range ns {
		if n == nil {
			lt.RemoveFromAllTable(n)
			continue
		}
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			if _, err := n.Ping(); err != nil {
				lt.RemoveFromAllTable(n)
			}
		}(n)
	}
	wg.Wait()
}

//RejoinList joins all node in nodelist.
func (lt *Manager) RejoinList() {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	var wg sync.WaitGroup
	for _, n := range lt.nodes[""] {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			_, err := n.join()
			if err != nil {
				log.Println(err)
			}
		}(n)
	}
	wg.Wait()
}
