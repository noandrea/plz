package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// keys for accessing data
const (
	keyCounters = "counters"
	keyDistrib  = "distributions"
)

// the path to the json data file to serve
var dataFilePath string
var listenAddress string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve the zip api",
	Long:  ``,
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&listenAddress, "listen", ":2007", "The address to listen to (default 0.0.0.0:2007)")
	serveCmd.Flags().StringVar(&dataFilePath, "data", "data.json", "The json data file to serve (default data.json)")

}

// listen starts the web server
func serve(cmd *cobra.Command, args []string) {
	fmt.Println(welcome)
	// open the database
	log.Info("serving data from ", dataFilePath)
	data, err := loadData(dataFilePath)
	if err != nil {
		fmt.Println("Error starting PLZ", err)
		log.Error(err)
		return
	}
	// echo start
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// health check :)
	e.GET("/status", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"version":   rootCmd.Version,
			"zip_codes": len(data[keyCounters]),
		})
	})
	e.GET("/zip/buildings", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, data[keyCounters])
	})
	e.GET("/zip/buildings/:code", func(c echo.Context) (err error) {
		zip := c.Param("code")
		if r, found := data[keyCounters][zip]; found {
			return c.JSON(http.StatusOK, r)
		}
		return c.JSON(http.StatusNotFound, map[string]string{})
	})
	e.GET("/zip/buildings/history", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, data[keyDistrib])
	})
	e.GET("/zip/buildings/:code/history", func(c echo.Context) (err error) {
		zip := c.Param("code")
		if r, found := data[keyDistrib][zip]; found {
			return c.JSON(http.StatusOK, r)
		}
		return c.JSON(http.StatusNotFound, map[string]string{})
	})
	err = e.Start(listenAddress)
	if err != nil {
		fmt.Println("Error starting PLZ", err)
		log.Error(err)
	}
}

// load json data to a map
func loadData(path string) (data map[string]map[string]interface{}, err error) {
	start := time.Now()
	// open the file
	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}
	// defer the closing
	defer jsonFile.Close()
	// print size
	if i, e := jsonFile.Stat(); e == nil {
		log.Debug("File ", path, " size(b):", i.Size())
	} else {
		log.Debug("Cannot get file stats for ", path, ": ", err)
	}
	// read it all in !!
	raw, err := ioutil.ReadAll(jsonFile)
	// init the data
	data = make(map[string]map[string]interface{}, 2)
	// fill it up
	err = json.Unmarshal(raw, &data)
	log.Debug("data loaded in ", time.Since(start))
	return
}
