package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// massageCmd represents the massage command
var massageCmd = &cobra.Command{
	Use:   "massage",
	Short: "Transform data from ESRI CSV to JSON suitable for serving zip info",
	Long:  ``,
	Run:   massage,
}

// will contain the stuff that we need
var csvPath, jsonPath string

func init() {
	rootCmd.AddCommand(massageCmd)
	massageCmd.Flags().StringVarP(&csvPath, "in", "i", "data.csv", "The input csv file (default data.csv)")
	massageCmd.Flags().StringVarP(&jsonPath, "out", "o", "data.json", "The output json file (default data.json)")
}

// massage process input dataset to create an API optimized
// json to be served with the serve command
func massage(cmd *cobra.Command, args []string) {
	// print the welcome screen
	fmt.Println(welcome)
	// begin processing
	start := time.Now()
	// open the csv file
	csvF, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Cannot open the file", csvPath, "for reading (use --debug for more info)")
		log.Debug("error ", err)
		return
	}
	// print the size
	csvI, err := csvF.Stat()
	if err == nil {
		fmt.Printf("Input file size (gb) is %.4f\n", (float64(csvI.Size()) / 1e9))
	}

	defer csvF.Close()
	//
	iZip, iYear := 15, 33
	// open a csv file reader
	csvR := csv.NewReader(csvF)
	csvR.TrimLeadingSpace = true
	// skip header
	_, err = csvR.Read()
	if err != nil {
		fmt.Println("Error reading CSV! (use --debug for more info)")
		log.Debug(err)
		return
	}

	// output
	dist := make(map[string]map[string]int)
	// read all the lines
	var x int
	for {
		line, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Error reading line ", err)
			continue
		}

		// TODO: this slicing is an overly optimistic assumption that the row will always be there
		z, y := line[iZip], line[iYear][0:4]
		if c, xs := dist[z]; xs {
			c[y]++
		} else {
			dist[z] = map[string]int{y: 1}
		}
		// print progress
		x++
		fmt.Printf("\rrecords %d", x)
	}
	fmt.Println()
	fmt.Println("Data read after ", time.Since(start))
	// transform
	// make two maps counts/distrib
	// make the years a sorted list
	// prepare aggregation

	type yc struct {
		year  string
		count int
	}

	distrib := make(map[string][]map[string]int, len(dist))
	counter := make(map[string]map[string]int, len(dist))
	for z, ycm := range dist {
		// sort and aggregate
		ycs, i := make([]yc, len(ycm)), 0
		sum := 0
		for y, c := range ycm {
			// transform the counters to a slice
			ycs[i] = yc{y, c}
			sum += c
			i++
		}
		// sort by year ascending
		sort.SliceStable(ycs, func(i, j int) bool {
			return ycs[i].year < ycs[j].year
		})
		// make it nice for the desired output
		ycj := make([]map[string]int, len(ycs))
		for i, yc := range ycs {
			ycj[i] = map[string]int{yc.year: yc.count}
		}
		// add the distribution to the result
		distrib[z] = ycj
		// add the sum to the result
		counter[z] = map[string]int{"count": sum}
	}
	fmt.Println("Data processed after ", time.Since(start))
	// transform everything in a big json blob
	outB, err := json.Marshal(map[string]interface{}{
		keyDistrib:  distrib,
		keyCounters: counter,
	})
	if err != nil {
		fmt.Println("Failed to encode data to Json", err)
		return
	}
	// finally write the output
	jsonF, err := os.Create(jsonPath)
	if err != nil {
		fmt.Println("Cannot open the file", jsonPath, "for writing (use --debug for more info)")
		log.Debug(err)
		return
	}
	defer jsonF.Close()

	n, err := jsonF.Write(outB)
	if err != nil {
		fmt.Println("Error writing JSON! (use --debug for more info)")
		log.Debug(err)
		return
	}

	fmt.Printf("Output size (gb) is %.4f\n", (float64(n) / 1e9))
	fmt.Println("Output written at", jsonPath)
	fmt.Println("Completed in ", time.Since(start))
}
