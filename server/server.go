package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	path := "data"

	clientConfig := torrent.Config{}
	clientConfig.Seed = true
	clientConfig.Debug = true
	clientConfig.DisableTrackers = true
	clientConfig.ListenAddr = "0.0.0.0:6666"
	clientConfig.DHTConfig = dht.ServerConfig{
		StartingNodes: serverAddrs,
	}
	clientConfig.DataDir = path
	clientConfig.DisableAggressiveUpload = false
	client, _ := torrent.NewClient(&clientConfig)

	dir, _ := os.Open(path)
	defer dir.Close()

	fi, _ := dir.Readdir(-1)
	for _, x := range fi {
		if !x.IsDir() && x.Name() != ".torrent.bolt.db" {
			d := makeMagnet(path, x.Name(), client)
			fmt.Println(d)
		}
	}

	fmt.Println(len(client.Torrents()))
	select {}
}

func makeMagnet(dir string, name string, cl *torrent.Client) string {
	mi := metainfo.MetaInfo{}
	mi.SetDefaults()
	info := metainfo.Info{PieceLength: 1024 * 1024}
	info.BuildFromFilePath(filepath.Join(dir, name))
	mi.InfoBytes, _ = bencode.Marshal(info)
	cl.AddTorrent(&mi)
	magnet := mi.Magnet(name, mi.HashInfoBytes()).String()
	return magnet
}

func serverAddrs() (addrs []dht.Addr, err error) {
	for _, s := range []string{
    "client:6666",
	} {
		ua, err := net.ResolveUDPAddr("udp4", s)
		if err != nil {
			continue
		}
		addrs = append(addrs, dht.NewAddr(ua))
	}
	if len(addrs) == 0 {
		err = errors.New("nothing resolved")
	}
	return
}
