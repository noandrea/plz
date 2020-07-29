package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// will contain the stuff that we need
var csvPath, jsonPath string
var summary bool

// massageCmd represents the massage command
var massageCmd = &cobra.Command{
	Use:   "massage",
	Short: "Transform data from ESRI CSV to JSON suitable for serving zip info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// print the welcome screen
		fmt.Println(welcome)
		massage(csvPath, jsonPath, summary)
	},
}

func init() {
	rootCmd.AddCommand(massageCmd)
	massageCmd.Flags().StringVarP(&csvPath, "input", "i", "data.csv", "The input csv file (default data.csv)")
	massageCmd.Flags().StringVarP(&jsonPath, "output", "o", "data.json", "The output json file (default data.json)")
	massageCmd.Flags().BoolVarP(&summary, "summary", "s", false, "Print only execution summary")
}

// massage process input dataset to create an API optimized
// json to be served with the serve command
func massage(csvPath, jsonPath string, summary bool) (err error) {
	// begin processing
	start := time.Now()
	// open the csv file
	csvF, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Cannot open the file", csvPath, "for reading:", err)
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
		fmt.Println("Error reading CSV!", err)
		return
	}

	// output
	aggregate := make(map[string]map[string]int)
	// read all the lines
	var x int
	for {
		line, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("WARNING: reading line ", err)
			continue
		}

		// TODO: this slicing is an overly optimistic assumption that the row will always be there
		z, y := line[iZip], line[iYear][0:4]
		if c, xs := aggregate[z]; xs {
			c[y]++
		} else {
			aggregate[z] = map[string]int{y: 1}
		}
		// print progress
		x++
		if !summary {
			fmt.Printf("\rrecords %d", x)
		}

	}
	fmt.Println("Processed", x, "records")
	fmt.Println("Data read after ", time.Since(start))
	// transform
	// make two maps counts/distrib
	// make the years a sorted list
	// prepare aggregation

	type yc struct {
		Year  string `json:"year,omitempty"`
		Count int    `json:"count,omitempty"`
	}

	distrib := make(map[string][]yc, len(aggregate))
	counter := make(map[string]map[string]int, len(aggregate))
	for z, ycm := range aggregate {
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
			return ycs[i].Year < ycs[j].Year
		})
		// add the distribution to the result
		distrib[z] = ycs
		// add the sum to the result
		counter[z] = map[string]int{"total": sum}
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
		fmt.Println("Cannot open the file", jsonPath, ":", err)
		return
	}
	defer jsonF.Close()

	n, err := jsonF.Write(outB)
	if err != nil {
		fmt.Println("Error writing JSON! ", err)
		return
	}

	fmt.Printf("Output size (gb) is %.4f\n", (float64(n) / 1e9))
	fmt.Println("Output written at", jsonPath)
	fmt.Println("Completed in ", time.Since(start))
	return
}
