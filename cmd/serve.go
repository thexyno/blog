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
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thexyno/xynoblog/db"
	"github.com/thexyno/xynoblog/server"
)

const cssDirKey = "cssdir"
const staticDirKey = "staticdir"
const fontDirKey = "fontdir"
const listenKey = "listen"
const releaseModeKey = "releaseMode"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the xynoblog server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cssDir := viper.GetString(cssDirKey)
		fontDir := viper.GetString(fontDirKey)
		mediaDir := viper.GetString(mediaDirKey)
		staticDir := viper.GetString(staticDirKey)
		listen := viper.GetString(listenKey)
		dbURI := viper.GetString(dbURIKey)
		releaseMode := viper.GetBool(releaseModeKey)
		if _, err := os.Stat(mediaDir); err != nil {
			if os.IsNotExist(err) {
				if err = os.Mkdir(mediaDir, 0755); err != nil {
					log.Panic(err)
				}
			} else {
				log.Panic(err)
			}
		}

		log.Print(viper.GetString("fontdir"))
		log.Printf("Starting Xynoblog on %s", listen)
		database := db.NewDb(dbURI)
		database.Seed()
		log.Printf("Fontdir: %s, CSSDir: %s, StaticDir: %s, MediaDir: %s", fontDir, cssDir, staticDir, mediaDir)
		mux := server.Mux(database, fontDir, cssDir, staticDir, mediaDir)
		log.Printf("Started Xynoblog on %s", listen)
		if releaseMode {
			gin.SetMode(gin.ReleaseMode)
		}
		mux.Run(listen)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.Flags().String(fontDirKey, "", "Directory containing JetBrainsMono[wght].ttf")
	viper.BindPFlag(fontDirKey, serveCmd.Flags().Lookup(fontDirKey))

	serveCmd.Flags().String(cssDirKey, "./cssdist", "Directory containing output.css (default should be fine)")
	viper.BindPFlag(cssDirKey, serveCmd.Flags().Lookup(cssDirKey))
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	serveCmd.Flags().String(staticDirKey, wd, "Directory containing the data folder")
	viper.BindPFlag(staticDirKey, serveCmd.Flags().Lookup(staticDirKey))

	serveCmd.Flags().String(listenKey, ":8080", "host/port to listen on")
	viper.BindPFlag(listenKey, serveCmd.Flags().Lookup(listenKey))

	serveCmd.Flags().Bool(releaseModeKey, false, "Set true for production")
	viper.BindPFlag(releaseModeKey, serveCmd.Flags().Lookup(releaseModeKey))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
