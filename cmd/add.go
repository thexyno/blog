/*
Copyright © 2022 Philipp Hochkamp <git@phochkamp.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/server"
	"github.com/thexyno/xynoblog/util"
	"gopkg.in/yaml.v3"
)

type mdheader struct {
	Title   string    `yaml:"title"`
	Id      string    `yaml:"id"`
	Created time.Time `yaml:"created"`
	Updated time.Time `yaml:"updated"`
	Tags    []string  `yaml:"tags"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a post",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dbc := db.NewDb(viper.GetString(dbURIKey))
		mediaDir := viper.GetString(mediaDirKey)
		dbc.Seed()
		for _, arg := range args {
			filepath.Walk(arg, func(fpath string, info os.FileInfo, err error) error {
				if !info.IsDir() && strings.HasSuffix(fpath, ".md") {
					md, err := ioutil.ReadFile(fpath)
					if err != nil {
						log.Error(err)
						return err
					}
					mdsplit := strings.Split(string(md), "---")
					yml := mdsplit[1]
					var header mdheader

					err = yaml.Unmarshal([]byte(yml), &header)
					if err != nil {
						log.Error(err)
						return err
					}
					mdRest := strings.Join(mdsplit[2:], "")
					mdRest, err = util.AddMedia(mdRest, mediaDir, path.Dir(fpath))
					if err != nil {
						log.Error(err)
						return err
					}
					server.RenderSimple([]byte(mdRest)) // Panics when md is broken

					if err := dbc.Add(db.Post{
						Title:   header.Title,
						Id:      db.PostId(header.Id),
						Content: mdRest,
						Created: header.Created,
						Updated: header.Updated,
						Tags:    header.Tags,
					}); err != nil {
						log.Error(err)
						return err
					}
				}
				return nil
			})

		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
