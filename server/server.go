package main

import (
	"github.com/anacrolix/torrent/metainfo"
	"fmt"
	"path/filepath"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent/bencode"
	"os"
)

func main() {
	path:="torrentData"

	clientConfig:=torrent.Config{}
	clientConfig.Seed=true
	clientConfig.Debug=true
	clientConfig.DisableTrackers=true
	clientConfig.ListenAddr="127.0.0.1:6666"
	clientConfig.DHTConfig=dht.ServerConfig{
		StartingNodes:dht.GlobalBootstrapAddrs,
	}
	clientConfig.DataDir=path
	clientConfig.DisableAggressiveUpload=false
	client,_:=torrent.NewClient(&clientConfig)

	dir,_:=os.Open(path)
	defer dir.Close()

	fi,_ :=dir.Readdir(-1)
	for _,x:=range fi{
		if !x.IsDir() && x.Name()!=".torrent.bolt.db"{
			d:=makeMagnet(path,x.Name(),client)
			fmt.Println(d)
		}
	}

	fmt.Println(len(client.Torrents()))
	select{}
}

func makeMagnet(  dir string, name string,cl *torrent.Client) string {
	mi := metainfo.MetaInfo{}
	mi.SetDefaults()
	info := metainfo.Info{PieceLength: 1024*1024}
	info.BuildFromFilePath(filepath.Join(dir, name))
	mi.InfoBytes, _ = bencode.Marshal(info)
	cl.AddTorrent(&mi)
	magnet := mi.Magnet(name, mi.HashInfoBytes()).String()
	return magnet
}